# ADR-0004: Immutable declarations use `const`

Status: accepted — Date: 2026-10-17

Supersedes the `let` spelling in §4.1, §4.3, §4.5, §4.7–§4.9, §4.11, §5.1–§5.3
and Appendix D of the *Atron Master Implementation Guide* v1.0. It does not
supersede any part of ADR-0001, ADR-0002 or ADR-0003.

## Context

The baseline guide gives Atron two declaration keywords: `let` for an immutable
binding and `var` for a mutable one. The maintainer reviewed that choice before
the lexer was designed and rejected it.

Two problems motivate the change.

The first is readability at the point of use. `let` is the immutable keyword in
Rust and Swift, but the mutable one in JavaScript, TypeScript and Kotlin. A
reader of Atron code therefore cannot infer mutability from `let` without
knowing which language the author came from, and the guide's own warnings lean
on that inference: W0002 is literally `'var x' is never reassigned; consider
'let'`. A keyword whose meaning depends on the reader's previous language is a
poor default for a language that advertises predictability.

The second is that `const` states the stronger guarantee directly. Atron has no
`const` today, and the keyword table count (22) is a published invariant, so the
change can be a one-for-one substitution rather than a keyword addition. `const`
is also the spelling used by Go, C, C++, C#, JavaScript and PHP for a binding
that cannot be reassigned, so the substitution reduces, rather than adds,
cross-language ambiguity.

The change is cheap now and expensive later. Before M1 there is no lexer, no
parser, no AST, no example program that anyone can run, and no release. After M1
the token kind is in golden files; after M2 it is in the AST; after M10 it is in
published code. The decision is therefore recorded at the last moment at which
it costs a documentation edit.

## Decision

`const` replaces `let` as the keyword for an immutable declaration.

```atron
const name = "Atron"    // immutable binding; initializer required
var count = 0           // mutable binding
var pending             // initialized to nil
```

The following rules are ratified with it:

1. `const name = expression` creates an immutable binding and **requires** an
   initializer. `const name` alone is a compile-time error.
2. `var name = expression` creates a mutable binding. `var name` without an
   initializer initializes to `nil`.
3. Immutability applies to the **binding, not recursively to the referenced
   object**. `const xs = []` still permits `xs.push(1)`; only reassignment of
   `xs` itself is rejected.
4. Assignment to a `const` binding remains diagnostic **E0303**
   (`cannot assign to immutable binding 'x'`). Re-declaration in the same scope
   remains **E0302**; use before declaration remains **E0301**.
5. The lexical scope, shadowing, capture and hoisting rules are **unchanged**
   from the guide §4.3. `const` and `var` are both **not hoisted**: a name becomes
   visible at its declaration, and use before it is E0301. The guide's hoisting
   rule names exactly two declaration kinds — `fn` and `struct` — and this ADR
   adds nothing to it, because renaming the immutable keyword does not change
   which declarations are hoisted.
6. `let` is **removed from the keyword table** and lexes as an ordinary
   identifier (`Ident`). A legacy `let name = 1` is therefore not a declaration
   at all: the parser reports **E0201** (`expected statement`) with a migration
   hint naming `const`.
7. There is **no** `Let` token kind, **no** `LetStmt` AST node and **no**
   compatibility alias at any layer.
8. The keyword count stays at **22**. `Const` occupies the slot `Let` held, so
   the 62-kind token inventory is unchanged in size.

The declaration AST node is a single struct for both spellings, so the parser
and the resolver cannot diverge on which keyword produced which node:

```go
type VarDecl struct {
    Name    *Ident
    Pub     bool
    Init    Expr
    Mutable bool // false => const, true => var
}
```

The parser enforces the `const` initializer requirement when it builds this
node; the resolver and later stages read only `Mutable`. Nothing downstream
needs to know which keyword produced the node, so the two spellings cannot
diverge in how they are scoped, captured or reported.

`impl` is deliberately **not** in the hoisting rule above, and this ADR does not
add it. The guide hoists `fn` and `struct` and says nothing about `impl`; whether
an `impl` block must be visible before its textual position, and how method
registration interacts with declaration order, is a **separate open design
question**. It is recorded as a future issue rather than settled here by
implication, because adding a third hoisted kind would be a language change in
its own right and does not follow from renaming `let` to `const`.

## Alternatives considered

**Keep `let` and add `const` alongside.** Rejected. Two spellings for one
concept, and the reader still has to guess which one a codebase uses. It also
makes 23 keywords, changing a number the guide publishes and the specification
now fixes.

**Keep `let` as a deprecated alias for `const`.** Rejected. An alias has to be
either a second token kind (contradicting "no `Let` token kind") or a
context-sensitive keyword (contradicting the fixed keyword table and making the
declaration site ambiguous). Since no Atron code exists in the wild, there is
nothing to keep compatible with; the alias would exist only to preserve a
removed keyword.

**Adopt `let mut` / `let` (Rust style).** Rejected. It keeps the keyword that
causes the cross-language confusion, adds a second token in the declaration
position, and departs further from the guide than a plain substitution does.

**Keep `let` and rename nothing.** This is the status quo and was the strongest
alternative: it needs no ADR and no specification churn. It was rejected because
the cost of the change is bounded entirely by this milestone — the guide is not
in version control and can be superseded by tracked documents, and no
implementation exists to migrate.

**Make `const` deep or transitively immutable.** Rejected for 0.1. It would
require either a deep copy on every `const` or a new permission model in the
resolver, and it contradicts the guide's stated rule and its own example
(`const xs = []` followed by `xs.push(1)`).

## Consequences

Positive:

- Mutability is readable at the declaration without knowing the author's
  background, and matches the majority spelling for "cannot be reassigned".
- The keyword table stays at 22 entries and the token inventory stays at 62
  kinds, so no published invariant changes.
- One `VarDecl` node with a `Mutable` flag removes a class of parser/resolver
  divergence and keeps M2 from needing two near-identical AST types.
- E0201 with a migration hint gives a legacy `let` user a precise, actionable
  error instead of a puzzling "undefined name `let`".

Negative and risks:

- **Every example in the guide that uses `let` becomes stale.** They are updated
  in the tracked specification, but the PDF is excluded from version control by
  `.gitignore` and is deliberately **not** modified. A reader holding the PDF
  will still see `let`; the specification is the normative document from M1
  onward and states that explicitly.
- Diagnostic help text in the guide (`declared with 'let'`) and the W0002
  suggestion (`consider 'let'`) must be reworded by the implementation that
  introduces them (M3 and M6), not by a mechanical search-and-replace.
- `const` is a Go keyword, so Go identifiers in the implementation need another
  name where `const` would have been natural, and the Atron keyword table cannot
  be derived mechanically from `go/token`.
- `let` being a legal identifier is two-edged: legacy-looking code lexes
  cleanly and parses further before failing, so the E0201 message and its hint
  carry the whole diagnostic burden.

Compatibility: this is a **pre-1.0 breaking change with no released surface**.
No published program, tag, module version or binary is affected, because none
exists. It is nonetheless recorded in `CHANGELOG.md` under `Unreleased →
Changed` rather than adopted silently.

Affected milestones: **M1** (token kind `Const`, keyword table, goldens),
**M2** (`VarDecl`, `const`-without-initializer parse error, `let` → E0201),
**M3** (symbol kind, E0303 and W0002 wording, hoisting tests), **M8** (the
formatter prints `const`), and any documentation that quotes a declaration.

## Follow-up

- Recorded as a language change in issue
  [#3](https://github.com/CosmicLang/atron-lang/issues/3), which also lists the
  guide conflicts and the rejected alternatives.
- Encoded in the normative specification: `docs/spec/lexical.md` (keyword
  table, identifier rule) and `docs/spec/grammar.md` (`const`/`var`
  declarations, entry point).
- **No implementation exists, and this record does not authorise M1.** Lexer and
  parser work starts only under a separate authorisation, and the M1 goldens
  must encode the rules above.

## Maintainer authority and chronology

This decision was taken by the maintainer of `CosmicLang/atron-lang` before M1
began. It is **not** the outcome of a community vote, discussion or consensus
process, and no such process is claimed retrospectively. The public record in
issue #3 was opened *after* the decision, so that the change is auditable and
can be objected to. A substantive objection is answered by a superseding ADR,
not by editing this one to look as though it was never made.
