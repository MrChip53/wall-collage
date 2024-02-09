package main

import (
	"flag"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

var folder string
var hidden bool
var collage bool
var onlyHidden bool
var solid string
var sleep int

var screenWidth = 1920
var screenHeight = 1080

func main() {
	flag.StringVar(&solid, "s", "#000000", "solid color for wallpaper")
	flag.BoolVar(&hidden, "h", false, "use hidden files in folder")
	flag.StringVar(&folder, "f", "", "wallpaper folder path")
	flag.BoolVar(&collage, "c", false, "create a collage of three images in folder")
	flag.BoolVar(&onlyHidden, "o", false, "only use hidden files in folder")
	flag.IntVar(&sleep, "t", 0, "time to change wallpaper in seconds")
	flag.Parse()

	if folder == "" {
		fmt.Println("Please provide a folder path")
		return
	} else if strings.HasSuffix(folder, "/") {
		folder = folder[:len(folder)-1]
	}

	err := makeFolder(folder + "/wall-collage")
	if err != nil {
		fmt.Println("Error creating wall-collage folder:", err)
	}

	var imgPaths = make([]string, 0)

	err = filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != folder {
			return filepath.SkipDir
		}
		if !info.IsDir() && isImageFile(path) && ((isHiddenFile(path) && (hidden || onlyHidden)) || (!isHiddenFile(path) && !onlyHidden)) {
			imgPaths = append(imgPaths, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error:", err)
	}

	for {
		err = setWallpaper(imgPaths)
		if err != nil {
			fmt.Println("Error setting wallpaper:", err)
		}
		if sleep == 0 {
			break
		}
		time.Sleep(time.Duration(sleep) * time.Second)
	}
}

func isHiddenFile(path string) bool {
	return strings.HasPrefix(filepath.Base(path), ".")
}

func isImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		return true
	}
	return false
}

func setWallpaper(imgPaths []string) error {
	imgPath := imgPaths[rand.Intn(len(imgPaths))]
	if collage {
		p, err := createCollage(getRandomImages(imgPaths, 3))
		if err != nil {
			fmt.Println("Error creating collage:", err)
		} else {
			imgPath = p
		}
	}

	mode := "full"
	if onlyHidden || isHiddenFile(imgPath) {
		mode = "fill"
	}

	bgCmd := fmt.Sprintf("hsetroot -solid \"%s\" -%s \"%s\"", solid, mode, imgPath)
	err := exec.Command("sh", "-c", bgCmd).Run()
	if err != nil {
		return err
	}
	return nil
}

func getRandomImage(imgPaths []string, allowHidden bool, allowVisible bool, list []string) (string, bool) {
	p := imgPaths[rand.Intn(len(imgPaths))]
	isHidden := isHiddenFile(p)
	if (isHidden && !allowHidden) || (isHidden && len(list) > 0) || (!isHidden && !allowVisible) || slices.Contains(list, p) {
		return getRandomImage(imgPaths, allowHidden, allowVisible, list)
	}
	return p, isHidden
}

func getRandomImages(imgPaths []string, num int) []string {
	var result = make([]string, 0)
	for i := 0; i < 3; i++ {
		r, isHidden := getRandomImage(imgPaths, hidden || onlyHidden, !onlyHidden, result)
		result = append(result, r)
		if isHidden {
			return result
		}
	}
	return result
}

func createCollage(imgPaths []string) (string, error) {
	if len(imgPaths) == 1 {
		return imgPaths[0], nil
	}

	imageFiles := make([]image.Image, 0)
	for _, imgPath := range imgPaths {
		file, err := os.Open(imgPath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			return imgPaths[0], err
		}
		img = resize.Resize(640, 1080, img, resize.Lanczos3)
		imageFiles = append(imageFiles, img)
	}

	canvas := image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))

	draw.Draw(canvas, image.Rect(0, 0, 640, 1080), imageFiles[0], image.ZP, draw.Src)
	draw.Draw(canvas, image.Rect(640, 0, 1280, 1080), imageFiles[1], image.ZP, draw.Src)
	draw.Draw(canvas, image.Rect(1280, 0, 1920, 1080), imageFiles[2], image.ZP, draw.Src)

	output, err := os.Create(folder + "/wall-collage/collage.png")
	if err != nil {
		return imgPaths[0], err
	}
	defer output.Close()

	err = png.Encode(output, canvas)
	if err != nil {
		return imgPaths[0], err
	}

	return folder + "/wall-collage/collage.png", nil
}

func makeFolder(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
