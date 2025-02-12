package service

import (
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"

	"github.com/edwvee/exiffix"
	"github.com/nfnt/resize"
)

var (
	screenWidth  = 1920
	screenHeight = 1080
)

func getImageHeightWidth(imgPath string) (int, int, error) {
	file, err := os.Open(imgPath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	img, _, err := exiffix.Decode(file)
	if err != nil {
		return 0, 0, err
	}

	return img.Bounds().Dy(), img.Bounds().Dx(), nil
}

func (s *service) createCollage(imgPaths []string) (string, bool, error) {
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

	draw.Draw(canvas, image.Rect(0, 0, 640, 1080), imageFiles[0], image.Point{}, draw.Src)
	draw.Draw(canvas, image.Rect(640, 0, 1280, 1080), imageFiles[1], image.Point{}, draw.Src)
	draw.Draw(canvas, image.Rect(1280, 0, 1920, 1080), imageFiles[2], image.Point{}, draw.Src)

	output, err := os.Create(tmpFolder + "/collage.png")
	if err != nil {
		return imgPaths[0], true, err
	}
	defer output.Close()

	err = png.Encode(output, canvas)
	if err != nil {
		return imgPaths[0], true, err
	}

	return tmpFolder + "/collage.png", false, nil
}
