package preflight

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"

	"github.com/koobie777/ark-android-forge/internal/config"
)

var ErrUnsupported = errors.New("unsupported")

// Result describes the outcome of a check.
type Result struct {
	Name     string
	Passed   bool
	Optional bool
	Details  string
	Err      error
}

// Run executes all preflight checks, mirroring the legacy bash workflow.
func Run(ctx context.Context, cfg *config.Config, logger zerolog.Logger) error {
	logger.Info().Msg("starting preflight checks")
	checks := []func(context.Context, *config.Config) Result{
		commandCheck("Java", "java", []string{"-version"}, false),
		commandCheck("Repo", "repo", []string{"--version"}, false),
		commandCheck("Git", "git", []string{"--version"}, false),
		commandCheck("ccache", "ccache", []string{"--version"}, true),
		workspaceCheck(),
		ulimitCheck(),
	}

	var failures []Result
	for _, checker := range checks {
		result := checker(ctx, cfg)
		printResult(result)
		if !result.Passed && !result.Optional {
			failures = append(failures, result)
		}
		if result.Err != nil {
			logger.Error().Err(result.Err).Str("check", result.Name).Msg("preflight check error")
		}
	}

	if len(failures) > 0 {
		return fmt.Errorf("preflight failed (%d checks)", len(failures))
	}

	logger.Info().Msg("preflight passed")
	return nil
}

func commandCheck(name, binary string, args []string, optional bool) func(context.Context, *config.Config) Result {
	return func(ctx context.Context, cfg *config.Config) Result {
		path, err := exec.LookPath(binary)
		if err != nil {
			details := fmt.Sprintf("%s command not found in PATH", binary)
			return Result{Name: name, Passed: false, Optional: optional, Details: details, Err: err}
		}
		return Result{Name: name, Passed: true, Optional: optional, Details: fmt.Sprintf("found at %s", path)}
	}
}

func workspaceCheck() func(context.Context, *config.Config) Result {
	return func(ctx context.Context, cfg *config.Config) Result {
		dir := cfg.Build.Workspace
		if dir == "" {
			return Result{Name: "Workspace", Passed: false, Details: "build.workspace not set"}
		}

		if err := os.MkdirAll(dir, 0o755); err != nil {
			return Result{Name: "Workspace", Passed: false, Details: "failed to create workspace", Err: err}
		}

		testFile := filepath.Join(dir, fmt.Sprintf(".arkforge-touch-%d", time.Now().UnixNano()))
		if err := os.WriteFile(testFile, []byte("ok"), 0o644); err != nil {
			return Result{Name: "Workspace", Passed: false, Details: "unable to write inside workspace", Err: err}
		}
		_ = os.Remove(testFile)

		return Result{Name: "Workspace", Passed: true, Details: fmt.Sprintf("ready at %s", dir)}
	}
}

func ulimitCheck() func(context.Context, *config.Config) Result {
	return func(ctx context.Context, cfg *config.Config) Result {
		limit, err := readUlimit()
		if err != nil {
			return Result{Name: "Ulimit", Passed: true, Optional: true, Details: "not supported on this platform"}
		}

		if limit < 8192 {
			return Result{Name: "Ulimit", Passed: false, Details: fmt.Sprintf("nofile=%d (<8192). Increase ulimit -n.", limit)}
		}

		return Result{Name: "Ulimit", Passed: true, Details: fmt.Sprintf("nofile=%d", limit)}
	}
}

func printResult(result Result) {
	status := "PASS"
	if !result.Passed {
		if result.Optional {
			status = "WARN"
		} else {
			status = "FAIL"
		}
	}

	fmt.Printf("[%s] %-12s - %s\n", status, result.Name, result.Details)
}
