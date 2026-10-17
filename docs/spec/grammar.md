# Atron grammar specification (M2 scope, ratified at M1)

Status: **ratified design contract; not yet implemented.** Normative for M2.
Lexical rules are in [`lexical.md`](lexical.md).

Baseline: *Atron Master Implementation Guide* v1.0 (17 October 2026), §4.3,
§4.11 and §5.2. Approved deviations:
[`ADR-0004`](../decisions/0004-const-declarations.md) (`const`) and
[`ADR-0005`](../decisions/0005-explicit-entrypoint.md) (explicit entry point,
single grammar, contextual validation).

No parser, AST or semantic analyzer exists. This document describes the syntax
that M2 must implement, and nothing more.

## 1. Declarations and the declaration node

`const` replaces `let` in the baseline grammar. The full set of declaration
alternatives is:

```ebnf
program        = { declaration } EOF ;
terminator     = NEWLINE | ";" | lookahead("}") | EOF ;
declaration    = import_decl | struct_decl | impl_decl | fn_decl | var_decl | statement ;

import_decl    = "import" ( IDENT | STRING ) [ "as" IDENT ] terminator ;
var_decl       = [ "pub" ] "const" IDENT "=" expression terminator
               | [ "pub" ] "var" IDENT [ "=" expression ] terminator ;
fn_decl        = [ "pub" ] "fn" IDENT "(" [ params ] ")" fn_body ;
fn_body        = block | "=>" expression terminator ;
params         = IDENT { "," IDENT } [ "," ] ;
struct_decl    = [ "pub" ] "struct" IDENT "{" [ field_list ] "}" ;
field_list     = IDENT { ( "," | NEWLINE ) IDENT } [ "," ] ;
impl_decl      = "impl" IDENT "{" { fn_decl } "}" ;
```

The two deviations from §4.11 are that `var_decl` uses `const` where the guide
uses `let`, and that nothing else in the production changes.

Rules:

- `const` **requires** `= expression`. `const name` with no initializer is a
  parse error at the declaration, reported as **E0201** with a hint that `const`
  needs an initializer. (`var` may omit it and defaults to `nil`.)
- The declaration builds one AST node, never two:

  ```go
  type VarDecl struct {
      Name    *Ident
      Pub     bool
      Init    Expr
      Mutable bool // false => const, true => var
  }
  ```

  There is no `LetStmt`, no `ConstStmt` and no `VarStmt`.
- `pub` is permitted on `const`, `var`, `fn` and `struct`, and is rejected
  elsewhere. A `pub` on a non-exportable declaration is a syntax error.
- `let` is **not** a keyword. `let name = 1` therefore parses as an expression
  statement beginning with an identifier, fails where `name` follows, and is
  reported as **E0201** (`unexpected token, expected <X>`) with a migration hint
  naming `const`. Do not add special-case recovery for `let`: the general
  statement error path plus the hint is the contract.
- Function declarations may be expression-bodied (`=>`) anywhere the grammar
  allows `fn_decl`, including when the function is named `main`. Deciding that
  an expression-bodied `main` is not an entry point belongs to the contextual
  layer (§4), **not** to the parser.

## 2. Statements, precedence and terminators

The statement, block and expression grammar is unchanged from the guide:

```ebnf
statement      = if_stmt | while_stmt | for_stmt | return_stmt | break_stmt
               | continue_stmt | block | assign_stmt | expr_stmt ;
block          = "{" { declaration } "}" ;
if_stmt        = "if" expression block [ "else" ( if_stmt | block ) ] ;
while_stmt     = "while" expression block ;
for_stmt       = "for" IDENT [ "," IDENT ] "in" expression block ;
```

Precedence and associativity are as ratified in the guide: `or`, then `and`,
then equality, then comparison, then `..`/`..=`, then `+ -`, then `* / %`, then
the right-associative `**`, then prefix `-`/`not`, then postfix `( [ . ?`.
There is no `!`, `&&` or `||`; the keywords are `not`, `and`, `or`.

Terminator handling, restated because two of the M1 corrections touched it:

- Inside `(...)`, `[...]`, and the braces of a map or struct literal, the parser
  **must ignore** `Newline` tokens.
- Inside a block's braces, `Newline` **is** a terminator and is not ignored.
- A `;` where a terminator is optional is not silently accepted as whitespace:
  `fn main(); {` is invalid, because only a single `Newline` may separate a
  header from its block (§3.1).
- `else` may follow a `}` on the same line or after a single `Newline`; the

## 3. Block layout, one grammar, and contextual validation

### 3.1 A block-opening `{` on its own line

Both of these are valid and equivalent:

```atron
fn main() {
    println("Atron")
}
```

```atron
fn main()
{
    println("Atron")
}
```

The lexer emits a `Newline` after `)` in both cases and never suppresses one
because `{` follows (`lexical.md` §7). The **parser** may consume exactly one
`Newline` immediately before a required block-opening `{`, and only in a
recognised header context. That list is **closed**:

| Header | May skip one `Newline` before `{` |
|---|---|
| `fn` declaration header | yes |
| `if` header | yes |
| `else` | yes |
| `while` header | yes |
| `for` header | yes |
| `impl` header | yes |
| `struct` header | **no** — a struct body is not a block |
| a map or struct **literal** | **no** — literal braces are expression context |
| any statement inside a block | **no** — `Newline` is a terminator there |

The skip is **exactly one** token, and only when the next token is `{`. It never
consumes `;`, never consumes two `Newline`s, and never consumes anything when a
statement follows the header instead of a body. A regression corpus must assert
that `fn main()` followed by `whatever()` followed by `{}` is **rejected**, so
the rule cannot degrade into "ignore newlines before any brace".

The formatter canonicalises to `{` on the header line. Accepting the alternative
input layout is a parser property, not a formatter property.

### 3.2 One grammar, no script mode

There is **one grammar** for every `.atrn` file. The parser builds a general
module AST using the same declaration and statement productions regardless of
how the file will be used, and `sema` applies ordinary scope, duplicate-name and
undefined-name rules without assuming the module is executable.

Executable-ness is decided **after** parsing, by a separate **contextual
program-validation layer**:

| Layer | Knows about `main`? |
|---|---|
| `internal/lexer` | no — `main` is an `Ident` |
| `internal/parser` | no — a parameterized or expression-bodied `fn main` is ordinary syntax |
| `internal/sema` | no — E0301/E0302/E0303 apply uniformly |
| contextual program validation | **yes** — the only layer that does |

That layer is called by the CLI driver (M4) and later by the module loader (M7).
It does **not** live in `internal/parser` or `internal/sema`, so the REPL, the
test runner, `Eval` and the inspection commands keep working without a `main` by
construction, rather than by flags threaded through the parser.


## 4. The executable entry point

An executable root must contain exactly one module-scope declaration

```atron
fn main() {
    // body
}
```

with all of the following properties:

| Property | Rule | Violation |
|---|---|---|
| Name | exactly `main` | E0310 |
| Parameters | none | E0310 |
| Body | a block, not `=> expression` | E0310 |
| Scope | module level | E0310 — a nested `main` does not count |
| Form | a named function declaration, not a binding | E0310 |
| Count | exactly one | E0302 if two are declared in one scope |
| Visibility | `pub` is irrelevant | permitted |

- `main` stays an **ordinary identifier**. A `fn main()` nested inside another
  function is legal and irrelevant to the entry point.
- `const main = fn() { ... }` is a valid declaration of a binding named `main`
  and is **not** an entry point.
- `fn main(a) { ... }` and `fn main() => 1` are **valid syntax** and must not
  produce a parse error. They are rejected only by executable-context
  validation, with **E0310**.
- **E0310** is reserved here and is new: the baseline catalogue ends at E0309.
  E0201 is explicitly **not** reused for these cases.

### 4.1 Executable roots versus libraries

| Aspect | Executable root | Imported library |
|---|---|---|
| `main` required | **yes** | no |
| Top-level statements | **rejected** | permitted (guide §4.9 semantics) |
| Module-scope `const`/`var` initializers | permitted | permitted |
| Module initialization | imported modules, then root initializers in order, then `main()` once | once per canonical path, via the import cache |
| A function named `main` | the entry point | never invoked automatically |

A further contextual rule: a **standalone top-level statement** at module scope
of an executable root — a bare call, loop or conditional — is a contextual
error, because it would execute outside any function. It is **not** an error in
a library. The asymmetry is deliberate; see ADR-0005, "Consequences".

Contexts that require **no** `main`: the REPL, test discovery and execution,
embedded `Eval`, and the inspection commands (`fmt`, `tokens`, AST dump,
disassembly). `atron run` requires it. `atron check` performs generic validation
and adds entry-point validation only in an executable-target context; **no CLI
flag for that is defined here**, because inventing one before the CLI phase
would document a command that does not exist.


## 5. Conformance examples (§4.12 of the guide, corrected)

The guide requires its §4.12 programs in `examples/` and `testdata/run/` with
`.golden` files, and they are the acceptance criterion for M4–M7. They appear
here corrected: `const` replaces `let`, and every **executable** root gains a
`fn main()` wrapper. **Every established output is preserved exactly**,
including the accented text in `hello.atrn`. `config.atrn` stays a library, so
it gains no `main` and keeps the guide's module semantics — including
permission for top-level statements, which this module does not happen to use.

### 5.1 `hello.atrn` → `Olá, Atron!`

```atron
fn main() {
    println("Olá, Atron!")
}
```

### 5.2 `fib.atrn` → `0 1 1 2 3 5 8 13 21 34`

```atron
fn fib(n) {
    if n < 2 { return n }
    return fib(n - 1) + fib(n - 2)
}

fn main() {
    var parts = []
    for i in 0..10 { parts.push(str(fib(i))) }
    println(parts.join(" "))
}
```

### 5.3 `wordfreq.atrn` → `a: 3`, `b: 2`, `c: 1`, one per line

```atron
fn main() {
    const text = "a b a c b a"
    var freq = {}
    for w in text.split(" ") {
        freq[w] = freq.get(w, 0) + 1
    }
    for k, v in freq { println("{k}: {v}") }
}
```

Insertion order is preserved, not sorted, and the golden file encodes it.

### 5.4 `config.atrn` — an imported library, no `main`

```atron
import fs
import json

pub fn load(path) {
    const text = fs.read_file(path)?
    const data = json.parse(text)?
    return data
}
```

This module is imported and never executed as a root. It needs no `main` and is
exempt from the executable-root restriction on top-level statements.

### 5.5 `main.atrn` → `loaded N keys`

```atron
import "./config" as config

fn main() {
    const cfg = config.load("app.json")
    if is_error(cfg) {
        println("could not load config: {cfg.message}")
        exit(1)
    }
    println("loaded {len(cfg)} keys")
}
```

The declaration initializer and the `if` move inside the body because top-level
statements are rejected in an executable root. The observable behaviour,
including the `exit(1)` path, is unchanged.

### 5.6 `shapes.atrn` → `[3.14159, 6]`

```atron
struct Circle { r }
struct Rect { w, h }

impl Circle { fn area(self) => 3.14159 * self.r ** 2 }
impl Rect   { fn area(self) => self.w * self.h }

fn main() {
    const shapes = [Circle{ r: 1.0 }, Rect{ w: 2, h: 3 }]
    const areas = shapes.map(fn(s) => s.area())
    println(areas)
}
```

### 5.7 Forms that are **not** entry points

```atron
// Rejected by executable-context validation with E0310 — NOT by the parser:
fn main(args) { println(args) }      // parameterized
fn main() => println("x")            // expression-bodied

// Rejected by the parser with E0201 plus a migration hint:
const ok = 1
let legacy = 2                       // `let` is an identifier, not a keyword
```


## 6. Diagnostics referenced by this document

| Code | Meaning | Layer |
|---|---|---|
| `E0101`–`E0105` | lexical errors | lexer — see `lexical.md` |
| `E0201` | unexpected token / expected `<X>`; also a legacy `let` and `const` without an initializer | parser |
| `E0301` | undefined name | sema |
| `E0302` | duplicate declaration in one scope | sema |
| `E0303` | assignment to an immutable (`const`) binding | sema |
| **`E0310`** | **missing or invalid executable entry point** | **contextual validation** |

`E0310` is the only code this specification introduces. Adding it to the tracked
error catalogue requires explicit maintainer approval, and it must not be
allocated to any other feature.

## 7. Conformance checklist for M2

| # | Requirement |
|---|---|
| 1 | `const x = 1` produces one `VarDecl` with `Mutable == false`; `var x = 1` produces `Mutable == true`. |
| 2 | `const x` with no initializer is E0201 with a hint; it never silently becomes a `nil` initializer. |
| 3 | There is no `LetStmt` type; `let x = 1` is E0201 with a `const` hint and does not parse as a declaration. |
| 4 | `fn main()` followed by `{` on the next line parses identically to the same-line form. |
| 5 | A header followed by an ordinary statement and then `{}` is **rejected**; the skip covers one `Newline` before `{`, nothing else. |
| 6 | `fn main(); {` is **rejected**; `;` is never consumed as optional whitespace. |
| 7 | The parser accepts `fn main(a) {}` and `fn main() => 1` without a diagnostic. |
| 8 | `Newline` is still a terminator inside a block and still ignored inside parentheses, brackets and literal braces. |
| 9 | `ast.Print` produces S-expressions suitable for golden tests, one golden per §5 example. |
| 10 | The parser terminates on arbitrary input and never panics (`FuzzParser`). |

## 8. What this document does not claim

- It does not claim the parser, AST or contextual validation layer exists. **M2
  and M4 have not started.**
- It does not define type checking, evaluation order or a message catalogue;
  those are M3, M4 and M6.
- It does not define an `atron check` flag for executable-target validation.
- It does not amend the original guide PDF, which is excluded from version
  control; it records the approved deviations in tracked, reviewable form.

  parser consumes that `Newline` only when the next token is `else`.
