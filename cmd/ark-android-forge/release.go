package main

import (
	"github.com/spf13/cobra"

	"github.com/koobie777/ark-android-forge/internal/artifacts"
)

var releasePath string

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Generate release metadata (artifacts manifest)",
	RunE: func(cmd *cobra.Command, args []string) error {
		manifest := artifacts.Generate(appCtx.cfg)
		return artifacts.Write(releasePath, manifest)
	},
}

func init() {
	releaseCmd.Flags().StringVar(&releasePath, "output", "artifacts/manifest.yaml", "path to write manifest")
	rootCmd.AddCommand(releaseCmd)
}
