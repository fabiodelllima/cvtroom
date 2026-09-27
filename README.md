# CVTRoom

[![pipeline](https://gitlab.com/delimafabio/cvtroom/badges/main/pipeline.svg)](https://gitlab.com/delimafabio/cvtroom/-/pipelines)
[![coverage](https://gitlab.com/delimafabio/cvtroom/badges/main/coverage.svg)](https://gitlab.com/delimafabio/cvtroom/-/pipelines)

CVTRoom provides `cvtr`, a command-line tool that clips a time range from a YouTube video by downloading only that range rather than the full file.

## Tutorial

This tutorial takes you from installation to your first clip. First, install the two runtime dependencies, [yt-dlp](https://github.com/yt-dlp/yt-dlp) and [ffmpeg](https://ffmpeg.org/). Then build and install the binary, which requires Go 1.24 or later:

```bash
git clone https://gitlab.com/delimafabio/cvtroom.git
cd cvtroom
make install
```

Finally, clip the range between 1:05 and 1:30 of a video:

```bash
cvtr --url "https://youtu.be/dQw4w9WgXcQ" --start 1:05 --end 1:30
```

The clip is saved to the current directory under the title of the video.

## How-to guides

### Preview the command without downloading

Add `--dry-run` to print the yt-dlp command that cvtr would run, without executing it.

### Choose the output file name

Pass a yt-dlp output template through `--output`, for example `--output "intro.%(ext)s"`.

### Run the quality checks locally

`make test` runs the test suite under the race detector, and `make cover` reports test coverage. `make vet` performs static analysis, and `make vuln` scans the code and the standard library for known vulnerabilities.

## Reference

| Flag        | Required | Description                                                                                |
| ----------- | -------- | ------------------------------------------------------------------------------------------ |
| `--url`     | Yes      | Video URL on `youtube.com`, `m.youtube.com`, `music.youtube.com`, or `youtu.be`            |
| `--start`   | Yes      | Start of the range, as `SS`, `MM:SS`, or `HH:MM:SS`; seconds may include a fractional part |
| `--end`     | Yes      | End of the range, in the same formats; must be later than `--start`                        |
| `--output`  | No       | yt-dlp output template                                                                     |
| `--dry-run` | No       | Print the yt-dlp command without running it                                                |
| `--version` | No       | Print the version and exit                                                                 |

cvtr exits with code `0` on success, `1` on a runtime failure (a missing dependency or a yt-dlp error), and `2` on a usage error.

## Explanation

cvtr validates every input before it starts any external process and delegates extraction to yt-dlp, which downloads only the requested section and relies on ffmpeg for keyframe-accurate cuts. Because process execution sits behind a `Runner` interface, the entire CLI can be tested without network access. The project depends exclusively on the Go standard library.

Architectural decisions are recorded in [`docs/adr`](docs/adr), requirements in [`docs/requirements.md`](docs/requirements.md), and behavior scenarios in [`features`](features).

Use this tool only for content you are permitted to download, in accordance with the YouTube Terms of Service.
