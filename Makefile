.PHONY: build test vet fmt lint clean

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || echo unknown)

LDFLAGS := -s -w \
	-X 'github.com/leixiaotian1/ginGen/internal/version.Version=$(VERSION)' \
	-X 'github.com/leixiaotian1/ginGen/internal/version.Commit=$(COMMIT)' \
	-X 'github.com/leixiaotian1/ginGen/internal/version.BuildDate=$(DATE)'

build:
	go build -trimpath -ldflags "$(LDFLAGS)" -o bin/ginGen .

test:
	go test ./... -count=1

vet:
	go vet ./...

fmt:
	gofmt -w .

clean:
	rm -rf bin/
