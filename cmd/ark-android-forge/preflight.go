package main

import (
	"github.com/spf13/cobra"

	"github.com/koobie777/ark-android-forge/internal/preflight"
)

var preflightCmd = &cobra.Command{
	Use:   "preflight",
	Short: "Run host readiness checks",
	RunE: func(cmd *cobra.Command, args []string) error {
		return preflight.Run(cmd.Context(), appCtx.cfg, appCtx.logger)
	},
}

func init() {
	rootCmd.AddCommand(preflightCmd)
}
