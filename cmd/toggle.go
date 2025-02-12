package cmd

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"

	"wall-collage/client"
	"wall-collage/pb"
)

var toggleCmd = &cobra.Command{
	Use:   "toggle",
	Short: "Toggle settings",
}

var toggleCollageCmd = &cobra.Command{
	Use:   "collage",
	Short: "Toggle collages",
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := client.NewClient(socketPath)
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err := client.ToggleCollage(ctx, &pb.ToggleCollageRequest{})
		if err != nil {
			log.Fatalf("could not toggle collages: %v", err)
		}
	},
}

func init() {
	toggleCmd.AddCommand(toggleCollageCmd)

	rootCmd.AddCommand(toggleCmd)
}
