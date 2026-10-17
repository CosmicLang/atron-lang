# Changelog

All notable changes to the Atron project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Before 1.0, minor releases may contain breaking changes; the sections below state
whether a release is compatible with the previous one.

## [Unreleased]

Milestone M0 is complete but not yet approved, so its changes remain unreleased.

### Added

- Go module `github.com/CosmicLang/atron-lang`, declaring Go 1.23.0.
- `atron` command with `version` (build metadata injected through `-ldflags`)
  and `help` subcommands, `-h`/`--help` and `-v`/`--version` shorthands, usage
  text, and exit code 64 for incorrect command-line usage as fixed by chapter
  4.8 of the specification (`0` success, `1` runtime error, `2` compilation
  diagnostics, `64` CLI misuse).
- Package scaffolding with SPDX headers and package documentation for
  `internal/token`, `internal/lexer`, `internal/ast`, `internal/parser`,
  `internal/sema`, `internal/object`, `internal/bytecode`, `internal/compiler`,
  `internal/vm`, `internal/stdlib`, `internal/diag`, `internal/module`,
  `internal/repl`, `internal/format` and `pkg/atron`.
- `Makefile` with `build`, `test`, `lint`, `fmt`, `fmt-check`, `vet`, `cover`,
  `bench`, `tidy`, `verify`, `clean` and `help` targets; `verify` is the full
  local quality gate and now runs `tidy` first.
- `golangci-lint` configuration for v2 (schema version 2) enabling `govet`,
  `staticcheck`, `errcheck`, `ineffassign`, `unused`, `revive`, `gocritic`,
  `gosec` and `misspell`, with `gofmt` and `goimports` as formatters.
- CI workflow running gofmt, `go mod tidy` verification, `go vet`, race tests
  with coverage, a metadata-injected build, a CLI smoke test, `golangci-lint`,
  `govulncheck` and a DCO sign-off check on pull requests, across Linux, macOS
  and Windows on Go 1.23.x and the current stable Go.
- Repository documentation: `README.md`, `CONTRIBUTING.md`,
  `CODE_OF_CONDUCT.md`, `SECURITY.md`, `GOVERNANCE.md`, `NOTICE`, this
  changelog, and architecture decision records under `docs/decisions/`.
- `docs/PROGRESS.md`, the milestone progress tracker.
- Apache License 2.0 as the project license.

### Changed

- `SECURITY.md` names GitHub Private Vulnerability Reporting (PVR) as the
  project's only private vulnerability channel and documents the submission flow
  through GitHub Security Advisories. PVR **is enabled** on the repository
  (`{"enabled": true}`, verified 2026-10-17), so the documented flow is the
  operational one. The code-of-conduct address is deliberately not offered as a
  security contact.
- `CODE_OF_CONDUCT.md` publishes the enforcement contact, replacing the previous
  placeholder, and adds explicit confidentiality commitments and a
  conflict-of-interest path. The Contributor Covenant 2.1 text and the
  enforcement ladder are unchanged.
- Build metadata is injected into the `main` package symbols
  (`main.version`, `main.commit`, `main.date`). This is the only symbol name the
  Go linker accepts for a `package main` executable; any documentation showing a
  module-path `-X` target was corrected during M0.

### Fixed

- `atron` returned exit code `2` for incorrect command-line usage. Chapter 4.8
  of the specification reserves `2` for compilation diagnostics and assigns `64`
  to CLI misuse, so `atron bogus`, `atron` with no arguments and
  `atron version extra` now exit `64`. A wrapper script can once again tell a
  mis-typed command line apart from a failed compile.
- `atron version` now reports the metadata injected by the build instead of
  always printing `devel`/`unknown`. The initial implementation targeted a
  module-path linker symbol that the Go linker does not resolve for
  `package main`, so the injected values were silently ignored.
- `make tidy` reported success for an untidy module whenever `go.mod` was not
  yet tracked, because it compared `git diff` output and Git ignores untracked
  files. It now uses `go mod tidy -diff`, which reads the module files directly
  and fails with the required diff, and `verify` depends on it.

### Security

- No execution commands exist yet; the command line does not read `.atrn`
  files, the network or the environment.

<!--
Template for a release. Copy it above the previous release, replacing
"Unreleased" with the version and the release date.

## [0.1.0] - YYYY-MM-DD

### Added

### Changed

### Deprecated

### Removed

### Fixed

### Security
-->

[Unreleased]: https://github.com/CosmicLang/atron-lang/commits/main
