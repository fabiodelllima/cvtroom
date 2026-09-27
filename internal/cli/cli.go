// Package cli parses the cvtr arguments and maps failures to exit codes,
// keeping main.go limited to wiring the program to the operating system.
package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"gitlab.com/delimafabio/cvtroom/internal/clip"
	"gitlab.com/delimafabio/cvtroom/internal/source"
	"gitlab.com/delimafabio/cvtroom/internal/timecode"
)

// Exit codes: 2 follows the Unix convention for usage errors.
const (
	ExitOK      = 0
	ExitFailure = 1
	ExitUsage   = 2
)

// Version is overridden at build time via -ldflags.
var Version = "dev"

type options struct {
	url, start, end, output string
	dryRun, version         bool
}

// Run executes the command and returns the exit code.
func Run(ctx context.Context, args []string, stdout, stderr io.Writer, runner clip.Runner) int {
	opts, code, done := parseFlags(args, stderr)
	if done {
		return code
	}
	if opts.version {
		fmt.Fprintf(stdout, "cvtr %s\n", Version)
		return ExitOK
	}

	req, err := buildRequest(opts)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return ExitUsage
	}
	if opts.dryRun {
		fmt.Fprintln(stdout, "yt-dlp "+strings.Join(req.Args(), " "))
		return ExitOK
	}
	if err := clip.Execute(ctx, runner, req, stdout, stderr); err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return ExitFailure
	}
	return ExitOK
}

func parseFlags(args []string, stderr io.Writer) (options, int, bool) {
	var opts options
	fs := flag.NewFlagSet("cvtr", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.StringVar(&opts.url, "url", "", "YouTube video URL (required)")
	fs.StringVar(&opts.start, "start", "", "clip start as SS, MM:SS or HH:MM:SS (required)")
	fs.StringVar(&opts.end, "end", "", "clip end as SS, MM:SS or HH:MM:SS (required)")
	fs.StringVar(&opts.output, "output", "", "yt-dlp output template (optional)")
	fs.BoolVar(&opts.dryRun, "dry-run", false, "print the yt-dlp command without running it")
	fs.BoolVar(&opts.version, "version", false, "print the version and exit")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return opts, ExitOK, true
		}
		return opts, ExitUsage, true
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(stderr, "error: unexpected arguments: %s\n", strings.Join(fs.Args(), " "))
		return opts, ExitUsage, true
	}
	return opts, ExitOK, false
}

func buildRequest(opts options) (clip.Request, error) {
	if opts.url == "" || opts.start == "" || opts.end == "" {
		return clip.Request{}, errors.New("--url, --start and --end are required")
	}
	videoURL, err := source.ValidateYouTube(opts.url)
	if err != nil {
		return clip.Request{}, err
	}
	start, err := timecode.Parse(opts.start)
	if err != nil {
		return clip.Request{}, fmt.Errorf("--start: %w", err)
	}
	end, err := timecode.Parse(opts.end)
	if err != nil {
		return clip.Request{}, fmt.Errorf("--end: %w", err)
	}
	req := clip.Request{URL: videoURL, Start: start, End: end, Output: opts.output}
	if err := req.Validate(); err != nil {
		return clip.Request{}, err
	}
	return req, nil
}
