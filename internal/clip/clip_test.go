package clip

import (
	"context"
	"errors"
	"io"
	"os/exec"
	"reflect"
	"testing"
	"time"
)

type fakeRunner struct {
	missing map[string]bool
	runErr  error
	called  bool
	gotName string
	gotArgs []string
}

func (f *fakeRunner) LookPath(name string) (string, error) {
	if f.missing[name] {
		return "", exec.ErrNotFound
	}
	return "/usr/bin/" + name, nil
}

func (f *fakeRunner) Run(_ context.Context, name string, args []string, _, _ io.Writer) error {
	f.called, f.gotName, f.gotArgs = true, name, args
	return f.runErr
}

const videoURL = "https://www.youtube.com/watch?v=dQw4w9WgXcQ"

func validRequest() Request {
	return Request{URL: videoURL, Start: 10 * time.Second, End: 20500 * time.Millisecond}
}

func TestValidate(t *testing.T) {
	cases := []struct {
		name    string
		start   time.Duration
		end     time.Duration
		wantErr bool
	}{
		{"valid range", 0, time.Second, false},
		{"negative start", -time.Second, time.Second, true},
		{"end equals start", 5 * time.Second, 5 * time.Second, true},
		{"end before start", 10 * time.Second, 5 * time.Second, true},
	}
	for _, tc := range cases {
		err := Request{URL: videoURL, Start: tc.start, End: tc.end}.Validate()
		if tc.wantErr != errors.Is(err, ErrInvalidRange) {
			t.Errorf("%s: Validate() error = %v, wantErr %v", tc.name, err, tc.wantErr)
		}
	}
}

func TestArgs(t *testing.T) {
	req := validRequest()
	want := []string{"--download-sections", "*10-20.5", "--force-keyframes-at-cuts", "--no-playlist", "--", videoURL}
	if got := req.Args(); !reflect.DeepEqual(got, want) {
		t.Errorf("Args() = %v, want %v", got, want)
	}

	req.Output = "clip.%(ext)s"
	want = []string{"--download-sections", "*10-20.5", "--force-keyframes-at-cuts", "--no-playlist",
		"--output", "clip.%(ext)s", "--", videoURL}
	if got := req.Args(); !reflect.DeepEqual(got, want) {
		t.Errorf("Args() with output = %v, want %v", got, want)
	}
}

func TestExecuteSuccess(t *testing.T) {
	runner := &fakeRunner{}
	if err := Execute(context.Background(), runner, validRequest(), io.Discard, io.Discard); err != nil {
		t.Fatalf("Execute() returned error: %v", err)
	}
	if runner.gotName != "yt-dlp" || !reflect.DeepEqual(runner.gotArgs, validRequest().Args()) {
		t.Errorf("runner received %s %v", runner.gotName, runner.gotArgs)
	}
}

func TestExecuteRejectsInvalidRangeWithoutRunning(t *testing.T) {
	runner := &fakeRunner{}
	req := validRequest()
	req.End = req.Start
	if err := Execute(context.Background(), runner, req, io.Discard, io.Discard); !errors.Is(err, ErrInvalidRange) {
		t.Errorf("Execute() error = %v, want ErrInvalidRange", err)
	}
	if runner.called {
		t.Error("runner must not be called for an invalid range")
	}
}

func TestExecuteReportsMissingDependency(t *testing.T) {
	for _, dep := range Dependencies {
		runner := &fakeRunner{missing: map[string]bool{dep: true}}
		err := Execute(context.Background(), runner, validRequest(), io.Discard, io.Discard)
		if !errors.Is(err, ErrMissingDependency) {
			t.Errorf("missing %s: error = %v, want ErrMissingDependency", dep, err)
		}
		if runner.called {
			t.Errorf("missing %s: runner must not be called", dep)
		}
	}
}

func TestExecuteWrapsRunnerFailure(t *testing.T) {
	boom := errors.New("exit status 1")
	runner := &fakeRunner{runErr: boom}
	if err := Execute(context.Background(), runner, validRequest(), io.Discard, io.Discard); !errors.Is(err, boom) {
		t.Errorf("Execute() error = %v, want wrapped %v", err, boom)
	}
}

// ExecRunner tests use the Go toolchain itself, which is present wherever the
// tests run, to exercise os/exec without depending on yt-dlp.
func TestExecRunner(t *testing.T) {
	r := ExecRunner{}
	if _, err := r.LookPath("go"); err != nil {
		t.Fatalf("LookPath(go) returned error: %v", err)
	}
	if _, err := r.LookPath("cvtr-nonexistent-binary"); err == nil {
		t.Error("LookPath of a nonexistent binary must fail")
	}
	if err := r.Run(context.Background(), "go", []string{"version"}, io.Discard, io.Discard); err != nil {
		t.Errorf("Run(go version) returned error: %v", err)
	}
}
