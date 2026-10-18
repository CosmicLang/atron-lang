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

## 1. The complete grammar

`const` replaces `let` in the baseline grammar, and one grammar serves every
file (ADR-0005). What follows is the **complete** M2 grammar, transcribed from
§4.11 of the guide with exactly those two changes and no other. No production is
abridged, and no production is invented.

Declarations:

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

Statements and blocks:

```ebnf
statement      = if_stmt | while_stmt | for_stmt | return_stmt | break_stmt
               | continue_stmt | block | assign_stmt | expr_stmt ;
block          = "{" { declaration } "}" ;
if_stmt        = "if" expression block [ "else" ( if_stmt | block ) ] ;
while_stmt     = "while" expression block ;
for_stmt       = "for" IDENT [ "," IDENT ] "in" expression block ;
return_stmt    = "return" [ expression ] terminator ;
break_stmt     = "break" terminator ;
continue_stmt  = "continue" terminator ;
assign_stmt    = target assign_op expression terminator ;
assign_op      = "=" | "+=" | "-=" | "*=" | "/=" | "%=" ;
target         = IDENT | postfix "." IDENT | postfix "[" expression "]" ;
expr_stmt      = expression terminator ;
```

Expressions. The production chain **is** the precedence hierarchy of §4.4: each
nonterminal calls the next-higher level, so the eleven levels of §2 are encoded
here rather than restated (`not` sits below `equality`, `-` sits above `*` and
`/` and below `**`):

```ebnf
expression     = or_expr ;
or_expr        = and_expr { "or" and_expr } ;
and_expr       = not_expr { "and" not_expr } ;
not_expr       = "not" not_expr | equality ;
equality       = comparison [ ( "==" | "!=" ) comparison ] ;
comparison     = range [ ( "<" | "<=" | ">" | ">=" ) range ] ;
range          = additive [ ( ".." | "..=" ) additive ] ;
additive       = multiplicative { ( "+" | "-" ) multiplicative } ;
multiplicative = unary { ( "*" | "/" | "%" ) unary } ;
unary          = "-" unary | power ;
power          = postfix [ "**" unary ] ;
postfix        = primary { "(" [ args ] ")" | "[" expression "]" | "." IDENT | "?" } ;
args           = expression { "," expression } [ "," ] ;
```

Primary expressions, collection literals, struct literals and function literals:

```ebnf
primary        = INT | FLOAT | STRING | RAWSTRING | "true" | "false" | "nil" | IDENT
               | "(" expression ")" | list_lit | map_lit | struct_lit | fn_lit ;
list_lit       = "[" [ args ] "]" ;
map_lit        = "{" [ entry { "," entry } [ "," ] ] "}" ;
entry          = ( IDENT | STRING | "[" expression "]" ) ":" expression ;
struct_lit     = IDENT "{" [ init { "," init } [ "," ] ] "}" ;   (* not in headers *)
init           = IDENT ":" expression ;
fn_lit         = "fn" "(" [ params ] ")" ( block | "=>" expression ) ;
```

Three notes that the guide attaches to this grammar, preserved here:

- Inside a `STRING`, `{ expression }` is analysed recursively and `\{` produces
  a literal `{`. The lexer returns the raw text and the parser splits it out
  (`lexical.md` §8.3).
- A `{` at the start of a statement opens a block; in expression position it
  opens a map literal.
- A `Newline` immediately before the `}` of a literal — or before `)` or `]` —
  is ignored, per the terminator rules in §2.

**Every production above is parser-level and context-free.** Nothing in it
mentions `main`, executable roots or libraries; those are contextual validation
rules and appear only in §3.2 and §4, so that the grammar describes one language
and the context rules describe one use of it.

Rules for `const` and `var`:

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

## 2. Precedence, associativity and terminators

The statement and block productions are in §1. Precedence and associativity are
**exactly** those of guide §4.4, from lowest binding level (1) to highest (11):

| Level | Operators | Associativity | Note |
|---|---|---|---|
| 1 | `or` | left | lowest binding |
| 2 | `and` | left | |
| 3 | prefix `not` | right | `not a == b` ≡ `not (a == b)` |
| 4 | `==`, `!=` | **non-associative** | `a == b == c` is a syntax error (E0201) |
| 5 | `<`, `<=`, `>`, `>=` | **non-associative** | |
| 6 | `..`, `..=` | **non-associative** | `0..n+1` ≡ `0..(n+1)` |
| 7 | `+`, binary `-` | left | |
| 8 | `*`, `/`, `%` | left | |
| 9 | prefix `-` | right | `-2 ** 2` ≡ `-(2 ** 2)` = `-4` |
| 10 | `**` | right | `2 ** 3 ** 2` ≡ `2 ** (3 ** 2)`; `2 ** -1` is valid |
| 11 | call `f()`, index `a[i]`, field `a.b`, postfix `?` | left (postfix) | highest binding |

Two consequences are called out because they are the cases an implementation is
most likely to get wrong, and both are asserted by the conformance corpus:

- **`not` binds looser than every comparison.** Level 3 sits *below* levels 4–10,
  so `not a == b` parses as `not (a == b)`, never as `(not a) == b`. `not` is a
  prefix operator of its own level, not part of a shared "unary" level.
- **Prefix `-` binds tighter than `*`, `/`, `%` but looser than `**`.** Level 9
  sits above level 8 and below level 10, so `-2 ** 2` parses as `-(2 ** 2)` and
  evaluates to `-4`, while `-a * b` parses as `(-a) * b`. The two prefix
  operators are therefore at **different** binding powers and must not be
  collapsed into one `unary` level in the prose or in a precedence-climbing
  implementation. The EBNF in §1 encodes this correctly: `unary = "-" unary |
  power` sits between `multiplicative` and `power`, while `not_expr` sits between
  `and_expr` and `equality`.

Non-associativity is a **parse error**, not a silent left-fold: `a == b == c`,
`a < b < c` and `a..b..c` are each E0201, because the level 4, 5 and 6
productions admit at most one operator and do not recurse.

Assignment is a statement, never an expression: `a = b = 1` and `if (x = 1)` are
syntax errors. Valid targets are an identifier, `obj.field` and `obj[index]`.

There is no `!`, `&&` or `||`; the operators are the keywords `not`, `and`, `or`.

Terminator handling, restated because two of the M1 corrections touched it:

- Inside `(...)`, `[...]`, and the braces of a map or struct literal, the parser
  **must ignore** `Newline` tokens.
- Inside a block's braces, `Newline` **is** a terminator and is not ignored.
- A `;` where a terminator is optional is not silently accepted as whitespace:
  `fn main(); {` is invalid, because only a single `Newline` may separate a
  header from its block (§3.1).
- `else` may follow a `}` on the same line or after a single `Newline`; the
  parser consumes that `Newline` only when the next token is `else`. This is the
  one other closed point at which a `Newline` is skippable, and it is listed in
  the §3.1 table so that both skips live in one place.
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

The lexer's rule is purely local: a `Newline` is emitted when an actual
qualifying logical line break follows a token that can end a statement
(`lexical.md` §7). In the first form there is **no** line break after `)` — the
`{` follows on the same line — so **no** `Newline` is emitted there. In the
second form the line break does follow `)`, so exactly one `Newline` is
emitted. The lexer never suppresses a `Newline` because `{` happens to follow,
and it never invents one that the source does not contain.

The **parser** may consume exactly one `Newline` immediately before a required
block-opening `{`, and only in a recognised header context. That list is
**closed**:

| Header | May skip one `Newline` before `{` | May skip one `Newline` before `else` |
|---|---|---|
| `fn` declaration header | yes | — |
| `if` header | yes | — |
| `else` | yes | — |
| `while` header | yes | — |
| `for` header | yes | — |
| `impl` header | yes | — |
| `struct` header | **no** — a struct body is not a block | — |
| a map or struct **literal** | **no** — literal braces are expression context | — |
| any statement inside a block | **no** — `Newline` is a terminator there | — |
| `else` after a completed `if` block | — | yes — the second permitted skip |

The two skips are the **complete** set. The `else` skip is listed here rather
than in a separate section so that no third `Newline` exemption can be added
elsewhere without contradicting this table; it is also the only exemption that
follows a statement or block instead of preceding a required `{`.

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

Each of the following is a **complete, independent** source file. They are
separate on purpose: putting two `fn main` declarations in one file would
independently trigger **E0302** (duplicate declaration in one scope), which would
contaminate an E0310 example and leave the reader unable to tell which rule
fired.

`e0310_parameterized.atrn` — valid syntax, invalid entry point:

```atron
fn main(args) {
    println(args)
}
```

`e0310_expression_body.atrn` — valid syntax, invalid entry point:

```atron
fn main() => println("x")
```

`e0310_binding.atrn` — a binding named `main` is not an entry point:

```atron
const main = fn() {
    println("x")
}
```

`e0310_missing.atrn` — no entry point at all:

```atron
fn helper() {
    return 1
}
```

All four are **rejected by executable-context validation with E0310**, and only
there. None of them is a parser error: the parser accepts every one of these
files without a diagnostic.

For contrast, a genuine **syntax** error — reported by the parser, not by
entry-point validation — is E0201:

```atron
const ok = 1
const broken =                        // E0201: expected expression after '='
```

Legacy `let` is also an E0201 syntax error, because `let` is an ordinary
identifier rather than a keyword (ADR-0004):

```atron
let legacy = 2                        // E0201 with a `const` migration hint
```

And a **duplicate** declaration is E0302, a semantic error — deliberately shown
alone here, and with no second `main` in any of the E0310 examples above:

```atron
fn main() {
    println("first")
}

fn main() {
    println("second")
}
```

That last file is a duplicate-declaration error (**E0302**). It is *also* not a
valid executable root, but E0302 is the correct primary diagnostic: the
duplicate exists whether or not the file is executable, so it is reported by
`sema`, before any entry-point validation runs.


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

**Lexical versus parser responsibility for punctuation.** The lexer classifies
punctuation and reports nothing about its position; the parser decides whether
that punctuation is legal where it appears. Every `.` that does not continue a
numeric literal is a `Dot` token produced without a diagnostic — including a
bare `.`, `5.` and `.5` — and any error in those forms is **E0201** from this
layer, never `E0103`. `E0103` is reserved for malformed number literals and
integer overflow, and a `.` cannot cause it. Conversely, `a.b` is valid member
access at level 11 of §2 and must not be rejected. See `lexical.md` §4 and §8.2
for the token-level statement of the same rule.

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
| 11 | `-2 ** 2` parses as `-(2 ** 2)` and evaluates to `-4`; `-a * b` parses as `(-a) * b`; `2 ** 3 ** 2` parses as `2 ** (3 ** 2)`; `2 ** -1` is accepted. |
| 12 | `not a == b` parses as `not (a == b)`, and `not` never binds tighter than `==`. |
| 13 | `a == b == c`, `a < b < c` and `a..b..c` are each E0201; non-associativity is an error, not a left-fold. |
| 14 | `a = b = 1` and `if (x = 1)` are E0201; assignment never yields a value. |
| 15 | Every E0310 example in §5.7 is accepted by the parser and rejected only by executable-context validation; the deliberate two-`main` file is E0302, not E0310. |
| 16 | `.5` is rejected as **E0201** at the parser: the lexer produces `Dot`, `Int`(`5`) with no lexical diagnostic, so the error is a missing expression on either side of `Dot`, never E0103. `5.` and `a.` are likewise E0201, while `a.b` parses as member access. |

## 8. What this document does not claim

- It does not claim the parser, AST or contextual validation layer exists. **M2
  and M4 have not started.**
- It does not define type checking, evaluation order or a message catalogue;
  those are M3, M4 and M6.
- It does not define an `atron check` flag for executable-target validation.
- It does not amend the original guide PDF, which is excluded from version
  control; it records the approved deviations in tracked, reviewable form.
