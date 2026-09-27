// Command cvtr clips a time range from a YouTube video without downloading the full file.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/delimafabio/cvtroom/internal/cli"
	"gitlab.com/delimafabio/cvtroom/internal/clip"
)

func main() {
	// Ctrl+C cancels the context, which stops yt-dlp instead of leaving it orphaned.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	code := cli.Run(ctx, os.Args[1:], os.Stdout, os.Stderr, clip.ExecRunner{})
	stop()
	os.Exit(code)
}
