# 0002. Delegate downloading to yt-dlp

## Status
Accepted

## Context
YouTube changes its delivery formats often, and reimplementing stream extraction would turn a small tool into a permanent maintenance burden.

## Decision
We delegate downloading to yt-dlp through `--download-sections`, which fetches only the requested range, and rely on ffmpeg for keyframe-accurate cuts. cvtr owns validation, argument construction and error reporting; yt-dlp owns extraction.

## Consequences
cvtr stays small and benefits from yt-dlp updates for free. Users must install yt-dlp and ffmpeg, so the CLI checks both before running. Process execution sits behind a `Runner` interface so tests never touch the network.
