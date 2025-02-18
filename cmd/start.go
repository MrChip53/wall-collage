package cmd

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"

	"wall-collage/client"
	"wall-collage/pb"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the wall-collage slideshow",
	Long:  "Start the wall-collage slideshow",
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := client.NewClient(socketPath)
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err := client.Start(ctx, &pb.StartRequest{})
		if err != nil {
			log.Fatalf("could not start collage service: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
