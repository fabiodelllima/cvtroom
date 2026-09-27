# Changelog

All notable changes to this project are documented in this file. The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.1] - 2026-09-26

### Changed

- Revised the documentation for grammar, cohesion, and consistency.
- Upgraded the CI image to Go 1.26.

### Fixed

- Pinned govulncheck to v1.8.0 so that the security stage no longer depends on the latest release, which requires Go 1.26.

## [0.1.0] - 2026-09-26

### Added

- The `cvtr` command-line interface with `--url`, `--start`, `--end`, `--output`, `--dry-run`, and `--version`.
- Timestamp parsing, YouTube URL validation, and construction of the yt-dlp command.
- A GitLab CI pipeline with formatting checks, static analysis, tests under the race detector, and vulnerability scanning.

[0.1.1]: https://gitlab.com/delimafabio/cvtroom/-/compare/v0.1.0...v0.1.1
[0.1.0]: https://gitlab.com/delimafabio/cvtroom/-/tags/v0.1.0
