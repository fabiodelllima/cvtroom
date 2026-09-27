# 0001. Rewrite the CLI in Go

## Status
Accepted

## Context
The proof of concept validated the idea of clipping only a range of a video. The tool must be easy to install on any machine and its behaviour must be covered by fast, deterministic tests.

## Decision
We rewrite the CLI in Go. A single static binary removes the need for an interpreter or virtual environment, the standard library covers flag parsing, URL handling and process execution, and `go test -race` gives fast feedback without extra tooling.

## Consequences
Distribution becomes a single file per platform. The project depends only on the standard library, which keeps the supply-chain surface minimal. Contributors need a Go toolchain to build from source.
