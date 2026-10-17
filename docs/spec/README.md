# Specification

This directory will hold the normative specification of the Atron language,
versioned alongside the implementation.

**Status: empty at M0.** The single source of truth today is the *Atron Master
Implementation Guide v1.0* (17 October 2026), which is not mirrored here.

The specification is written in English (ADR-0002) and is split by chapter of
the guide:

| File (planned) | Content | Milestone |
|---|---|---|
| `lexical.md` | source encoding, tokens, terminators, literals, comments | M1 |
| `grammar.md` | grammar, precedence, statements and expressions | M2 |
| `types.md` | value types, equality, truthiness, conversions | M4 |
| `semantics.md` | scoping, closures, modules, evaluation order | M3/M4/M7 |
| `stdlib.md` | native modules and their contracts | M7 |
| `errors.md` | diagnostic codes and message catalogue | M6 |
| `bytecode.md` | instruction set and chunk format | M4 |

Rules:

- The specification states **observable behaviour**. Anything a user could write
  a test for belongs here; anything else belongs in an ADR or in package
  documentation.
- A behaviour change lands in the same pull request as the specification change
  that describes it. Code that contradicts the specification is a bug; the
  specification is amended deliberately, never to match accidental behaviour.
- Proposals for language changes follow the process in
  [`../../CONTRIBUTING.md`](../../CONTRIBUTING.md) and are recorded as ADRs under
  [`../decisions/`](../decisions/).
