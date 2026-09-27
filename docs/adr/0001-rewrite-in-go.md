# 0001. Rewrite the CLI in Go

## Status

Accepted

## Context

The proof of concept validated the approach of downloading only a range of a video. The next version must be easy to install on any machine, and its behavior must be covered by fast, deterministic tests.

## Decision

We rewrite the CLI in Go. A single static binary removes the need for an interpreter or a virtual environment; the standard library covers flag parsing, URL handling, and process execution; and `go test -race` provides fast feedback without additional tooling.

## Consequences

Distribution is reduced to a single file per platform. Because the project depends only on the standard library, its supply-chain surface remains minimal. In exchange, contributors need a Go toolchain to build from source.
