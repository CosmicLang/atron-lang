# ADR-0005: Executable roots declare an explicit `fn main()`

Status: accepted — Date: 2026-10-17

Introduces a requirement the *Atron Master Implementation Guide* v1.0 does not
contain. It does not supersede any existing ADR.

## Context

The baseline guide treats every `.atrn` file as a script. §4.9 says each file is
a module, evaluated once and cached by canonical path, and its `main.atrn`
example is a bare `println(...)` at top level. Under that model "what happens
when I run this file?" has the same answer as "what happens when anything
imports it?": the module body runs.

That model is workable for small examples and becomes ambiguous as programs
grow, for three reasons.

1. **The entry point is not locatable.** With top-level statements, control
   enters at the first declaration and runs downward, interleaved with the
   initialization of the declarations themselves. A reader cannot point at the
   line where the program starts, and a tool cannot report "the entry point is
   here".
2. **Execution and import are indistinguishable.** A module that prints or
   mutates globals at top level behaves differently depending on *who imported
   it and how many times*. The import cache makes the behaviour deterministic
   but not obvious, and library code acquires an implicit side effect the author
   may not have intended.
3. **Static validation has nothing to validate.** "Does this program have a
   valid entry point?" is only a decidable question if a declaration is
   designated as one. Without that, a typo in a function name is a runtime
   "undefined name" instead of a compile-time error.

The guide itself hints at the problem: its E0304 (`return` outside a function)
and E0306 (`?` outside a function) exist precisely because top-level code is not
inside a function. Adding an explicit entry point gives those rules a natural
home rather than a special case.

## Decision

An executable root must contain exactly one top-level named function declaration
`fn main() { ... }`.

```atron
// hello.atrn — an executable root
fn main() {
    println("Hello, world!")
}
```

The entry point contract:

- **Name** `main`. `main` stays an **ordinary identifier**; it is not a keyword
  and not reserved.
- **Zero parameters.** `fn main(args) { ... }` is not an entry point.
- **Block body.** `fn main() => expression` is valid *syntax* but not a valid
  entry point.
- **Module scope.** Exactly one such declaration at the top level of the root
  module.
- **Invoked once** per executable run, with no arguments.
- **`pub fn main()` is permitted.** Visibility is orthogonal to entry-point
  validity.
- A **nested** function named `main` is not an entry point.
- A **binding** such as `const main = fn() { ... }` does not satisfy the
  contract; the required form is a named function declaration.

Validation is **contextual**, not syntactic:

| Context | `main` required |
|---|---|
| `atron run` (executable root) | **Yes** |
| `atron check` in an executable-target context | **Yes** |
| Imported library module | No |
| REPL | No |
| Test discovery and execution | No |
| Embedded `Eval` (`pkg/atron`) | No |
| `atron fmt`, `atron tokens`, AST dump, disassembly | No — inspection must not require an entry point |

Diagnostics:

- **E0310** (new, reserved here) reports a missing or invalid executable entry
  point: absent, parameterized, expression-bodied, or not a module-scope named
  function.
- **E0302** remains the duplicate-declaration diagnostic and therefore already
  covers two `fn main` declarations in one scope.
- **E0201** is **not** reused. `fn main(a) { ... }` and `fn main() => 1` are
  syntactically valid; reporting a parse error for them would make a parser
  diagnostic depend on how the file is later invoked, and would mislead a reader
  of a library that legitimately contains a function named `main`.

One grammar serves every file. There is no second parser and no
"script mode": the parser builds the same module AST regardless of invocation
context, the semantic analyzer applies ordinary scope and duplicate rules
without assuming executability, and a **separate contextual program-validation
layer** decides whether the module is being used as an executable root. That
layer is called by the CLI driver at M4 and reused by the module loader at M7;

### Execution and initialization contract

- **Executable root:** module-scope declarations, imports and declaration
  initializers are permitted. **Standalone top-level statements are rejected** —
  a bare call, loop or conditional at module scope of an executable root is a
  contextual error, because it would execute outside any function.
- **Initialization order:** resolve and initialize required imported modules
  first (following the module dependency and import-cycle rules), then evaluate
  the root's module-level declaration initializers in source order, then invoke
  `main()` exactly once.
- **Imported libraries** are exempt from `main` and keep the guide's existing
  semantics, **including permitted top-level statements**, executed once per
  canonical path under the import cache. A library function named `main` is
  never invoked automatically, and executable-context restrictions are never
  imposed on a library.
- **Module initialization stays once-per-module** per the guide's caching
  contract; this decision adds an entry point, it does not add a second
  initialization path.

### Return value and exit status

The runtime **discards** the value returned by `main`. A `main` that completes
normally exits with code `0`. There is **no implicit conversion** of a returned
integer or error into a process exit status; `exit(n)` remains the only
mechanism for setting it.

The established CLI and runtime classification is unchanged:

| Code | Meaning |
|---|---|
| `0` | successful completion |
| `1` | unrecoverable runtime error or panic |
| `2` | compilation or static diagnostic failure |
| `64` | invalid CLI usage |

A missing or invalid entry point is a **static** failure, so `atron run` on such
a file exits `2`, not `1`.

### Scope of this decision by milestone

M4 implements only the **single-file executable model**: parse one file, validate
the entry point, run it. M4 does **not** gain an import resolver or a module
cache, and must not acquire an M7 dependency. Full import resolution, caching and
library-module semantics land in M7. Until then, no claim is made that imports
behave as described above.

it lives outside `internal/parser` and `internal/sema`.

## Alternatives considered

**Script-style top level with no entry point (the baseline).** Rejected for the
reasons in *Context*: the start of execution is not locatable, and validating
"is this program runnable?" is impossible without a designated declaration.

**Script style plus an *optional* `main`.** Rejected as the worst of both: the
reader still cannot tell from the module whether the entry point is `main` or
the top of the file, and every tool needs both code paths. Ambiguity that is
only resolved at runtime is exactly what this language is trying to avoid.

**`main` as a reserved keyword.** Rejected. It would break programs that legally
use `main` as a variable or field name, and it buys nothing: selecting the entry
point is a module-scope name lookup, not a lexical property. It would also
change the guide's fixed 22-keyword table, which this ADR does not touch.

**`const main = fn() { ... }` as the entry point.** Rejected. It makes the entry
point an initializer expression, so "when does `main` run?" depends on
declaration-initializer ordering, and a helper named `main` in a library becomes
ambiguous.

**Reuse E0201 for a wrong-shaped `main`.** Rejected, as described under
*Decision*: the code would report a parse error for syntactically valid syntax
and would make a parser diagnostic depend on invocation context.

**Separate executable and library parsers (`main`-mode vs library-mode).**
Rejected. Duplicated grammar, duplicated recovery tests and inevitable drift;
the difference is contextual, so it belongs in a contextual validation layer.

**Return an `int` from `main` and use it as the exit code (C/Go style).**
Rejected. It couples `main`'s type to the OS, conflicts with `nil` returns from a
dynamically typed language, silently truncates out-of-range values, and
duplicates `exit(n)`, which already exists and is explicit.

**Wrap the program in an implicit generated `main`.** Rejected. It hides the
invocation site from the reader and the debugger, and it makes the "duplicate
`main`" rule depend on whether the user wrote one.

## Consequences

Positive:

- Execution has one statically locatable start, which the formatter, the test
  runner, the REPL filter and the documentation can all point at.
- "This program has no entry point" becomes a compile-time diagnostic (E0310)
  with an exit code of `2`, instead of a runtime "undefined name" or a silent
  no-op.
- Library modules stay free of an entry-point obligation, so the guide's module
  semantics and import cache survive unchanged.
- Keeping validation out of `parser` and `sema` means the REPL, tests, `Eval` and
  the inspection commands work by construction rather than by special-casing.

Negative and risks:

- **Every executable example in the guide becomes wrong** and must gain a
  `fn main()` wrapper; `docs/spec/grammar.md` shows the corrected forms and
  preserves the expected program output.
- E0310 is **net-new** and must be added to the tracked error catalogue; the
  guide's Appendix D stops at E0309. It is reserved here so no other feature can
  claim the number first.
- "Top-level statements are rejected in an executable root but permitted in a
  library" is a context-dependent rule, which is the kind of asymmetry that
  surprises users. It is deliberate and documented in one place, but it requires
  precise tests in M4 so the error never appears for a library.
- A `main` that returns a value now has that value silently discarded, which is a
  new way to write a program that computes something nobody reads. Unused-value
  warnings do not exist yet; that is a known M6 gap.

Compatibility: pre-1.0, no released surface, no migration burden. Recorded in
`CHANGELOG.md` under `Unreleased → Changed`.

Affected milestones: **M2** (module AST accepts any `fn main`; no entry-point
rejection in the parser), **M3** (E0302 duplicate rule; no `main` assumption),
**M4** (contextual validation, single-file execution, exit code `2`, E0310
emitted by the driver), **M7** (the loader reuses the same layer for roots only),
**M8** (test runner and REPL do **not** require `main`), **M6** (E0310 catalogue
entry and rendering).

## Follow-up

- Recorded as a language change in issue
  [#3](https://github.com/CosmicLang/atron-lang/issues/3).
- The contextual validation layer has no home yet: it is a **planned** package
  boundary owned by the M4 driver, not an implemented one. Naming it before M4
  is design, not code.
- **No implementation exists, and this record does not authorise M1 or M4.**
- Whether `atron check` should default to executable-target validation, and how
  that context is selected without inventing an undocumented flag, is deferred
  to the CLI documentation phase. This ADR deliberately defines no flag.

## Maintainer authority and chronology

This decision was taken by the maintainer of `CosmicLang/atron-lang` before M1
began, as a deliberate **addition** to the baseline guide rather than a
correction of it. It is **not** the outcome of a community vote, discussion or
consensus process, and none is claimed retrospectively. Issue #3 was opened
after the decision to make it auditable; objections are answered by a
superseding ADR rather than by rewriting this record.

