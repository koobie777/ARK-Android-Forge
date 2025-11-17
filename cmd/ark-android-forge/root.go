package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/koobie777/ark-android-forge/internal/config"
	"github.com/koobie777/ark-android-forge/internal/execx"
	"github.com/koobie777/ark-android-forge/internal/ui"
)

type appState struct {
	cfg    *config.Config
	logger zerolog.Logger
	runner *execx.Runner
}

var (
	rootCmd = &cobra.Command{
		Use:   "ark-android-forge",
		Short: "Modular Android build orchestrator for The ARK Ecosystem",
		Long:  "ARKFORGE is the Go-based rewrite of the original shell orchestrator, providing modular Android build automation.",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			return bootstrapOnce()
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			if nonInteractive {
				cmd.Println("Non-interactive mode enabled; select an explicit sub-command.")
				return nil
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			return ui.RunMenu(ctx, appCtx.cfg, appCtx.logger)
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	appCtx = &appState{}

	configPath     string
	jsonLogs       bool
	nonInteractive bool
)

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "forge.yaml", "path to forge configuration file")
	rootCmd.PersistentFlags().BoolVar(&jsonLogs, "json", false, "enable structured JSON logging")
	rootCmd.PersistentFlags().BoolVar(&nonInteractive, "non-interactive", false, "disable interactive menu and require explicit sub-command")

	cobra.OnInitialize(func() {
		if (appCtx.logger == zerolog.Logger{}) {
			initLogger()
		}
	})
}

func bootstrapOnce() error {
	if (appCtx.logger == zerolog.Logger{}) {
		initLogger()
	}

	if appCtx.cfg == nil {
		cfg, err := config.Load(configPath)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		appCtx.cfg = cfg
	}

	if appCtx.runner == nil {
		appCtx.runner = execx.NewRunner(appCtx.logger)
	}

	return nil
}

func initLogger() {
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	if jsonLogs {
		appCtx.logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
		return
	}

	appCtx.logger = zerolog.New(output).With().Timestamp().Logger()
}
