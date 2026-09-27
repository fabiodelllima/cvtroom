# 0002. Delegate downloading to yt-dlp

## Status

Accepted

## Context

YouTube changes its delivery formats frequently, so reimplementing stream extraction would turn a small tool into a permanent maintenance burden.

## Decision

We delegate downloading to yt-dlp through its `--download-sections` option, which fetches only the requested range, and rely on ffmpeg for keyframe-accurate cuts. Responsibilities are divided accordingly: cvtr owns validation, argument construction, and error reporting, while yt-dlp owns extraction.

## Consequences

cvtr remains small and inherits improvements to yt-dlp at no cost. In exchange, users must install yt-dlp and ffmpeg, so the CLI checks for both before running. Process execution sits behind a `Runner` interface, which keeps the tests independent of the network.
