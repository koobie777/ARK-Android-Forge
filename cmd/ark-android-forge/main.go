package main

import (
	"os"

	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	if err := rootCmd.Execute(); err != nil {
		logger.Error().Err(err).Msg("ark-android-forge execution failed")
		os.Exit(1)
	}
}
