# Atron

Atron is an original, open-source programming language implemented in Go. It is
dynamically and strongly typed, its sources use the `.atrn` extension, and its
toolchain ships as a single `atron` binary.

**Status: milestone M0 — repository foundation.** There is no language
implementation yet. The repository currently contains the Go module, package
scaffolding, the `atron` command with `version` and `help`, tests, build tooling,
CI configuration and project documentation. Nothing in this repository can
compile or execute an `.atrn` program, and no such claim is made anywhere below.

---

## What Atron is (design intent)

Atron is designed as a small, predictable language:

- **Compiled, not interpreted ad hoc.** Source is compiled to bytecode and run
  by a stack virtual machine, so there is a single execution path.
- **Dynamically and strongly typed.** No implicit coercions; errors are
  reported instead of silently converting values.
- **Deterministic.** Program output must not depend on map iteration order or
  other incidental runtime details.
- **Embeddable.** The runtime is designed to hold no mutable global state, so
  several interpreters can coexist in one Go process.

These are design decisions recorded in the specification and the ADRs. None of
them is implemented yet beyond package scaffolding.

## Status of the language

| Capability | Status |
|---|---|
| `atron version` (build metadata via `-ldflags`) | Implemented and tested |
| `atron help`, usage text, invalid-command handling | Implemented and tested |
| Package layout, `doc.go` scaffolding, SPDX headers | Implemented |
| Build tooling, linting, CI, project docs, ADRs | Implemented |
| Lexer, tokens | Planned — milestone M1 |
| AST, parser | Planned — milestone M2 |
| Semantic analysis (resolver) | Planned — milestone M3 |
| Runtime values, bytecode, compiler, VM | Planned — milestone M4 |
| Compound types and structs | Planned — milestone M5 |
| Errors and diagnostics renderer | Planned — milestone M6 |
| Modules, standard library, embedding API | Planned — milestone M7 |
| Formatter, test runner, REPL | Planned — milestone M8 |
| Performance work | Planned — milestone M9 |
| 0.1.0 release | Planned — milestone M10 |

The authoritative source of truth for language syntax and semantics is the Atron
Master Implementation Guide v1.0 (17 October 2026). Its summary will be
published under `docs/spec/` as the language takes shape.

## Installation

At M0 there is no release and no published module version. The remote `main`
branch contains only the repository's initial commit, so the module is **not**
fetchable and `go install ...@main` fails today. The only supported way to obtain
a binary is to build from a local clone:

```bash
# 1. Build from a clone, with build metadata injected by the Makefile.
git clone https://github.com/CosmicLang/atron-lang.git
cd atron-lang
make build
./bin/atron version

# 2. Plain Go build without the Makefile (metadata shows "devel").
go build -o bin/atron ./cmd/atron
```

Once M0 is committed and pushed, the development version becomes installable:

```bash
go install github.com/CosmicLang/atron-lang/cmd/atron@main
```

Prebuilt binaries, checksums and package-manager formulae are planned for the
0.1.0 release (milestone M10) and are **not** available today.

## Usage

```console
$ atron version
atron version devel
  commit:   unknown
  built:    unknown
  go:       go1.23.0
  platform: darwin/arm64
```

```console
$ atron help
atron — the Atron programming language toolchain

Usage:
  atron <command> [arguments]

Commands:
  version   Print build information (version, commit, build date, Go version)
  help      Print this help text, or help for a single command
...
```

Exit codes follow chapter 4.8 of the Atron specification:

| Code | Meaning |
| ---- | ------- |
| `0`  | successful execution |
| `1`  | unrecoverable runtime error or panic (reserved; not produced at M0) |
| `2`  | compilation diagnostics (reserved; not produced at M0) |
| `64` | incorrect command-line usage (unknown command, unexpected argument) |

Only `0` and `64` can occur at M0; `1` and `2` belong to later milestones.

There is no `atron run` yet. The Hello World program below illustrates the
intended syntax of the language and **cannot be executed by this repository
today**:

```atron
println("Hello, world!")
```

## Building and testing

Requirements: Go 1.23 or newer, `git`, and `make` (optional but used by the
documented commands). The language core uses the Go standard library only.

```bash
make build     # build bin/atron with version/commit/date injected
make test      # go test -race ./...
make lint      # golangci-lint run
make fmt       # gofmt formatting in place (go fmt ./...)
make cover     # coverage profile + HTML report in coverage/
make bench     # benchmarks
make verify    # gofmt check, vet, race tests, lint, build
```

Without `make`:

```bash
gofmt -l .
go vet ./...
go test -race ./...
golangci-lint run
```

`golangci-lint` v2.15.0 is the pinned linter version:

```bash
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0
```

## Repository layout

```
cmd/atron/          command-line entry point (subcommand dispatch)
internal/token/     token kinds, positions, keywords          (M1)
internal/lexer/     lexical analysis and terminators          (M1)
internal/ast/       syntax tree and S-expression printer      (M2)
internal/parser/    recursive descent + Pratt parser          (M2)
internal/sema/      resolver: scopes, symbols, captures       (M3)
internal/object/    runtime values                            (M4)
internal/bytecode/  opcodes, chunks, disassembler             (M4)
internal/compiler/  AST to bytecode                           (M4)
internal/vm/        stack virtual machine                     (M4)
internal/stdlib/    native modules                            (M7)
internal/diag/      diagnostics model and renderer            (M1/M6)
internal/module/    module loading, cache and resolution      (M7)
internal/repl/      interactive shell                         (M8)
internal/format/    canonical formatter                       (M8)
pkg/atron/          public embedding API (the only public API)
testdata/           golden files and fixtures                (M1+)
examples/           runnable example programs                (M4+)
editors/            TextMate grammar and VS Code extension    (M10)
docs/               progress, ADRs, specification, benchmarks
```

`testdata/`, `examples/` and `editors/` hold only a `README.md` explaining what
will land there; they are part of the chapter 3.3 layout but contain no content
at M0.

Everything under `internal/` is private to the module. `pkg/atron` will be the
only public Go API and will follow strict semantic compatibility from 1.0.

## Project documentation

- [`docs/PROGRESS.md`](docs/PROGRESS.md) — current milestone, completed work,
  pending items and verification evidence.
- [`docs/decisions/`](docs/decisions/) — architecture decision records.
- [`CONTRIBUTING.md`](CONTRIBUTING.md) — build, test, commit and proposal
  process.
- [`GOVERNANCE.md`](GOVERNANCE.md) — decision making and trademark policy.
- [`SECURITY.md`](SECURITY.md) — private vulnerability reporting.
- [`CHANGELOG.md`](CHANGELOG.md) — release history (Unreleased section).

## Important notices

- **Name collision.** The command name `atron` is not unique: an unrelated
  Python tool (`atron-cli`) also installs an `atron` executable. Do not install
  Atron over an existing unrelated tool of the same name, and check `which
  atron` after installation. See ADR-0003.
- **Trademark and name clearance.** Public distribution requires trademark and
  name-availability checks that have **not** been performed. See ADR-0003 and
  `GOVERNANCE.md`.
- **No stability promise yet.** The language, the bytecode format and
  `pkg/atron` may change without notice before 1.0.

## License

Licensed under the Apache License, Version 2.0. See [`LICENSE`](LICENSE) for the
full text and [`NOTICE`](NOTICE) for attribution. By contributing you agree to
the Developer Certificate of Origin; see [`CONTRIBUTING.md`](CONTRIBUTING.md).
