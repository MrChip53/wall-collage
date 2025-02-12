package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"wall-collage/pb"
	"wall-collage/service"
)

const socketPath = "/tmp/wall-collage.sock"

var folderPath string

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Start the wall-collage daemon",
	Long:  "Start the wall-collage daemon",
	Run: func(cmd *cobra.Command, args []string) {
		if folderPath == "" {
			fmt.Println("Please provide a folder path")
			return
		} else if strings.HasSuffix(folderPath, "/") {
			folderPath = folderPath[:len(folderPath)-1]
		}

		if err := os.RemoveAll(socketPath); err != nil {
			panic(err)
		}

		listener, err := net.Listen("unix", socketPath)
		if err != nil {
			panic(err)
		}
		defer listener.Close()

		service, err := service.NewService(folderPath)
		if err != nil {
			panic(err)
		}

		_, err = service.Start(context.Background(), &pb.StartRequest{})
		if err != nil {
			log.Fatalf("Could not start collage service: %v", err)
		}

		grpcServer := grpc.NewServer()
		pb.RegisterWallCollageServer(grpcServer, service)

		log.Printf("Server listening on socket %s", socketPath)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	},
}

func init() {
	daemonCmd.Flags().StringVarP(&folderPath, "folder", "f", "", "Folder path to scan for images")
	rootCmd.AddCommand(daemonCmd)
}
