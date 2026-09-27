.RECIPEPREFIX := >
BINARY  := bin/cvtr
MODULE  := gitlab.com/delimafabio/cvtroom
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X $(MODULE)/internal/cli.Version=$(VERSION)

.PHONY: build test cover fmt vet vuln install clean

build:
> CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/cvtr

test:
> go test -race ./...

cover:
> go test -race -covermode=atomic -coverprofile=coverage.out ./...
> go tool cover -func=coverage.out

fmt:
> gofmt -w .

vet:
> go vet ./...

vuln:
> go run golang.org/x/vuln/cmd/govulncheck@latest ./...

install:
> CGO_ENABLED=0 go install -trimpath -ldflags "$(LDFLAGS)" ./cmd/cvtr

clean:
> rm -rf bin coverage.out
