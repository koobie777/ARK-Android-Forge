package main

import (
	"github.com/spf13/cobra"

	"github.com/koobie777/ark-android-forge/internal/ui"
)

var menuCmd = &cobra.Command{
	Use:   "menu",
	Short: "Launch the interactive ARKFORGE command deck",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ui.RunMenu(cmd.Context(), appCtx.cfg, appCtx.logger)
	},
}

func init() {
	rootCmd.AddCommand(menuCmd)
}
