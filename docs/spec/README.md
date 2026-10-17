# Specification

This directory holds the normative specification of the Atron language,
versioned alongside the implementation.

**Status: two chapters ratified, nothing implemented.** `lexical.md` and
`grammar.md` are the ratified design contract for M1 and M2 respectively; no
lexer, parser or semantic analyzer exists, and ratification is not an
authorisation to start M1. The remaining chapters are unwritten.

The baseline is the *Atron Master Implementation Guide* v1.0 (17 October 2026),
which is **not** mirrored here: `.gitignore` excludes the PDF from version
control, so it cannot be amended. Where a ratified decision differs from the
guide, the difference is recorded in this specification and in an ADR under
[`../decisions/`](../decisions/), and the specification wins.

The specification is written in English (ADR-0002) and is split by chapter of
the guide:

| File | Content | Milestone | State |
|---|---|---|---|
| [`lexical.md`](lexical.md) | source encoding, tokens, line endings, identifiers, literals, comments, ASI, lexical diagnostics | M1 | **ratified, not implemented** |
| [`grammar.md`](grammar.md) | grammar, declarations, precedence, block layout, entry point, contextual restrictions | M2 | **ratified, not implemented** |
| `types.md` | value types, equality, truthiness, conversions | M4 | planned |
| `semantics.md` | scoping, closures, modules, evaluation order | M3/M4/M7 | planned |
| `stdlib.md` | native modules and their contracts | M7 | planned |
| `errors.md` | diagnostic codes and message catalogue | M6 | planned |
| `bytecode.md` | instruction set and chunk format | M4 | planned |

Ratified decisions that deviate from the baseline guide:

| Decision | ADR | Language-change record |
|---|---|---|
| Immutable declarations use `const`, not `let` | [ADR-0004](../decisions/0004-const-declarations.md) | [#3](https://github.com/CosmicLang/atron-lang/issues/3) |
| Executable roots declare an explicit `fn main()` | [ADR-0005](../decisions/0005-explicit-entrypoint.md) | [#3](https://github.com/CosmicLang/atron-lang/issues/3) |

Rules:

- The specification states **observable behaviour**. Anything a user could write
  a test for belongs here; anything else belongs in an ADR or in package
  documentation.
- A behaviour change lands in the same pull request as the specification change
  that describes it. Code that contradicts the specification is a bug; the
  specification is amended deliberately, never to match accidental behaviour.
- A ratified chapter states plainly what is **not** implemented. A specification
  is a contract for future work, not evidence that the work was done.
- Proposals for language changes follow the process in
  [`../../CONTRIBUTING.md`](../../CONTRIBUTING.md) and are recorded as ADRs under
  [`../decisions/`](../decisions/).
