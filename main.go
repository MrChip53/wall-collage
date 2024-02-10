package main

import (
	"flag"
	"fmt"
	"github.com/edwvee/exiffix"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"os/exec"
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
	lockfile := lockFile()
	defer unlockFile(lockfile)

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
		os.Exit(1)
	}

	imgPaths, err := scanFolder(folder)
	if err != nil {
		fmt.Println("Error scanning folder:", err)
		os.Exit(1)
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

func setWallpaper(imgPaths []string) error {
	var single bool
	var imgPath string
	var err error

	imgPath = imgPaths[rand.Intn(len(imgPaths))]
	if collage {
		imgPath, single, err = createCollage(getRandomImages(imgPaths, 3))
		if err != nil {
			fmt.Println("Error creating collage:", err)
			imgPath = imgPaths[rand.Intn(len(imgPaths))]
		}
	}

	mode := "full"
	if single || onlyHidden || isHiddenFile(imgPath) {
		mode = "fill"
	}

	bgCmd := fmt.Sprintf("hsetroot -solid \"%s\" -%s \"%s\"", solid, mode, imgPath)
	err = exec.Command("sh", "-c", bgCmd).Run()
	if err != nil {
		return err
	}
	return nil
}

func createCollage(imgPaths []string) (string, bool, error) {
	//imgPaths = []string{
	//	"/home/mrchip/.wallpapers/.20210602_163950.jpg",
	//}

	if len(imgPaths) == 1 {
		return imgPaths[0], true, nil
	}

	imageFiles := make([]image.Image, 0)
	for _, imgPath := range imgPaths {
		file, err := os.Open(imgPath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		img, _, err := exiffix.Decode(file)
		if err != nil {
			return imgPaths[0], true, err
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
		return imgPaths[0], true, err
	}
	defer output.Close()

	err = png.Encode(output, canvas)
	if err != nil {
		return imgPaths[0], true, err
	}

	return folder + "/wall-collage/collage.png", false, nil
}
