// Package clip describes a clip and delegates the download to yt-dlp.
package clip

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"time"

	"gitlab.com/delimafabio/cvtroom/internal/timecode"
)

var (
	// ErrInvalidRange reports a negative start or an end that does not follow the start.
	ErrInvalidRange = errors.New("invalid range")
	// ErrMissingDependency reports an external tool missing from PATH.
	ErrMissingDependency = errors.New("missing dependency")
)

// Dependencies lists the required external binaries: yt-dlp downloads the
// section and relies on ffmpeg to cut at the requested points.
var Dependencies = []string{"yt-dlp", "ffmpeg"}

// Request describes a clip whose source has already been validated.
type Request struct {
	URL    string
	Start  time.Duration
	End    time.Duration
	Output string
}

// Validate checks that the range is consistent.
func (r Request) Validate() error {
	if r.Start < 0 {
		return fmt.Errorf("%w: start (%s) must not be negative", ErrInvalidRange, r.Start)
	}
	if r.End <= r.Start {
		return fmt.Errorf("%w: end (%s) must be after start (%s)", ErrInvalidRange, r.End, r.Start)
	}
	return nil
}

// Args builds the yt-dlp command line. The "--" before the URL ends option
// parsing, so no input can ever be interpreted as a flag.
func (r Request) Args() []string {
	section := fmt.Sprintf("*%s-%s", timecode.Seconds(r.Start), timecode.Seconds(r.End))
	args := []string{
		"--download-sections", section,
		"--force-keyframes-at-cuts",
		"--no-playlist",
	}
	if r.Output != "" {
		args = append(args, "--output", r.Output)
	}
	return append(args, "--", r.URL)
}

// Runner abstracts process execution so tests run without network or yt-dlp.
type Runner interface {
	LookPath(name string) (string, error)
	Run(ctx context.Context, name string, args []string, stdout, stderr io.Writer) error
}

// ExecRunner is the real implementation, backed by os/exec.
type ExecRunner struct{}

// LookPath locates the binary in PATH.
func (ExecRunner) LookPath(name string) (string, error) {
	return exec.LookPath(name)
}

// Run executes the process; the context allows interruption with Ctrl+C.
func (ExecRunner) Run(ctx context.Context, name string, args []string, stdout, stderr io.Writer) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

// CheckDependencies fails early, with a clear message, if a binary is missing.
func CheckDependencies(r Runner) error {
	for _, name := range Dependencies {
		if _, err := r.LookPath(name); err != nil {
			return fmt.Errorf("%w: %s not found in PATH", ErrMissingDependency, name)
		}
	}
	return nil
}

// Execute validates the request, checks dependencies and downloads the clip.
func Execute(ctx context.Context, r Runner, req Request, stdout, stderr io.Writer) error {
	if err := req.Validate(); err != nil {
		return err
	}
	if err := CheckDependencies(r); err != nil {
		return err
	}
	if err := r.Run(ctx, "yt-dlp", req.Args(), stdout, stderr); err != nil {
		return fmt.Errorf("yt-dlp failed: %w", err)
	}
	return nil
}
