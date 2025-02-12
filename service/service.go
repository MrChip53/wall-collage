package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"time"

	"wall-collage/notif"
	"wall-collage/pb"
)

const tmpFolder = "/tmp/wall-collage"

type service struct {
	pb.UnimplementedWallCollageServer

	folders []string

	isRunning     bool
	delay         int
	bgColor       string
	enableCollage bool

	imgs []string

	notificationService *notif.NotificationService

	collageCtx    context.Context
	collageCancel context.CancelFunc
}

func NewService(folderPath string) (*service, error) {
	err := makeFolder(tmpFolder)
	if err != nil {
		return nil, err
	}

	return &service{
		isRunning:     false,
		delay:         5,
		bgColor:       "#000000",
		enableCollage: true,
		folders:       []string{folderPath},
	}, nil
}

func (s *service) startCollageService() error {
	if s.isRunning {
		return nil
	}

	s.scanFolders()

	if s.notificationService == nil {
		ns, err := notif.NewNotificationService()
		if err != nil {
			return err
		}
		s.notificationService = ns
	}

	s.isRunning = true
	s.collageCtx, s.collageCancel = context.WithCancel(context.Background())
	go s.collageLoop(s.collageCtx)

	return nil
}

func (s *service) stopCollageService() error {
	if !s.isRunning {
		return nil
	}

	s.isRunning = false
	s.collageCancel()

	return nil
}

func (s *service) collageLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := s.setWallpaper(s.imgs)
			if err != nil {
				log.Printf("Error setting wallpaper: %v", err)
			}
			if s.delay == 0 {
				break
			}
			time.Sleep(time.Duration(s.delay) * time.Second)
		}
	}
}

func (s *service) setWallpaper(imgPaths []string) error {
	var single bool
	var imgPath string
	var err error

	if len(imgPaths) == 0 {
		bgCmd := fmt.Sprintf("hsetroot -solid \"%s\"", s.bgColor)
		err = exec.Command("sh", "-c", bgCmd).Run()
		if err != nil {
			return err
		}

		return nil
	}

	imgPath = imgPaths[rand.Intn(len(imgPaths))]
	if s.enableCollage {
		imgPath, single, err = s.createCollage(s.getRandomImages(imgPaths, 3))
		if err != nil {
			log.Printf("Error creating collage: %v", err)
			imgPath = imgPaths[rand.Intn(len(imgPaths))]
		}
	}

	mode := "full"
	if single {
		mode = "fill"
	}

	bgCmd := fmt.Sprintf("hsetroot -solid \"%s\" -%s \"%s\"", s.bgColor, mode, imgPath)
	err = exec.Command("sh", "-c", bgCmd).Run()
	if err != nil {
		return err
	}
	return nil
}

func (s *service) getRandomImages(imgPaths []string, num int) []string {
	result := make([]string, 0)
	for i := 0; i < num; i++ {
		r, single := s.getRandomImage(imgPaths, result)
		result = append(result, r)
		if !s.enableCollage || single {
			return result
		}
	}
	return result
}

func (s *service) getRandomImage(imgPaths []string, list []string) (string, bool) {
	p := imgPaths[rand.Intn(len(imgPaths))]

	h, w, err := getImageHeightWidth(p)
	if err != nil {
		return s.getRandomImage(imgPaths, list)
	}

	single := !s.enableCollage || isSingleFile(p) || w > h
	if (single && len(list) > 0) || slices.Contains(list, p) {
		return s.getRandomImage(imgPaths, list)
	}

	return p, single
}

func (s *service) scanFolders() {
	if len(s.folders) == 0 {
		s.imgs = make([]string, 0)
		return
	}

	imgs := make([]string, 0)

	for _, folder := range s.folders {
		imgPaths, err := s.scanFolder(folder)
		if err != nil {
			log.Printf("Error scanning folder: %v", err)
			continue
		}
		imgs = append(imgs, imgPaths...)
	}

	s.imgs = imgs
}

func (s *service) scanFolder(folder string) ([]string, error) {
	imgPaths := make([]string, 0)

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != folder {
			return filepath.SkipDir
		}
		if !info.IsDir() && isImageFile(path) {
			imgPaths = append(imgPaths, path)
		}
		return nil
	})
	return imgPaths, err
}

func (s *service) sendNotification(title, body string) {
	err := s.notificationService.Notify(title, body)
	if err != nil {
		log.Printf("Error sending notification: %v", err)
	}
}

func (s *service) Start(ctx context.Context, in *pb.StartRequest) (*pb.StartResponse, error) {
	s.sendNotification("Wall Collage", "Starting Wall Collage")
	err := s.startCollageService()
	if err != nil {
		return nil, err
	}
	return &pb.StartResponse{}, nil
}

func (s *service) Stop(ctx context.Context, in *pb.StopRequest) (*pb.StopResponse, error) {
	s.sendNotification("Wall Collage", "Stopping Wall Collage")
	err := s.stopCollageService()
	if err != nil {
		return nil, err
	}
	return &pb.StopResponse{}, nil
}

func (s *service) Status(ctx context.Context, in *pb.StatusRequest) (*pb.StatusResponse, error) {
	log.Printf("Status request received")
	return &pb.StatusResponse{
		IsRunning:       s.isRunning,
		Delay:           int32(s.delay),
		BackgroundColor: s.bgColor,
		Collage:         s.enableCollage,
	}, nil
}

func (s *service) SetDelay(ctx context.Context, in *pb.SetDelayRequest) (*pb.SetDelayResponse, error) {
	s.sendNotification("Wall Collage", fmt.Sprintf("Setting delay to %d seconds", in.Delay))
	s.delay = int(in.Delay)
	return &pb.SetDelayResponse{}, nil
}

func (s *service) SetBackgroundColor(ctx context.Context, in *pb.SetBackgroundColorRequest) (*pb.SetBackgroundColorResponse, error) {
	s.sendNotification("Wall Collage", fmt.Sprintf("Setting background color to %s", in.Color))
	s.bgColor = in.Color
	return &pb.SetBackgroundColorResponse{}, nil
}

func (s *service) ToggleCollage(ctx context.Context, in *pb.ToggleCollageRequest) (*pb.ToggleCollageResponse, error) {
	s.sendNotification("Wall Collage", fmt.Sprintf("Setting collage to %t", !s.enableCollage))
	s.enableCollage = !s.enableCollage
	return &pb.ToggleCollageResponse{Enabled: s.enableCollage}, nil
}

func (s *service) ListFolders(ctx context.Context, in *pb.ListFoldersRequest) (*pb.ListFoldersResponse, error) {
	log.Printf("List folders request received")
	return &pb.ListFoldersResponse{Folders: s.folders}, nil
}

func (s *service) AddFolder(ctx context.Context, in *pb.AddFolderRequest) (*pb.AddFolderResponse, error) {
	if slices.Contains(s.folders, in.Folder) {
		return &pb.AddFolderResponse{}, nil
	}

	s.folders = append(s.folders, in.Folder)
	s.scanFolders()
	s.sendNotification("Wall Collage", fmt.Sprintf("Added folder %s", in.Folder))
	return &pb.AddFolderResponse{}, nil
}

func (s *service) RemoveFolder(ctx context.Context, in *pb.RemoveFolderRequest) (*pb.RemoveFolderResponse, error) {
	if int(in.FolderIndex) > len(s.folders)-1 {
		return nil, fmt.Errorf("Folder index out of range")
	}

	newFolders := make([]string, 0)

	for i, folder := range s.folders {
		if i == int(in.FolderIndex) {
			continue
		}

		newFolders = append(newFolders, folder)
	}

	s.folders = newFolders
	s.scanFolders()
	s.sendNotification("Wall Collage", fmt.Sprintf("Removed folder %d", in.FolderIndex))
	return &pb.RemoveFolderResponse{}, nil
}
