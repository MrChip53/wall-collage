package cmd

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"

	"wall-collage/client"
	"wall-collage/pb"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the configuration of the wall-collage daemon",
	Long:  "Set the configuration of the wall-collage daemon",
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := client.NewClient(socketPath)
		defer conn.Close()

		// Optional delay Flags
		if delay, err := cmd.Flags().GetInt("delay"); err == nil && delay > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, err := client.SetDelay(ctx, &pb.SetDelayRequest{Delay: int32(delay)})
			if err != nil {
				log.Fatalf("could not set delay: %v", err)
			}
		}

		// Optional background color Flags
		if bgColor, err := cmd.Flags().GetString("bg-color"); err == nil && bgColor != "" {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, err := client.SetBackgroundColor(ctx, &pb.SetBackgroundColorRequest{Color: bgColor})
			if err != nil {
				log.Fatalf("could not set background color: %v", err)
			}
		}
	},
}

func init() {
	setCmd.Flags().Int("delay", 0, "Delay between images in seconds")
	setCmd.Flags().String("bg-color", "", "Background color")
	setCmd.Flags().String("folder-path", "", "Folder path to scan for images")

	rootCmd.AddCommand(setCmd)
}
