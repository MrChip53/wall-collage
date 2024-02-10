package main

import (
	"github.com/edwvee/exiffix"
	"math/rand"
	"os"
	"slices"
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

func getRandomImages(imgPaths []string, num int) []string {
	var result = make([]string, 0)
	for i := 0; i < 3; i++ {
		r, isHidden, single := getRandomImage(imgPaths, hidden || onlyHidden, !onlyHidden, result)
		result = append(result, r)
		if (isHidden && !collage) || single {
			return result
		}
	}
	return result
}

func getRandomImage(imgPaths []string, allowHidden bool, allowVisible bool, list []string) (string, bool, bool) {
	p := imgPaths[rand.Intn(len(imgPaths))]

	isHidden := isHiddenFile(p)
	h, w, err := getImageHeightWidth(p)
	if err != nil {
		return getRandomImage(imgPaths, allowHidden, allowVisible, list)
	}
	single := (isHidden && !collage) || isSingleFile(p) || w > h
	if (isHidden && !allowHidden) || (single && len(list) > 0) || (!isHidden && !allowVisible) || slices.Contains(list, p) {
		return getRandomImage(imgPaths, allowHidden, allowVisible, list)
	}
	return p, isHidden, single
}
