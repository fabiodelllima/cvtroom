# CVTRoom

[![pipeline](https://gitlab.com/delimafabio/cvtroom/badges/main/pipeline.svg)](https://gitlab.com/delimafabio/cvtroom/-/pipelines)
[![coverage](https://gitlab.com/delimafabio/cvtroom/badges/main/coverage.svg)](https://gitlab.com/delimafabio/cvtroom/-/pipelines)

`cvtr` is a command-line tool that clips a time range from a YouTube video, downloading only that range instead of the full file.

## Tutorial

Install the runtime dependencies, [yt-dlp](https://github.com/yt-dlp/yt-dlp) and [ffmpeg](https://ffmpeg.org/), then build and install the binary (Go 1.24 or later):

```bash
git clone https://gitlab.com/delimafabio/cvtroom.git
cd cvtroom
make install
```

Clip seconds 65 to 90 of a video:

```bash
cvtr --url "https://youtu.be/dQw4w9WgXcQ" --start 1:05 --end 1:30
```

The clip is saved in the current directory, named after the video title.

## How-to guides

**Preview the command before downloading.** Add `--dry-run` to print the yt-dlp command without running it.

**Choose the output file name.** Pass a yt-dlp output template: `--output "intro.%(ext)s"`.

**Run the checks locally.** `make test` runs the suite with the race detector, `make cover` prints coverage, `make vet` and `make vuln` run static analysis and the vulnerability scan.

## Reference

| Flag        | Required | Description                                                                    |
| ----------- | -------- | ------------------------------------------------------------------------------ |
| `--url`     | yes      | Video URL on `youtube.com`, `m.youtube.com`, `music.youtube.com` or `youtu.be` |
| `--start`   | yes      | Clip start as `SS`, `MM:SS` or `HH:MM:SS`; seconds may be fractional           |
| `--end`     | yes      | Clip end, same formats; must be after `--start`                                |
| `--output`  | no       | yt-dlp output template                                                         |
| `--dry-run` | no       | Print the yt-dlp command and exit                                              |
| `--version` | no       | Print the version and exit                                                     |

Exit codes: `0` success, `1` runtime failure (missing dependency, yt-dlp error), `2` usage error.

## Explanation

cvtr validates every input before starting any process and delegates extraction to yt-dlp, which fetches only the requested section and uses ffmpeg for keyframe-accurate cuts. Process execution sits behind a `Runner` interface, so the whole CLI is tested without network access. The project depends only on the Go standard library.

Design decisions are recorded in [`docs/adr`](docs/adr), requirements in [`docs/requirements.md`](docs/requirements.md) and behaviour scenarios in [`features`](features).

Use this tool only for content you are allowed to download, in line with YouTube's Terms of Service.
