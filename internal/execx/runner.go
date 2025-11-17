package execx

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*[A-Za-z]`)

// Command represents an external command invocation.
type Command struct {
	Name    string
	Args    []string
	Dir     string
	Env     map[string]string
	Timeout time.Duration
	DryRun  bool
}

// Runner executes commands with logging, timeouts, and retries.
type Runner struct {
	logger         zerolog.Logger
	defaultTimeout time.Duration
}

// NewRunner returns a configured Runner.
func NewRunner(logger zerolog.Logger) *Runner {
	return &Runner{
		logger:         logger,
		defaultTimeout: 2 * time.Hour,
	}
}

// Run executes the provided command until completion.
func (r *Runner) Run(ctx context.Context, cmd Command) error {
	if cmd.DryRun {
		r.logger.Info().
			Str("cmd", cmd.Name).
			Strs("args", cmd.Args).
			Msg("dry-run: skip execution")
		return nil
	}

	timeout := cmd.Timeout
	if timeout <= 0 {
		timeout = r.defaultTimeout
	}

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	execCmd := exec.CommandContext(ctx, cmd.Name, cmd.Args...)
	execCmd.Dir = cmd.Dir
	execCmd.Env = os.Environ()
	for k, v := range cmd.Env {
		execCmd.Env = append(execCmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	stdoutLog := newLogWriter(r.logger, zerolog.InfoLevel)
	stderrLog := newLogWriter(r.logger, zerolog.ErrorLevel)
	stdout := io.MultiWriter(&stdoutBuf, stdoutLog)
	stderr := io.MultiWriter(&stderrBuf, stderrLog)

	execCmd.Stdout = stdout
	execCmd.Stderr = stderr

	start := time.Now()
	r.logger.Info().
		Str("cmd", execCmd.String()).
		Msg("execx: starting command")

	if err := execCmd.Start(); err != nil {
		return fmt.Errorf("start command: %w", err)
	}

	err := execCmd.Wait()
	stdoutLog.Flush()
	stderrLog.Flush()
	duration := time.Since(start)

	event := r.logger.Info()
	if err != nil {
		event = r.logger.Error().Err(err)
	}
	event.
		Dur("duration", duration).
		Str("cmd", execCmd.String()).
		Msg("execx: command finished")

	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("command timeout after %s", timeout)
	}

	return err
}

type logWriter struct {
	logger zerolog.Logger
	level  zerolog.Level
	buf    bytes.Buffer
}

func newLogWriter(logger zerolog.Logger, level zerolog.Level) *logWriter {
	return &logWriter{
		logger: logger,
		level:  level,
	}
}

func (w *logWriter) Write(p []byte) (int, error) {
	total := len(p)
	for len(p) > 0 {
		idx := bytes.IndexByte(p, '\n')
		if idx == -1 {
			w.buf.Write(p)
			break
		}
		w.buf.Write(p[:idx])
		w.emit()
		p = p[idx+1:]
	}
	return total, nil
}

func (w *logWriter) Flush() {
	if w.buf.Len() == 0 {
		return
	}
	w.emit()
}

func (w *logWriter) emit() {
	line := strings.TrimSpace(stripANSI(w.buf.String()))
	if line != "" {
		w.logger.WithLevel(w.level).Msg(line)
	}
	w.buf.Reset()
}

func stripANSI(line string) string {
	return ansiRegexp.ReplaceAllString(line, "")
}
