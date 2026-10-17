# SPDX-License-Identifier: Apache-2.0
#
# Makefile for the Atron project.
#
# Targets required by milestone M0 (chapter 7 of the Atron Master
# Implementation Guide): build, test, lint, fmt, cover, bench. The `verify`
# target runs the full local gate used before reporting a milestone.

# Repository identity ---------------------------------------------------------
MODULE      := github.com/CosmicLang/atron-lang
BINARY      := atron
CMD_PKG     := ./cmd/atron

# Build output ----------------------------------------------------------------
BIN_DIR     := bin
BIN         := $(BIN_DIR)/$(BINARY)
COVER_DIR   := coverage
COVER_FILE  := $(COVER_DIR)/coverage.out
COVER_HTML  := $(COVER_DIR)/coverage.html

# Tool versions ---------------------------------------------------------------
GOLANGCI_LINT_VERSION := v2.15.0

# Build metadata injected through -ldflags ------------------------------------
#   VERSION  overridable: make build VERSION=0.1.0
#   COMMIT   the current Git commit, or "unknown" outside a Git checkout
#   DATE     RFC 3339 UTC timestamp; honours SOURCE_DATE_EPOCH for reproducible
#            builds (https://reproducible-builds.org/specs/source-date-epoch/)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo devel)
COMMIT  ?= $(shell git rev-parse HEAD 2>/dev/null || echo unknown)
ifdef SOURCE_DATE_EPOCH
DATE    ?= $(shell date -u -d "@$(SOURCE_DATE_EPOCH)" +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || date -u -r "$(SOURCE_DATE_EPOCH)" +%Y-%m-%dT%H:%M:%SZ)
else
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
endif

# Build metadata is injected into the variables declared in
# cmd/atron/version.go. That file belongs to package main, so the linker symbol
# is "main.<name>" rather than the import path of the package: the linker
# resolves -X targets by symbol name, and package main is recorded as "main".
# If the metadata variables are ever moved into an importable package,
# VERSION_PKG must be updated to match.
VERSION_PKG := main
LDFLAGS     := -s -w \
	-X '$(VERSION_PKG).version=$(VERSION)' \
	-X '$(VERSION_PKG).commit=$(COMMIT)' \
	-X '$(VERSION_PKG).date=$(DATE)'

GOFLAGS_BUILD := -trimpath

GO      ?= go
GOBIN   := $(shell $(GO) env GOPATH)/bin
LINT    := $(shell command -v golangci-lint 2>/dev/null || echo $(GOBIN)/golangci-lint)

.DEFAULT_GOAL := build

.PHONY: all build test lint fmt fmt-check vet cover bench tidy verify clean help

all: build

## build: compile the atron binary into bin/ with injected build metadata
build:
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS_BUILD) -ldflags "$(LDFLAGS)" -o $(BIN) $(CMD_PKG)
	@echo "built $(BIN) ($(VERSION) @ $(COMMIT))"

## test: run the full test suite with the race detector
test:
	$(GO) test -race ./...

## lint: run golangci-lint with the repository configuration
lint:
	@if [ ! -x "$(LINT)" ]; then \
		echo "golangci-lint not found; install with:"; \
		echo "  go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)"; \
		exit 1; \
	fi
	$(LINT) run

## fmt: format all Go sources in place
fmt:
	$(GO) fmt ./...

## fmt-check: fail if any Go file is not gofmt-formatted
fmt-check:
	@out="$$(gofmt -l .)"; \
	if [ -n "$$out" ]; then \
		echo "the following files are not gofmt-formatted:"; \
		echo "$$out"; \
		exit 1; \
	fi
	@echo "gofmt: clean"

## vet: run go vet on all packages
vet:
	$(GO) vet ./...

## cover: run tests with coverage and write an HTML report
cover:
	@mkdir -p $(COVER_DIR)
	$(GO) test -race -covermode=atomic -coverprofile=$(COVER_FILE) ./...
	$(GO) tool cover -func=$(COVER_FILE) | tail -1
	$(GO) tool cover -html=$(COVER_FILE) -o $(COVER_HTML)
	@echo "coverage report: $(COVER_HTML)"

## bench: run all benchmarks
bench:
	$(GO) test -run '^$$' -bench . -benchmem ./...

## tidy: verify go.mod and go.sum are tidy without modifying the working tree
# `go mod tidy -diff` (Go 1.23+) prints the changes go.mod/go.sum would need and
# exits non-zero if there are any. It is used instead of `go mod tidy` followed
# by `git diff --quiet` because `git diff` silently ignores untracked files: in
# the initial repository state, where go.mod has never been committed, that
# check passed unconditionally and could not detect an untidy module.
tidy:
	@$(GO) mod tidy -diff >/dev/null || { \
		echo "go.mod/go.sum are not tidy; run 'go mod tidy' and commit the result:"; \
		$(GO) mod tidy -diff || true; \
		exit 1; \
	}
	@echo "go.mod: tidy"

## verify: run the complete local quality gate
verify: tidy fmt-check vet test lint build

## clean: remove build and coverage artifacts
clean:
	rm -rf $(BIN_DIR) $(COVER_DIR)

## help: list available targets
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /'
