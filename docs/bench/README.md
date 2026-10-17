# Benchmarks

This directory holds published benchmark results for Atron.

**Status: empty at M0.** There is nothing to measure: the VM does not exist yet.
No performance number is stated anywhere in this repository, because any such
number would be fabricated.

## What lands here

| File (planned) | Content | Milestone |
|---|---|---|
| `0001-value-representation.md` | ADR-0001 validation: `any` vs 32-byte struct | M4 |
| `0002-vm-dispatch.md` | switch vs computed dispatch, if measured | M9 |
| `0003-baselines.md` | fib(30), sum loop, list append, map insert | M9 |
| `raw/` | raw `go test -bench` output and machine description | M9 |

## Rules

- Every result records the **machine and toolchain**: `go version`, `GOOS`,
  `GOARCH`, CPU model, and whether the binary was built with `-trimpath`.
- Raw output is stored verbatim next to the summary. A summary without raw data
  is not evidence.
- Benchmarks that inform a decision are cited from the corresponding ADR.
- Results are never edited to look better. A regression published honestly is
  more useful than a flattering number that nobody can reproduce.
- Comparative claims require both variants to be measured on the same machine in
  the same session.
