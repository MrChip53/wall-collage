package cmd

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"

	"wall-collage/client"
	"wall-collage/pb"
)

var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Sets a random image from the folders as wallpaper",
	Long:  "Sets a random image from the folders as wallpaper",
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := client.NewClient(socketPath)
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err := client.Random(ctx, &pb.RandomRequest{})
		if err != nil {
			log.Fatalf("could not start collage service: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(randomCmd)
}
