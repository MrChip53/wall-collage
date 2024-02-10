package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func scanFolder(folder string) ([]string, error) {
	var imgPaths = make([]string, 0)

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
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
	return imgPaths, err
}

func isImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		return true
	}
	return false
}

func isSingleFile(path string) bool {
	return strings.Contains(filepath.Base(path), ".single.")
}

func isHiddenFile(path string) bool {
	return strings.HasPrefix(filepath.Base(path), ".")
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

func lockFile() *os.File {
	lockfile, err := os.OpenFile("/tmp/wall-collage.lock", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Failed to open lock file:", err)
		os.Exit(1)
	}

	err = syscall.Flock(int(lockfile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		fmt.Println("Another instance is already running.")
		os.Exit(1)
	}
	return lockfile
}

func unlockFile(lockfile *os.File) {
	err := syscall.Flock(int(lockfile.Fd()), syscall.LOCK_UN)
	if err != nil {
		fmt.Println("Failed to unlock file:", err)
	}
	lockfile.Close()
}
