package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "wall-collage",
	Short: "Wall collage will create a collage of images and set them as your wallpaper",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
