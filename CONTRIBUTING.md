# Contributing to Atron

Thanks for your interest. This document describes how to build, test and propose
changes. At milestone M0 the project is small enough that these rules are short;
they will grow as the language does.

## Before you start

At the time of writing the repository has **no language implementation**: only
the `atron` command with `version` and `help` exists. Work on the lexer,
parser, VM and the rest is organised in milestones documented in
[`docs/PROGRESS.md`](docs/PROGRESS.md). Please open an issue before starting a
large change so effort is not duplicated.

## Requirements

- Go 1.23 or newer (the module declares `go 1.23.0`).
- `git` with `user.name` and `user.email` configured, because every commit must
  carry a `Signed-off-by` trailer.
- `make` (used by the documented commands; running the underlying `go` commands
  directly is equally valid).
- `golangci-lint` v2.15.0 for `make lint`.

```bash
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0
```

## Build and test

```bash
make build      # bin/atron with version/commit/date injected
make test       # go test -race ./...
make lint       # golangci-lint run
make cover      # coverage report in coverage/
make bench      # benchmarks
make verify     # gofmt check, vet, race tests, lint, build
```

`make verify` is the same gate CI enforces. A change is not ready for review
until it passes locally on your machine.

## Code conventions

- **Formatting.** `gofmt` only; the build fails on unformatted files. Imports
  are grouped by `goimports` (standard library, then external).
- **License headers.** Every Go file starts with
  `// SPDX-License-Identifier: Apache-2.0` followed by a blank line. Do not add
  the long Apache header.
- **Package documentation.** Each package has a `doc.go` describing its
  responsibility and the milestone that introduces it. Update it when the
  package's scope changes.
- **Exported identifiers** carry doc comments starting with the identifier name,
  as required by `revive`.
- **Error strings** are lowercase and carry no trailing punctuation. Diagnostics
  that reach the user are produced by the diagnostics package, not concatenated
  at the call site.
- **No mutable global state.** The runtime must be safe to instantiate more than
  once in a process; this is a design requirement, not a style preference.
- **Determinism.** Nothing may depend on Go map iteration order, pointer values
  or wall-clock time in observable output. Sort keys explicitly.
- **Linting.** Warnings are not disabled to make a build green. If a check must
  be suppressed, add a narrowly scoped `//nolint` comment **with a reason**.

## Tests

- Table-driven tests are the default for the front end and the runtime.
- Test names describe behaviour (`TestParseRejectsUnterminatedString`), not
  implementation details.
- Every bug fix comes with a regression test that fails before the fix.
- Golden files live under `testdata/` and are updated together with the code
  that produces them.
- Tests must not depend on the working directory, the network or the clock.

## Commit messages

Atron uses [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <summary in the imperative mood>

<body: what changed and why, wrapped at 72 columns>

Signed-off-by: Your Name <you@example.com>
```

Accepted types: `feat`, `fix`, `docs`, `test`, `refactor`, `perf`, `build`,
`ci`, `chore`. Scopes mirror the package layout: `cli`, `token`, `lexer`,
`ast`, `parser`, `sema`, `object`, `bytecode`, `compiler`, `vm`, `stdlib`,
`diag`, `module`, `repl`, `format`, `docs`, `ci`.

Rules:

- One logical change per commit; do not mix formatting-only changes with
  behaviour changes.
- The summary line is at most 72 characters and has no trailing period.
- Reference issues in the body (`Refs: #42`), not in the summary.
- Do not rewrite published history on `main`.

## Developer Certificate of Origin

Atron uses the DCO instead of a contributor license agreement. Every commit must
end with a `Signed-off-by` trailer matching the commit author:

```bash
git commit -s -m "feat(lexer): recognise identifiers"
```

The trailer certifies the Developer Certificate of Origin, version 1.1:

```
By making a contribution to this project, I certify that:

(a) The contribution was created in whole or in part by me and I have the right
    to submit it under the open source license indicated in the file; or
(b) The contribution is based upon previous work that, to the best of my
    knowledge, is covered under an appropriate open source license and I have
    the right under that license to submit that work with modifications,
    whether created in whole or in part by me, under the same open source
    license (unless I am permitted to submit under a different license), as
    indicated in the file; or
(c) The contribution was provided directly to me by some other person who
    certified (a), (b) or (c) and I have not modified it.
(d) I understand and agree that this project and the contribution are public and
    that a record of the contribution (including all personal information I
    submit with it, including my sign-off) is maintained indefinitely and may be
    redistributed consistent with this project or the open source license(s)
    involved.
```

## Pull requests

1. Fork the repository and create a topic branch from `main`.
2. Keep the change focused; unrelated clean-ups belong in separate PRs.
3. Make sure `make verify` passes.
4. Fill in the pull request template: what changed, why, and how it was tested.
5. Link the issue the change addresses, if any.
6. Expect review comments about tests, determinism and error messages. Those
   three areas are where most review feedback lands.

Changes that alter the language's observable behaviour require an ADR under
`docs/decisions/` **in the same pull request**, and must update the relevant
documentation. A behaviour change without a corresponding spec and CHANGELOG
entry will not be merged.

## Proposing language changes

Language changes are more consequential than code changes, because programs
written by other people depend on them. The process is:

1. Open an issue using the *Language change proposal* template and describe the
   problem, not only the proposed syntax.
2. Discussion happens in the issue. Consensus is sought, not voted on.
3. If accepted, the change is recorded in an ADR and the specification is
   updated before implementation.
4. Implementation follows the milestone plan in `docs/PROGRESS.md`.

Breaking changes before 1.0 are allowed when they buy a materially better
language, and are recorded in `CHANGELOG.md` under *Changed*. After 1.0 they
require a deprecation cycle.

## Reporting bugs

Use the bug report template. A useful report includes the exact `atron version`
output, the source that triggered the problem, the expected and actual
behaviour, and the platform. Security issues must **not** be filed publicly; see
[`SECURITY.md`](SECURITY.md).

## Code of conduct

Participation is governed by [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md).

## License of contributions

By contributing you agree that your contribution is licensed under the Apache
License, Version 2.0, and you grant the project the rights described in that
license. Do not contribute code you cannot license this way.
