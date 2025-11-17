package main

import (
	"github.com/spf13/cobra"

	"github.com/koobie777/ark-android-forge/internal/android"
)

var (
	syncManifest string
	syncForce    bool
	syncDryRun   bool
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Run repo sync for the configured workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := android.SyncOptions{
			Manifest: syncManifest,
			Force:    syncForce,
			DryRun:   syncDryRun,
		}
		return android.RepoSync(cmd.Context(), appCtx.runner, appCtx.cfg, opts)
	},
}

func init() {
	syncCmd.Flags().StringVar(&syncManifest, "manifest", "", "custom manifest name to sync")
	syncCmd.Flags().BoolVar(&syncForce, "force", false, "force sync (repo --force-sync)")
	syncCmd.Flags().BoolVar(&syncDryRun, "dry-run", false, "log the command without executing it")
	rootCmd.AddCommand(syncCmd)
}
