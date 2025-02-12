package cmd

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"wall-collage/client"
	"wall-collage/pb"
)

var folderCmd = &cobra.Command{
	Use:   "folder",
	Short: "Folder settings",
}

var listFolderCmd = &cobra.Command{
	Use:   "ls",
	Short: "List folders",
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := client.NewClient(socketPath)
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		r, err := client.ListFolders(ctx, &pb.ListFoldersRequest{})
		if err != nil {
			log.Fatalf("could not list folders: %v", err)
		}

		for i, folder := range r.GetFolders() {
			fmt.Printf("%d: %s\n", i, folder)
		}
	},
}

var addFolderCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new folder",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatalf("usage: wall-collage folder add <folder>")
		}

		client, conn := client.NewClient(socketPath)
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err := client.AddFolder(ctx, &pb.AddFolderRequest{Folder: args[0]})
		if err != nil {
			log.Fatalf("could not add folder: %v", err)
		}
	},
}

var removeFolderCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a folder by index",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatalf("usage: wall-collage folder rm <index>")
		}

		idx, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("index must be an integer")
		}

		client, conn := client.NewClient(socketPath)
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		r, err := client.RemoveFolder(ctx, &pb.RemoveFolderRequest{FolderIndex: int32(idx)})
		if err != nil {
			log.Fatalf("could not remove folder: %v", err)
		}

		for i, folder := range r.GetFolders() {
			fmt.Printf("%d: %s\n", i, folder)
		}
	},
}

func init() {
	folderCmd.AddCommand(listFolderCmd)
	folderCmd.AddCommand(addFolderCmd)
	folderCmd.AddCommand(removeFolderCmd)

	rootCmd.AddCommand(folderCmd)
}
