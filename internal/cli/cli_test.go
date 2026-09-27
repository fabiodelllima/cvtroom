package cli

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os/exec"
	"strings"
	"testing"
)

type fakeRunner struct {
	missing map[string]bool
	runErr  error
	called  bool
	gotArgs []string
}

func (f *fakeRunner) LookPath(name string) (string, error) {
	if f.missing[name] {
		return "", exec.ErrNotFound
	}
	return "/usr/bin/" + name, nil
}

func (f *fakeRunner) Run(_ context.Context, _ string, args []string, _, _ io.Writer) error {
	f.called, f.gotArgs = true, args
	return f.runErr
}

const videoURL = "https://www.youtube.com/watch?v=dQw4w9WgXcQ"

func run(t *testing.T, runner *fakeRunner, args ...string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), args, &stdout, &stderr, runner)
	return code, stdout.String(), stderr.String()
}

func TestRunSuccess(t *testing.T) {
	runner := &fakeRunner{}
	code, _, stderr := run(t, runner, "--url", videoURL, "--start", "1:05", "--end", "1:30")
	if code != ExitOK {
		t.Fatalf("code = %d, stderr = %q", code, stderr)
	}
	if !runner.called || !strings.Contains(strings.Join(runner.gotArgs, " "), "*65-90") {
		t.Errorf("runner args = %v", runner.gotArgs)
	}
}

func TestRunDryRunDoesNotExecute(t *testing.T) {
	runner := &fakeRunner{}
	code, stdout, _ := run(t, runner, "--url", videoURL, "--start", "10", "--end", "20", "--dry-run")
	if code != ExitOK || runner.called {
		t.Fatalf("code = %d, runner called = %v", code, runner.called)
	}
	want := "yt-dlp --download-sections *10-20 --force-keyframes-at-cuts --no-playlist -- " + videoURL + "\n"
	if stdout != want {
		t.Errorf("stdout = %q, want %q", stdout, want)
	}
}

func TestRunVersion(t *testing.T) {
	code, stdout, _ := run(t, &fakeRunner{}, "--version")
	if code != ExitOK || stdout != "cvtr dev\n" {
		t.Errorf("code = %d, stdout = %q", code, stdout)
	}
}

func TestRunHelp(t *testing.T) {
	code, _, stderr := run(t, &fakeRunner{}, "-h")
	if code != ExitOK || !strings.Contains(stderr, "-url") {
		t.Errorf("code = %d, stderr = %q", code, stderr)
	}
}

func TestRunUsageErrors(t *testing.T) {
	cases := map[string][]string{
		"missing flags":  {"--url", videoURL},
		"unknown flag":   {"--bogus"},
		"positional arg": {"--url", videoURL, "--start", "1", "--end", "2", "extra"},
		"invalid url":    {"--url", "https://vimeo.com/1", "--start", "1", "--end", "2"},
		"invalid start":  {"--url", videoURL, "--start", "abc", "--end", "2"},
		"invalid end":    {"--url", videoURL, "--start", "1", "--end", "1:60"},
		"reversed range": {"--url", videoURL, "--start", "20", "--end", "10"},
	}
	for name, args := range cases {
		runner := &fakeRunner{}
		code, _, stderr := run(t, runner, args...)
		if code != ExitUsage {
			t.Errorf("%s: code = %d, want %d (stderr %q)", name, code, ExitUsage, stderr)
		}
		if runner.called {
			t.Errorf("%s: runner must not be called", name)
		}
	}
}

func TestRunRuntimeFailures(t *testing.T) {
	cases := map[string]*fakeRunner{
		"missing yt-dlp": {missing: map[string]bool{"yt-dlp": true}},
		"yt-dlp fails":   {runErr: errors.New("exit status 1")},
	}
	for name, runner := range cases {
		code, _, stderr := run(t, runner, "--url", videoURL, "--start", "1", "--end", "2")
		if code != ExitFailure || !strings.HasPrefix(stderr, "error:") {
			t.Errorf("%s: code = %d, stderr = %q", name, code, stderr)
		}
	}
}
