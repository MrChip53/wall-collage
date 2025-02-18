package cmd

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"

	"wall-collage/client"
	"wall-collage/pb"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the wall-collage slideshow",
	Long:  "Stop the wall-collage slideshow",
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := client.NewClient(socketPath)
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err := client.Stop(ctx, &pb.StopRequest{})
		if err != nil {
			log.Fatalf("could not start collage service: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
