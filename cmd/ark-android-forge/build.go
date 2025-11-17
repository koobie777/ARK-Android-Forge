package main

import (
	"github.com/spf13/cobra"

	"github.com/koobie777/ark-android-forge/internal/android"
)

var (
	buildDevice  string
	buildTarget  string
	buildVariant string
	buildRepo    string
	buildDryRun  bool
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Run envsetup + lunch + m for a device",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := android.BuildOptions{
			Device:       buildDevice,
			Target:       buildTarget,
			Variant:      buildVariant,
			RepoOverride: buildRepo,
			DryRun:       buildDryRun,
		}
		return android.Build(cmd.Context(), appCtx.runner, appCtx.cfg, opts)
	},
}

func init() {
	buildCmd.Flags().StringVar(&buildDevice, "device", "", "device codename to build (defaults to fleet primary)")
	buildCmd.Flags().StringVar(&buildTarget, "target", "", "build target (recovery, bootimage, etc.)")
	buildCmd.Flags().StringVar(&buildVariant, "variant", "userdebug", "lunch variant (user, userdebug, eng)")
	buildCmd.Flags().StringVar(&buildRepo, "repo", "", "override repository directory inside workspace")
	buildCmd.Flags().BoolVar(&buildDryRun, "dry-run", false, "log command without running it")
	rootCmd.AddCommand(buildCmd)
}
