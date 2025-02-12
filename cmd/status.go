package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	"wall-collage/client"
	"wall-collage/pb"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of the wall-collage",
	Long:  "Get the status of the wall-collage",
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := client.NewClient(socketPath)
		defer conn.Close()

		_, err := client.Status(context.Background(), &pb.StatusRequest{})
		if err != nil {
			log.Fatalf("could not get status: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
