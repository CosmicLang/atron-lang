# Changelog

All notable changes to the Atron project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Before 1.0, minor releases may contain breaking changes; the sections below state
whether a release is compatible with the previous one.

## [Unreleased]

Milestone M0 is **formally accepted** by the maintainer as technically complete
on 2026-10-17 (see `docs/audits/m0-acceptance.md`). Its changes remain unreleased:
no tag, version or binary has been published, and publishing is M10. M1 has not
started.

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
- The normative lexical and grammar specifications (`docs/spec/lexical.md`,
  `docs/spec/grammar.md`), ratified as the M1 and M2 design contract. **Nothing
  they describe is implemented.**
- A *Language change proposal* issue template
  (`.github/ISSUE_TEMPLATE/language_change.md`). `CONTRIBUTING.md` had required
  proposals to use this template, but it had never existed in the repository.

### Changed

- **Language (breaking, unreleased):** immutable declarations use `const`
  instead of `let`. `const` requires an initializer, `var` still defaults to
  `nil`, immutability applies to the binding rather than to the referenced
  object, and assignment to a `const` binding remains E0303. `let` is no longer a
  keyword and lexes as an ordinary identifier, so a legacy `let x = 1` is
  reported by the parser as E0201 with a migration hint. There is no `Let` token
  kind, `LetStmt` AST node or compatibility alias; the keyword table stays at 22
  entries and the token inventory at 62 kinds. See
  [ADR-0004](docs/decisions/0004-const-declarations.md) and issue
  [#3](https://github.com/CosmicLang/atron-lang/issues/3). No released surface is
  affected, because no release exists.
- **Language (breaking, unreleased):** an executable root must contain exactly
  one module-scope `fn main() { ... }` with zero parameters and a block body.
  `main` remains an ordinary identifier; `pub fn main()` is permitted. A
  parameterized or expression-bodied `fn main` is valid syntax and is rejected
  only by executable-context validation, with the new diagnostic **E0310**;
  duplicate declarations remain E0302 and E0201 is not reused. One grammar
  serves executable roots and libraries, and a separate contextual layer decides
  executability, so the REPL, tests, `Eval` and inspection commands require no
  `main`. Imported libraries are exempt and keep the guide's module semantics.
  The value returned by `main` is discarded and normal completion exits `0`;
  `exit(n)` remains the only way to set the status. See
  [ADR-0005](docs/decisions/0005-explicit-entrypoint.md) and issue
  [#3](https://github.com/CosmicLang/atron-lang/issues/3).
- **Language (clarification, unreleased):** LF, CRLF and a lone CR each count as
  exactly one logical line break; columns count Unicode code points and CRLF
  advances the byte offset by two bytes while incrementing the line once. A lone
  CR is therefore not an illegal character. The spec also fixes that the lexer
  never suppresses an ASI `Newline` because `{` follows: the *parser* may
  consume a single `Newline` before a required block-opening `{` in a recognised
  header, so both `fn main() {` and a brace on the next line are accepted, while
  `;` is never treated as optional whitespace.
- The `docs/spec/` and `docs/decisions/` README indexes list the ratified
  chapters, the accepted ADRs and their language-change records, and the next
  free ADR number is now `0006`.
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

- `docs/spec/lexical.md` contradicted itself about the dot. §4 listed `Dot` as one
  of the twelve punctuation kinds and `grammar.md` relied on it for member access,
  yet §8.2 stated that a bare `.` and a leading `.` in `.5` were lexical error
  E0103. Both cannot hold: a lexer that reports E0103 for `.` could not emit the
  `Dot` that `a.b` needs. The specification now states one rule — a `.` that does
  not continue a numeric literal is exactly one `Dot` token with **no** diagnostic,
  including `.5`, `5.` and a bare `.`; `.5` stays invalid as a float literal and is
  rejected as an expression by the parser with E0201; E0103 is reserved for
  malformed number literals and integer overflow. The M1 conformance checklist
  gained the corresponding cases and had two duplicated entry numbers repaired.
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
