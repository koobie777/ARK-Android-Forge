package android

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/koobie777/ark-android-forge/internal/config"
	"github.com/koobie777/ark-android-forge/internal/execx"
)

// BuildOptions describe how to launch a build.
type BuildOptions struct {
	Device       string
	Target       string
	Variant      string
	RepoOverride string
	DryRun       bool
}

// Build runs envsetup + lunch + m/mka for the requested device.
func Build(ctx context.Context, runner *execx.Runner, cfg *config.Config, opts BuildOptions) error {
	if runner == nil {
		return fmt.Errorf("runner is nil")
	}
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	if opts.Device == "" {
		if len(cfg.Fleet) > 0 {
			opts.Device = cfg.Fleet[0].Codename
		} else {
			return fmt.Errorf("device required")
		}
	}
	if opts.Target == "" {
		opts.Target = cfg.Build.DefaultType
	}
	if opts.Variant == "" {
		opts.Variant = "userdebug"
	}

	repoName := opts.RepoOverride
	if repoName == "" {
		if device := cfg.DeviceByCodename(opts.Device); device != nil && device.Repository != "" {
			repoName = device.Repository
		} else {
			repoName = "android"
		}
	}

	sourceDir := filepath.Join(cfg.Build.Workspace, fmt.Sprintf("%s-%s", repoName, opts.Device))
	envsetup := filepath.Join(sourceDir, "build", "envsetup.sh")
	if _, err := os.Stat(envsetup); err != nil {
		return fmt.Errorf("envsetup missing in %s: %w", sourceDir, err)
	}

	script := fmt.Sprintf("set -euo pipefail; source build/envsetup.sh && lunch %s-%s && m %s -j%d",
		opts.Device, opts.Variant, opts.Target, cfg.Jobs)

	cmd := execx.Command{
		Name:   "bash",
		Args:   []string{"-lc", script},
		Dir:    sourceDir,
		DryRun: opts.DryRun,
		Env: map[string]string{
			"ARK_COMMANDER": cfg.Commander,
		},
	}

	return runner.Run(ctx, cmd)
}
