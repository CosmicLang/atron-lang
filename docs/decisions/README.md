# Architecture Decision Records

This directory holds the decision records for Atron. An ADR explains **why** a
decision of lasting consequence was made, what else was considered, and what it
costs. It is not documentation of how the code works — that belongs in the
package documentation and in `docs/spec/`.

## Index

| ADR | Title | Status |
|---|---|---|
| [0001](0001-value-representation.md) | Runtime value representation | Accepted (benchmark validation deferred to M4) |
| [0002](0002-language-of-record.md) | Language of code, messages and documentation | Accepted |
| [0003](0003-license-name-extension.md) | License, project name and source file extension | Accepted for development; not cleared for public distribution |
| [0004](0004-const-declarations.md) | Immutable declarations use `const` | Accepted (2026-10-17) |
| [0005](0005-explicit-entrypoint.md) | Executable roots declare an explicit `fn main()` | Accepted (2026-10-17) |

## When an ADR is required

- language syntax or semantics;
- the bytecode instruction set or the chunk format;
- package boundaries, and the public API of `pkg/atron`;
- the choice of a significant dependency;
- anything a future contributor would otherwise have to re-litigate.

Bug fixes, behaviour-preserving refactors, tests, documentation and CI changes
do not need one.

## Format

Every ADR uses the same shape, taken from Appendix B.4 of the guide:

```markdown
# ADR-NNNN: Short title

Status: proposed | accepted | superseded by ADR-MMMM — Date: YYYY-MM-DD

## Context
The problem and the constraints.

## Decision
What was decided, in one or two sentences.

## Alternatives
The options considered and why they were rejected.

## Consequences
Positive, negative, and impact on tests, performance or compatibility.
```

## Rules

- **Numbering is sequential and never reused.** The next ADR is
  `0006-short-title.md`.
- **English.** See ADR-0002.
- **Never delete an ADR.** A decision that is reversed is marked
  `superseded by ADR-MMMM` and links to its replacement. The history of a wrong
  decision is worth more than a tidy directory.
- **One decision per record.** If a document needs two headings that could each
  stand alone, it is two ADRs.
- **State what is unverified.** If a decision depends on a measurement, a legal
  review or a search that has not been done, say so in the record and say who
  must do it. ADR-0001 and ADR-0003 both do this, and both are read as
  provisional for that reason.
- **Record reality, not intention.** An ADR whose consequences describe work
  that has not happened yet must label the section *Follow-up* rather than
  implying completion.
