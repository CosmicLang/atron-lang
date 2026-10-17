# M0 acceptance record

Status: accepted — Date: 2026-10-17

This is a project-governance record, not an architecture decision. It states what
the maintainer accepted, the evidence that acceptance rests on, and what
acceptance does **not** authorize. Architecture decisions live in
[`docs/decisions/`](../decisions/).

## A. Acceptance declaration

The project maintainer formally declared **M0 — Repository Foundation accepted
as technically complete on 2026-10-17**.

This acceptance is a maintainer decision recorded in project documentation. It
is **not** a GitHub pull-request review, not a cryptographic signature, and not a
recorded approval event on any pull request: PR #1 was merged with zero submitted
reviews (`reviews: []`, `reviewDecision: ""`). No such artifact is claimed here
because none exists.

The declaration covers the technical milestone only. It does not establish
trademark clearance, release readiness, or permission to begin M1.

## B. Accepted implementation

The accepted work is milestone M0 as defined by chapter 7 of the Atron Master
Implementation Guide, delivered through **pull request #1** (`M0: establish Atron
repository foundation`) from `phase/m0-foundation` into `main`, and merged by the
maintainer.

Deliverables present in the accepted tree:

- `go.mod` declaring module `github.com/CosmicLang/atron-lang`, `go 1.23.0`.
- The chapter 3.3 directory structure, with a `doc.go` in each of the 15
  packages under `internal/` and `pkg/atron`; those packages contain **no
  executable code**, which is what fixes the M0/M1 boundary.
- `cmd/atron` implementing `version` and `help`.
- `Makefile` (`build`, `test`, `lint`, `fmt`, `cover`, `bench`, plus `fmt-check`,
  `vet`, `tidy`, `verify`, `clean`, `help`) and `.golangci.yml`.
- CI workflow per chapter 6.4: Ubuntu/macOS/Windows × Go 1.23.x and stable, plus
  lint, vulnerability check and DCO sign-off.
- `LICENSE` (Apache-2.0), `NOTICE`, `README.md`, `CONTRIBUTING.md`,
  `CODE_OF_CONDUCT.md`, `SECURITY.md`, `GOVERNANCE.md`, `CHANGELOG.md`.
- Issue templates (bug, feature, question), pull request template.
- ADR-0001, ADR-0002, ADR-0003 and this progress/acceptance documentation.

## C. Technical evidence

| Item | Value |
|---|---|
| Pull request | [#1](https://github.com/CosmicLang/atron-lang/pull/1) — **MERGED**, `mergedAt` 2026-10-17T22:45:04Z |
| PR head that was merged | `25c2b727aadb47d64763f4f54e267700fe5bf93d` |
| Pre-merge base | `40e612dc09bd6ee6541e54844c29d54c34844b98` |
| Merge commit | `78a19110b273d0623fb34a4fec293c54627963ee` |
| Merge strategy | merge commit (two parents: `40e612d`, `25c2b72`); not squash, not rebase |
| Post-merge workflow | [39076322979](https://github.com/CosmicLang/atron-lang/actions/runs/39076322979) — event `push`, **SUCCESS** |
| Final pre-merge workflow | [39076059769](https://github.com/CosmicLang/atron-lang/actions/runs/39076059769) — event `pull_request`, on `25c2b72`, nine of nine jobs |
| Test matrix | 6/6 — Ubuntu, macOS and Windows × Go 1.23.x and Go stable |
| Lint | `golangci-lint` — success, `0 issues` |
| Vulnerability check | success |
| DCO sign-off | success on the pull request; skipped on the post-merge `push` run by design |
| Go compatibility | module declares `go 1.23.0`; CI proves 1.23.x and the stable toolchain |

The post-merge run's `dco sign-off` job reports `skipped` rather than `success`
because its workflow condition is `github.event_name == 'pull_request'`. On a
`push` event it has zero steps and never runs. DCO was verified on PR #1 itself,
where every non-merge commit in the range carries a `Signed-off-by` trailer.


## D. Security status

The accepted `SECURITY.md` is the approved Atron policy, blob
`57d02660c984cb6eab55a6e8fe2e958a23dcaab1`, and it is what the merged `main`
carries. GitHub's generic starter template, added to `main` in `40e612d`, was
superseded rather than restored; the merged file contains none of that template's
text.

The policy documents coordinated disclosure through GitHub Private Vulnerability
Reporting, with intake at
`https://github.com/CosmicLang/atron-lang/security/advisories/new`, a disclosure
timing table and a fix-then-release procedure. It also keeps vulnerability
reporting separate from code-of-conduct moderation: the two "are handled through
different channels and are not interchangeable".

Private Vulnerability Reporting is **enabled**:
`GET /repos/CosmicLang/atron-lang/private-vulnerability-reporting` returns
`{"enabled": true}`. Activation was a repository-settings change made by the
maintainer, not a commit.

## E. Outstanding non-M0 obligations

These are outside M0's technical scope. None of them is satisfied by this
acceptance, and none should be read as waived.

- **Trademark and legal clearance are not completed.** No trademark search has
  been performed and no legal opinion obtained. The project remains *accepted for
  development, not cleared for public distribution* (ADR-0003).
- **The `atron` executable collides with the unrelated `atron-cli`.** Confirmed
  from that package's published metadata (`entry_points.txt` declares
  `atron=atron_cli.atron:cli`). Documented in `README.md` and `GOVERNANCE.md`.
- **Keep-versus-rename is still the maintainer's decision.** No rename was made
  during or after the merge; name, `.atrn` extension, executable name and module
  path are unchanged.
- **No public release is authorized.** No tag and no release exist; publishing is
  M10 scope.
- **Branch protection is not configured.** `main` has no protection rule and no
  required reviewers, so it currently accepts direct pushes and unreviewed
  merges. Hardening this is recommended as an administrative follow-up.

## F. Implementation boundary

**M0 acceptance does not authorize M1 implementation or a release.**

M1 (lexer) is not started and must not start on the strength of this record. No
M1 branch exists, and no lexer, token, parser, compiler, VM or other
implementation code has been written: every package under `internal/` and
`pkg/atron` still holds only a `doc.go`. M2 through M10 remain not started.

