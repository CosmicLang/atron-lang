# Atron project progress

Current milestone: **M0 — repository foundation**
Last updated: 2026-10-17

| Milestone | State | Notes |
|---|---|---|
| M0 — repository foundation | complete, awaiting maintainer approval | Local gates green; exit codes corrected to §4.8; PDF excluded; published as draft PR #1; CI green on all six matrix legs against `d776938`. See “Verification”, “M0 acceptance criteria” and “Known gaps”. |
| M1 — lexer | not started | Must not begin before M0 is approved. |
| M2 — AST and parser | not started | |
| M3 — resolver | not started | |
| M4 — values, bytecode, compiler, VM | not started | Runs the ADR-0001 benchmark validation. |
| M5 — compound types and structs | not started | |
| M6 — errors and diagnostics | not started | |
| M7 — modules, standard library, embedding API | not started | |
| M8 — formatter, test runner, REPL | not started | |
| M9 — performance | not started | |
| M10 — 0.1.0 release | not started | |

`complete` means every deliverable of the milestone exists and every acceptance
criterion was demonstrated by a command whose output is recorded below. `not
started` means exactly that: no code for the milestone exists.

## Decisions in force

| ADR | Title | Status |
|---|---|---|
| [0001](decisions/0001-value-representation.md) | Runtime value representation | Accepted; benchmark validation deferred to M4 |
| [0002](decisions/0002-language-of-record.md) | Language of code, messages and documentation | Accepted |
| [0003](decisions/0003-license-name-extension.md) | License, project name and source file extension | Accepted for development, **not** cleared for public distribution |

## Governance decisions (M0-R2)

The maintainer settled the two governance placeholders that M0-R1 left open. No
ADR was created: neither decision changes the language, the bytecode format or
the public API of `pkg/atron`, so under `GOVERNANCE.md` they are reviewed here
rather than recorded as architecture.

| Question | Decision | Effect in the repository |
|---|---|---|
| Private vulnerability channel | GitHub Private Vulnerability Reporting (PVR) is the only accepted channel | `SECURITY.md` documents the intended PVR flow. PVR is **disabled** (`{"enabled": false}`); the document says so and claims nothing more. Enabling it needs the *Admin* role and is a pending prerequisite. |
| Code-of-conduct moderation contact | `joao.gabriel.vale.theodoro@gmail.com` | Published in the Enforcement section of `CODE_OF_CONDUCT.md`, with confidentiality commitments and a conflict-of-interest path. |

The two channels are kept deliberately separate. The moderation address is a
conduct channel and is not offered as a security contact, and `SECURITY.md`
states that explicitly, so a vulnerability report is never routed to an
address that was published for a different purpose.

The moderation address is a personal mailbox rather than a role address. That is
acceptable for a single-maintainer project, but it couples the code of conduct
to one person: if the maintainer changes, the address must be updated together
with this document and `GOVERNANCE.md`.

## Authoritative implementation reference

The normative specification for this repository is the **Atron Master
Implementation Guide** (*Atron · Guia mestre de implementação para IA*), v1.0,
dated 17 October 2026, 45 pages, written in Portuguese.

| Property | Value |
|---|---|
| Local path | `Atron_Guia_Mestre_Implementacao_IA.pdf` (repository root) |
| Size | 262,356 bytes |
| SHA-256 | `e04d8fce92d258d780cf45d575107acdbd0eb4cb76736412d5035a4050ed1f4d` |
| Version / date | v1.0 · 17 October 2026 |
| Version control | **Not versioned.** Excluded by the anchored root rule `/.gitignore` → `/Atron_Guia_Mestre_Implementacao_IA.pdf` |

The document is the source of truth for every milestone: chapter 4 defines the
language semantics (including the exit-code table of §4.8), chapter 3 the
repository layout, chapter 6 the CI and quality gates, and chapter 7 the
milestone plan. Where this repository and the guide disagree, the guide wins and
the repository is corrected — that is exactly how the exit-code defect in
“Exit-code contract” was found.

Maintainer decision (M0-R1): the PDF stays on the maintainer's machine and out of
version control, because it is a distributed reference rather than a source
artifact. The rule is root-anchored so it cannot accidentally match a same-named
file elsewhere in the tree. Any change to this decision is a maintainer decision
and must not be made silently: replacing the document requires updating the hash
above.

Verify the exclusion at any time with:

```bash
git check-ignore -v Atron_Guia_Mestre_Implementacao_IA.pdf
git status --porcelain | grep -i '\.pdf$'   # must print nothing
```

## M0 deliverables

Definition from chapter 7 of the Atron Master Implementation Guide, with status:

- [x] `go mod init` with the module path `github.com/CosmicLang/atron-lang`,
      `go 1.23.0`.
- [x] Directory structure of chapter 3.3 with a `doc.go` in every package.
- [x] `cmd/atron` with `version` and `help`.
- [x] `Makefile` with `build`, `test`, `lint`, `fmt`, `cover` and `bench` (plus
      `fmt-check`, `vet`, `tidy`, `verify`, `clean`, `help`).
- [x] `.golangci.yml` enabling the linters required by chapter 3.4.
- [x] CI workflow per chapter 6.4.
- [x] `LICENSE` with the official Apache-2.0 text.
- [x] `NOTICE`.
- [x] Initial `README.md`.
- [x] `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, `SECURITY.md`, `GOVERNANCE.md`,
      `CHANGELOG.md`.
- [x] Issue templates (bug, feature, question) and pull request template.
- [x] `docs/PROGRESS.md` (this file).
- [x] ADR-0001, ADR-0002, ADR-0003.

## M0 acceptance criteria

| Criterion | Status | Evidence |
|---|---|---|
| `make test` green | Met locally | `go test -race -count=1 ./...` → `ok github.com/CosmicLang/atron-lang/cmd/atron` |
| `make lint` green | Met locally | `golangci-lint run` → `0 issues.` |
| CI green on three platforms | **Met** | Run [39074104773](https://github.com/CosmicLang/atron-lang/actions/runs/39074104773) on `d776938`, triggered by pull request: `test` green on ubuntu-latest, macos-latest and windows-latest, each against Go 1.23.12 and Go 1.27.2 (six of six matrix legs). Each leg ran gofmt, `go mod tidy -diff`, vet, race tests with coverage (94.3%), a metadata-injected build and the exit-64 CLI smoke test; no step was skipped. `lint`, `vulnerability check` and `dco sign-off` also green in the same run. Two earlier runs failed on Windows and their defects are fixed and recorded; only the final run is cited as evidence. |
| `atron version` shows version and commit injected by `-ldflags` | Met | See “Verification” and `TestCompiledBinaryVersionMetadata`. |
| Every `.go` file has an SPDX header | Met | `find . -name '*.go' \| xargs grep -L 'SPDX-License-Identifier: Apache-2.0'` → no output (20 files). |
| `go vet` clean | Met | `go vet ./...` → no output |
| Process exit codes match chapter 4.8 | Met | `atron bogus` → 64, `atron` (no args) → 64, `atron version extra` → 64, `atron help` → 0, `atron version` → 0. Asserted both in-process (`TestExitCodeContract`, `TestAppRun`) and against the compiled executable (`TestCompiledBinaryExitCodes`). |

## Exit-code contract

Chapter 4.8 of the guide fixes four process exit codes. The implementation
declares only the two that M0 can produce; `1` and `2` are documented as
reserved so no reader mistakes an unimplemented code for a working one.

| Code | Meaning | Produced at M0 |
|---|---|---|
| `0` | successful execution | yes (`version`, `help`) |
| `1` | unrecoverable runtime error or panic | no — needs the VM (M4) |
| `2` | compilation diagnostics | no — needs the front end (M2–M6) |
| `64` | incorrect command-line usage | yes (unknown command, unexpected argument, no command) |

An earlier revision of this milestone used `2` for usage errors. That was a
specification violation: `2` belongs to compilation diagnostics, and a wrapper
script could not distinguish a mis-typed flag from a failed compile. The
correction is recorded in `CHANGELOG.md`. The expectations in the tests are
written as numeric literals rather than as the `exitUsage` constant, so a future
edit to the constant fails the suite instead of silently redefining the
contract.


## Remote verification (M0-R4)

The branch was published and a draft pull request opened. Two defects appeared
only on the Windows matrix legs and were fixed inside M0:

| Run | Attempt | Outcome |
|---|---|---|
| [39073735761](https://github.com/CosmicLang/atron-lang/actions/runs/39073735761) | `dd36dcf` | Windows failed at **Verify gofmt** on both Go versions; every later step was skipped. The repository declared no line-ending policy, so a checkout with the Git-for-Windows default `core.autocrlf=true` produced CRLF and `gofmt` rejected it. Fixed by adding `.gitattributes` with `text=auto eol=lf`. |
| [39073871074](https://github.com/CosmicLang/atron-lang/actions/runs/39073871074) | `d3835ba` | gofmt, tidiness and vet now passed on Windows. Both legs then failed in **Test with race detector** with `# .out: no required module provides package .out`: the step was the only one in the job without an explicit shell, so the Windows default `pwsh` re-tokenised the command and the coverage filename became a package argument. The tests themselves passed at 94.3%. Fixed by pinning `shell: bash`, as the neighbouring steps already did. |
| [39074104773](https://github.com/CosmicLang/atron-lang/actions/runs/39074104773) | `d776938` | **Green.** Six of six platform/Go legs, plus lint, vulnerability check and DCO. Every step executed; none skipped. |

Both defects were in the repository's configuration rather than in the Go
sources, and neither was worked around by weakening a check: no test was
removed, no lint rule relaxed and no job disabled.

## Verification

Commands run at the end of M0-R1, from the repository root, on `darwin/arm64`
with Go 1.27.2 (the installed toolchain; the module declares `go 1.23.0` and CI
tests both 1.23.x and stable). Output is reproduced verbatim.

```
$ gofmt -l .                       # no output: every file is formatted
$ go vet ./...                     # no output
$ go test -race -count=1 ./...
ok      github.com/CosmicLang/atron-lang/cmd/atron
$ golangci-lint run
0 issues.
$ make verify
go.mod: tidy
gofmt: clean
go vet ./...
go test -race ./...
ok      github.com/CosmicLang/atron-lang/cmd/atron
/Users/joaotheodoro/go/bin/golangci-lint run
0 issues.
go build -trimpath -ldflags "-s -w -X 'main.version=bea08fa-dirty' ..." -o bin/atron ./cmd/atron
built bin/atron (bea08fa-dirty @ bea08fa8ef493c6fd499701e96390632e38c38ed)
```

`verify` now depends on `tidy`, so a module that is not tidy fails the single
gate documented as complete rather than only the standalone target.

```
$ ./bin/atron version
atron version bea08fa-dirty
  commit:   bea08fa8ef493c6fd499701e96390632e38c38ed
  built:    2026-10-17T21:20:33Z
  go:       go1.27.2
  platform: darwin/arm64

$ ./bin/atron help >/dev/null;    echo $?   # 0
$ ./bin/atron version >/dev/null; echo $?   # 0
$ ./bin/atron bogus >/dev/null 2>&1;        echo $?   # 64
$ ./bin/atron >/dev/null 2>&1;              echo $?   # 64
$ ./bin/atron version extra >/dev/null 2>&1; echo $?  # 64
```

The commit hash above is the pre-existing initial commit; `-dirty` appears
because M0's files are uncommitted, which is the correct behaviour of
`git describe --dirty`.

`make tidy` was verified in both directions. It now runs `go mod tidy -diff`
(Go 1.23+), which reads the module files directly and never touches the working
tree:

```
$ make tidy                                    # tidy module
go.mod: tidy

$ printf '\nrequire github.com/google/uuid v1.6.0\n' >> go.mod && make tidy
go.mod/go.sum are not tidy; run 'go mod tidy' and commit the result:
diff current/go.mod tidy/go.mod
--- current/go.mod
+++ tidy/go.mod
@@ -1,5 +1,3 @@
 module github.com/CosmicLang/atron-lang
 
 go 1.23.0
-
-require github.com/google/uuid v1.6.0
make: *** [tidy] Error 1
```

The previous implementation ran `go mod tidy` and then `git diff --quiet -- go.mod
go.sum`. Because `git diff` ignores untracked files, it reported success
unconditionally in the initial repository state, where `go.mod` has never been
committed — the exact state the repository is in. That defect is fixed rather
than documented.

Test coverage is limited to `cmd/atron`, the only package containing code: `94.3%
of statements`, with every function at 100% except `main`, which is exercised by
the binary itself. 10 test functions and 34 cases (24 of them table subtests)
pass. Every other package holds a `doc.go` and no executable statements, so a
coverage number for them would be meaningless; the per-package coverage targets
of chapter 6.3 apply to the milestones that implement them.

## Known gaps and pending items

| Item | Owner | Consequence |
|---|---|---|
| CI has run and is green | **Resolved (M0-R4)** | Final run [39074104773](https://github.com/CosmicLang/atron-lang/actions/runs/39074104773) on `d776938`: six of six platform/Go legs green, plus lint, vulnerability check and DCO. The branch is published and draft PR #1 is open; `main` is still `bea08fa` until the maintainer merges. |
| M0 work published as a draft PR | **Resolved (M0-R4)** | `phase/m0-foundation` is pushed at `d776938` and PR #1 targets `main` in draft state. Nothing is merged; the merge decision is the maintainer's. |
| Git identity configured | **Resolved (M0-R3)** | The maintainer supplied the author e-mail. `user.name` and `user.email` are set at repository scope only; global and system scope were left untouched, and `git commit -s` derives every trailer from that identity. |
| The initial commit `bea08fa` uses a machine-local author address | Maintainer (optional) | It predates the confirmed identity and carries `joaotheodoro@Mac-mini-de-Joao.local`. It is pre-existing history on `main` and is deliberately not rewritten. The DCO job compares the PR range, so it is not flagged; rewriting it, if ever wanted, is a separate maintainer decision. |
| No moderation contact in `CODE_OF_CONDUCT.md` | **Resolved (M0-R2)** | The confirmed moderation address is published in the Enforcement section. Contributor Covenant 2.1 text, confidentiality language and escalation ladder are unchanged and consistent with `GOVERNANCE.md`. |
| Private security reporting channel not yet active | Maintainer | The maintainer selected GitHub Private Vulnerability Reporting (PVR) as the only accepted private channel, and PVR is **not enabled**: `GET /repos/CosmicLang/atron-lang/private-vulnerability-reporting` returns `{"enabled": false}`. `SECURITY.md` documents the intended process and states the disabled status plainly; it does not claim the channel is operational. Enabling PVR is a repository-settings action requiring the *Admin* role, and is outside the scope of this directive. |
| Name and extension clearance not performed (ADR-0003) | Maintainer | The project is not cleared for public distribution. |
| The `atron` executable collides with the unrelated `atron-cli` | Maintainer | Documented in `README.md`; may motivate a rename. |
| `docs/pt-BR/`, `docs/spec/`, `docs/bench/`, `testdata/`, `examples/` and `editors/` contain only a `README.md` describing what will land there | Milestones M1+ | Created so the layout of chapter 3.3 exists. Each directory explains its future contents and states that it is empty; no placeholder source files were invented. |
| No `release.yml` workflow and no `.goreleaser.yaml` | Milestone M10 | Publishing is out of scope for M0. |
| ~~No Git commit has been created~~ | **Resolved (M0-R3)** | Superseded: eight DCO-signed commits now exist on `phase/m0-foundation`. |

## Useful commands

```bash
make build      # binary in bin/ with metadata injected
make test       # race tests
make lint       # golangci-lint
make cover      # coverage profile and HTML report
make bench      # benchmarks (none exist yet)
make verify     # gofmt check, vet, race tests, lint, build
./bin/atron version
./bin/atron help
```

Fuzzing targets do not exist yet; the first is `FuzzLexer` in M1.

## Next milestone

M1 (lexer) is **not started** and must not start until the M0 report is
approved. The report is delivered in the session response rather than in this
file, per Appendix B.2 of the guide.
