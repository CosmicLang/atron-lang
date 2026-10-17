# ADR-0002: Language of code, messages and documentation

Status: accepted — Date: 2026-10-17

## Context

The maintainer's first language is Portuguese, and the Atron Master
Implementation Guide is written in Portuguese. Contributors and users, however,
are expected to be international: a programming language with Portuguese-only
diagnostics and documentation limits its audience severely, and tooling
ecosystems (editors, package indexes, translation services) assume English
source strings.

Three arrangements were considered:

1. Everything in Portuguese, including code comments, diagnostics and docs.
2. Everything in English.
3. English for code, diagnostics and primary documentation, with Portuguese
   documentation maintained in parallel under `docs/pt-BR/`.

The guide selects (3) and notes that the maintainer may switch to (1) by editing
this ADR at M0. That option was considered and rejected — see Alternatives.

## Decision

- **Identifiers, comments and commit messages:** English. Identifiers follow Go
  conventions; comments explain intent in English so that any contributor can
  review the code.
- **Diagnostics (user-visible messages):** English by default. Messages are
  identified by stable codes (`E0101`, `R0012`, …) and resolved through a
  message catalogue, so a translation never changes the code and a program's
  behaviour never depends on the language setting. `ATRON_LANG` selects a
  catalogue; `NO_COLOR` and `ATRON_LANG` are read from the environment only by
  the diagnostic renderer, never by the language core.
- **Primary documentation:** English. `README.md`, `CONTRIBUTING.md`,
  `SECURITY.md`, `GOVERNANCE.md`, ADRs, `CHANGELOG.md` and the specification are
  authoritative in English.
- **Portuguese documentation:** maintained under `docs/pt-BR/`, starting with a
  translation of the README and of the user-facing parts of the specification.
  Translations are explicitly non-authoritative: if they disagree with the
  English text, the English text wins.
- **Language identifiers:** `en` is the default and the only catalogue at M0. A
  missing translation falls back to `en` rather than failing.

## Alternatives considered

**Portuguese everywhere.** Rejected because it shrinks the contributor pool,
because editor and tooling conventions expect English identifiers, and because
translation of an entire diagnostics catalogue later would be a breaking change
in user-visible output. The guide allows it; the cost is not worth the
convenience for a single maintainer.

**English diagnostics, translated documentation only.** Rejected as a
half-measure: the most-read user-facing text for a compiler is its error
messages, so a Portuguese-speaking user would still meet English at the exact
moment comprehension matters.

**Runtime-switchable message strings inside the compiler.** Rejected: it would
put presentation inside the language core, break the "no mutable global state"
rule, and make golden tests depend on the environment.

## Consequences

Positive:

- One source of truth for terminology.
- Diagnostic codes are stable and language-independent, so tests and
  documentation can reference them forever.
- Translation is additive: adding a catalogue cannot break existing output.

Negative:

- The maintainer writes English documentation, which is more effort than
  Portuguese.
- `docs/pt-BR/` can lag behind and must be reviewed for staleness; a stale
  translation is worse than none. The directory therefore states its
  non-authoritative status explicitly.
- Adding a language later requires a review of every diagnostic for
  translatability (placeholders, pluralisation, word order). The catalogue is
  designed for this from M1, but the first translation will still expose
  messages that are hard to render.

## Follow-up

- M1: diagnostics catalogue with codes and English messages; `ATRON_LANG` is
  parsed but only `en` is available.
- M6: renderer applies the catalogue; unknown or missing languages fall back to
  `en` with no warning noise.
- Post-1.0: `docs/pt-BR/` translation review on each minor release.
