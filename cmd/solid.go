package cmd

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"

	"wall-collage/client"
	"wall-collage/pb"
)

var solidCmd = &cobra.Command{
	Use:   "solid",
	Short: "Sets a solid color as wallpaper",
	Long:  "Sets a solid color as wallpaper",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatalf("solid command requires a color argument")
		}

		client, conn := client.NewClient(socketPath)
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err := client.SolidColor(ctx, &pb.SolidColorRequest{Color: args[0]})
		if err != nil {
			log.Fatalf("could not start collage service: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(solidCmd)
}
