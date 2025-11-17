package ui

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
	"unicode"

	"github.com/rs/zerolog"

	"github.com/koobie777/ark-android-forge/internal/config"
)

// RunMenu renders the interactive ARKFORGE menu that previously lived in bash.
func RunMenu(ctx context.Context, cfg *config.Config, logger zerolog.Logger) error {
	reader := bufio.NewReader(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		renderMenu(cfg)
		fmt.Print("Select ARKFORGE operation: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			logger.Info().Msg("Smart Build module is being ported to Go; stay tuned.")
		case "2":
			logger.Warn().Msg("Recovery Build is still in development.")
		case "3":
			logger.Info().Msg("ROM Build module is being ported to Go; stay tuned.")
		case "4":
			logger.Info().Msg("Boot/Recovery builder will be available after module rewrite.")
		case "5":
			logger.Info().Msg("Resume build is not yet implemented in Go.")
		case "6":
			logger.Info().Msg("Repo sync only flow is available via 'sync' command.")
		case "7":
			logger.Info().Msg("Device manager is coming to the Go orchestrator.")
		case "8":
			logger.Info().Msg("Repository manager rewrite in-progress.")
		case "9":
			logger.Info().Msg("Directory manager rewrite in-progress.")
		case "10":
			logger.Info().Msg("Configuration manager will arrive alongside the new API.")
		case "11":
			showFleetStatus(cfg)
		case "12":
			showUserGuide()
		case "13":
			logger.Info().Msg("Tmux manager is not required in Go mode; use 'tmux' CLI directly.")
		case "0":
			showExitMessage(cfg)
			return nil
		default:
			fmt.Printf("Unknown selection %q\n", input)
		}

		fmt.Print("Press Enter to return to the ARK Command Deck...")
		_, _ = reader.ReadString('\n')
	}
}

func renderMenu(cfg *config.Config) {
	clearScreen()
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	fmt.Println("===============================================")
	fmt.Println("             ARKFORGE COMMAND DECK             ")
	fmt.Println("===============================================")
	fmt.Printf("Commander: %s | Mode: %s | Jobs: %d\n", cfg.Commander, titleCase(cfg.Mode), cfg.Jobs)
	fmt.Printf("Time: %s UTC | Version: %s\n", now, cfg.Version)
	if inTmux() {
		fmt.Println("TMUX: active session detected")
	} else {
		fmt.Println("TMUX: no session detected (launch via 'tmux new -s arkforge')")
	}
	fmt.Printf("Workspace: %s\n", cfg.Build.Workspace)
	fmt.Println()
	fmt.Println("Primary Operations:")
	fmt.Println("  1) Smart Build      - Device discovery + repo selection + build")
	fmt.Println("  2) Recovery Build   - Build TWRP/OrangeFox recovery")
	fmt.Println("  3) ROM Build        - Compile a full ROM for a device")
	fmt.Println("  4) Boot/Recovery Images - Build boot/recovery from ROM source")
	fmt.Println("  5) Resume Build     - Continue interrupted builds")
	fmt.Println("  6) Repo Sync Only   - Sync repositories without building")
	fmt.Println()
	fmt.Println("Modules & System:")
	fmt.Println("  7) Device Manager         - Manage device database/configs")
	fmt.Println("  8) Repository Manager     - Manage ROM sources")
	fmt.Println("  9) Directory Manager      - Manage build/cache/output directories")
	fmt.Println("  10) Configuration Manager - ARK/Forge settings and customization")
	fmt.Println()
	fmt.Println("Fleet & Documentation:")
	fmt.Println("  11) Show Fleet Status  - List all ARK Fleet devices")
	fmt.Println("  12) User Guide         - Read the ARKFORGE user guide")
	fmt.Println("  13) Tmux Manager       - Manage ARK tmux sessions")
	fmt.Println()
	fmt.Println("  0) Exit ARKFORGE")
	fmt.Println()
}

func showFleetStatus(cfg *config.Config) {
	clearScreen()
	fmt.Println("============== ARK FLEET STATUS ==============")
	for _, device := range cfg.Fleet {
		fmt.Printf("- %-15s (%s) [%s]\n", device.Name, device.Codename, titleCase(device.Role))
		if device.Repository != "" {
			fmt.Printf("    Repository: %s\n", device.Repository)
		}
	}
	fmt.Println("==============================================")
}

func showUserGuide() {
	clearScreen()
	fmt.Println("============== USER GUIDE ====================")
	fmt.Println("The Go rewrite consolidates the bash modules into:")
	fmt.Println("- 'preflight' : Validate host prerequisites (Java, repo, ulimit, disk)")
	fmt.Println("- 'sync'      : Manage repo sync operations with manifests")
	fmt.Println("- 'build'     : Launch modular builds per device (coming soon)")
	fmt.Println("- 'release'   : Package artifacts and manifests (coming soon)")
	fmt.Println("Refer to README.md for detailed setup instructions.")
	fmt.Println("==============================================")
}

func showExitMessage(cfg *config.Config) {
	clearScreen()
	fmt.Println("===============================================")
	fmt.Printf(" Commander %s, the ARK stands ready for return.\n", cfg.Commander)
	fmt.Println("===============================================")
}

func clearScreen() {
	if runtime.GOOS == "windows" {
		fmt.Print("\033[H\033[2J")
	} else {
		fmt.Print("\033[H\033[2J")
	}
}

func inTmux() bool {
	_, ok := os.LookupEnv("TMUX")
	return ok
}

func titleCase(input string) string {
	if input == "" {
		return ""
	}
	runes := []rune(strings.ToLower(input))
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
