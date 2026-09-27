# Requirements

## Functional

- FR-01: Clip a YouTube video given `--url`, `--start` and `--end`.
- FR-02: Download only the requested range, never the full video.
- FR-03: Accept timecodes as `SS`, `MM:SS` or `HH:MM:SS`, with optional fractional seconds.
- FR-04: Accept `youtube.com`, `m.youtube.com`, `music.youtube.com` and `youtu.be` URLs pointing to a video.
- FR-05: Offer `--output` to forward an output template to yt-dlp.
- FR-06: Offer `--dry-run` to print the underlying command without running it.
- FR-07: Offer `--version`.

## Non-functional

- NFR-01: Validate all input before starting any external process.
- NFR-02: Exit codes: 0 on success, 1 on runtime failure, 2 on usage error.
- NFR-03: Fail fast with a clear message when yt-dlp or ffmpeg is missing.
- NFR-04: Ship as a single static binary with no runtime beyond yt-dlp and ffmpeg.
- NFR-05: Keep test coverage above 90% and run tests with the race detector in CI.
- NFR-06: Never let user input be parsed as a yt-dlp flag.
