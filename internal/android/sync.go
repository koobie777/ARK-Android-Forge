package android

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/koobie777/ark-android-forge/internal/config"
	"github.com/koobie777/ark-android-forge/internal/execx"
)

// SyncOptions configures repo sync runs.
type SyncOptions struct {
	Manifest string
	Force    bool
	DryRun   bool
}

// RepoSync performs a repo sync within the configured workspace.
func RepoSync(ctx context.Context, runner *execx.Runner, cfg *config.Config, opts SyncOptions) error {
	if runner == nil {
		return fmt.Errorf("runner is nil")
	}
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	if err := os.MkdirAll(cfg.Build.Workspace, 0o755); err != nil {
		return fmt.Errorf("prepare workspace: %w", err)
	}

	args := []string{"sync", "--current-branch", fmt.Sprintf("--jobs=%d", cfg.Jobs)}
	if opts.Manifest != "" {
		manifest := filepath.Base(opts.Manifest)
		args = append(args, fmt.Sprintf("--manifest-name=%s", manifest))
	}
	if opts.Force {
		args = append(args, "--force-sync")
	}

	cmd := execx.Command{
		Name:   "repo",
		Args:   args,
		Dir:    cfg.Build.Workspace,
		DryRun: opts.DryRun,
		Env: map[string]string{
			"ARK_COMMANDER": cfg.Commander,
		},
	}

	return runner.Run(ctx, cmd)
}
