---
name: 🗣️ Language change proposal
about: Propose a change to Atron syntax, semantics or the diagnostics contract
title: "[LANGUAGE]: "
labels: enhancement
assignees: ''
---

<!--
  Use this template for any proposal that changes what Atron source code means:
  new or removed syntax, a new keyword, a changed diagnostic code, a change to
  evaluation or module semantics.

  Read CONTRIBUTING.md → "Proposing a language change" before filing. A language
  change needs an ADR and a normative specification update; this issue is the
  public record that those documents link to.
-->

## Summary

One or two sentences describing the change.

## Motivation

What problem does this solve? Who does it affect, and what is the cost of doing
nothing?

## Proposed change

Be concrete enough that someone can implement it. Show before/after source where
possible.

```atron
// before

```

```atron
// after

```

## Conflicts with the baseline

The normative baseline is the *Atron Master Implementation Guide* v1.0, which is
**not** in version control (`.gitignore` excludes the PDF), so deviations have to
be recorded in this repository instead of in the guide.

- [ ] This change is compatible with the guide.
- [ ] This change **conflicts** with the guide, and the conflicting provision is
      quoted below together with its section and page.

> Guide §___, page ___, currently says:

## Alternatives considered

Other designs that were considered and why they were rejected.

## Compatibility and migration

- Is any released binary, tag or published module affected? (Before 1.0 the
  answer is normally "no".)
- Does existing source in the wild need to change? If so, what is the migration
  path, and is a diagnostic with a hint required?

## Milestones affected

Which milestone implements this, and which later milestones depend on it?

- [ ] M1 lexer
- [ ] M2 AST and parser
- [ ] M3 resolver
- [ ] M4 values, compiler, VM
- [ ] M5 compound types and structs
- [ ] M6 errors and diagnostics
- [ ] M7 modules, standard library, embedding API
- [ ] M8 formatter, test runner, REPL
- [ ] Documentation only

## Chronology disclosure

Be explicit about how this proposal originated. This matters: a decision that
was already made by the maintainer must not be described as the outcome of
community consensus.

- [ ] This is a new proposal open for discussion before any decision.
- [ ] The decision has **already been made** by the maintainer; this issue is a
      public record opened for audit and objection, not a retrospective vote.

## Governance checklist

- [ ] I have read `CONTRIBUTING.md`.
- [ ] I understand that a language change requires an ADR and a normative
      specification update in `docs/spec/`.
- [ ] I will link the ADR(s) and the specification pull request here once they
      exist.
