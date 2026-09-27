# Requirements

Each requirement follows the convention of ISO/IEC/IEEE 29148: it is stated with "shall" and can be verified individually.

## Functional requirements

- **FR-01:** The CLI shall clip a YouTube video given `--url`, `--start`, and `--end`.
- **FR-02:** The CLI shall download only the requested range, never the full video.
- **FR-03:** The CLI shall accept timestamps as `SS`, `MM:SS`, or `HH:MM:SS`, with optional fractional seconds.
- **FR-04:** The CLI shall accept video URLs on `youtube.com`, `m.youtube.com`, `music.youtube.com`, and `youtu.be`.
- **FR-05:** The CLI shall provide an `--output` flag that forwards an output template to yt-dlp.
- **FR-06:** The CLI shall provide a `--dry-run` flag that prints the underlying command without running it.
- **FR-07:** The CLI shall provide a `--version` flag that prints the version.

## Non-functional requirements

- **NFR-01:** The CLI shall validate all input before starting any external process.
- **NFR-02:** The CLI shall exit with code 0 on success, 1 on a runtime failure, and 2 on a usage error.
- **NFR-03:** The CLI shall fail fast, with a clear message, when yt-dlp or ffmpeg is missing.
- **NFR-04:** The project shall ship as a single static binary whose only runtime dependencies are yt-dlp and ffmpeg.
- **NFR-05:** The test suite shall keep coverage above 90% and shall run under the race detector in CI.
- **NFR-06:** The CLI shall never allow user input to be interpreted as a yt-dlp flag.
