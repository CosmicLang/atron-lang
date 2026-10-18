# Atron lexical specification (M1)

Status: **ratified design contract; not yet implemented.** Normative for M1.

Baseline: *Atron Master Implementation Guide* v1.0 (17 October 2026), §4.1 and
§5.1. This document is the normative statement of the lexer's observable
behaviour. Where it differs from the guide, it wins, and the difference is
recorded as an ADR: [`0004`](../decisions/0004-const-declarations.md) for the
`const` keyword and [`0005`](../decisions/0005-explicit-entrypoint.md) for the
entry point. The guide PDF is excluded from version control by `.gitignore` and
is not modified.

Grammar and semantics are in [`grammar.md`](grammar.md). Nothing in this
document describes behaviour that exists: no lexer code has been written, and
M1 is not authorised by this specification.

## 1. Source encoding

- Source files are **UTF-8**. A malformed byte sequence is invalid input; no
  transcoding and no lenient replacement is performed.
- `.atrn` is the source extension (ADR-0003).
- A **UTF-8 BOM** (`EF BB BF`) at the very start of the file is skipped and is
  **not** part of any token. A BOM anywhere else is an invalid character
  (E0101).
- A **shebang** line — the file starting with `#!` — is ignored up to and
  including its line break. It is not a comment token and is not preserved
  as lexical metadata. `#!` is not otherwise special.
- The lexer never modifies the source. All byte offsets refer to the input
  exactly as supplied.

## 2. Line endings and columns

Three line-ending representations are normalised **logically**, without altering
the source bytes:

| Sequence | Logical line breaks | Byte offset advance |
|---|---|---|
| LF (`\n`) | 1 | 1 |
| CRLF (`\r\n`) | 1 | 2 |
| CR (`\r`) alone | 1 | 1 |

Rules:

- Each of the three counts as **exactly one** logical line break. A standalone
  CR is **not** an illegal character merely because it is unpaired.
- Normalisation is **logical only**. It changes `Line`/`Col` bookkeeping; it
  never rewrites, deletes or substitutes a byte of the source. A CR byte stays
  in the source buffer and a CRLF stays two bytes, so every offset remains a
  true index into `src`.
- Because offsets stay byte-faithful, the lexeme invariant of §9 holds without
  exception for every **real, source-backed** token:

  ```
  Token.Lexeme == string(src[Pos.Offset:End.Offset])
  ```

  The only token this does not apply to is the synthetic `Newline` inserted by
  ASI (§7), which has empty `Lexeme` and no source span of its own.
- CR bytes therefore **can** appear inside a token's lexeme, and must. A
  multiline raw string (§8.4) that contains CRLF or standalone CR returns those
  actual bytes in its `Lexeme` and spans them with `Pos`/`End`, because the
  invariant above is checked by golden tests and by `src[Pos.Offset:End.Offset]`
  equality. Nothing in the lexer strips or rewrites them.
- Collected comments (§6) likewise preserve their **original source text**
  byte-for-byte, including any CR or CRLF inside a block comment, so that the
  formatter can reprint them exactly.
- A CR in a token lexeme is never rendered as a literal carriage return in a
  diagnostic message; escaping is a rendering concern (M6), never a change to
  the stored lexeme.
- **Columns count Unicode code points (runes), not bytes**, and start at 1.
  Lines start at 1. A tab advances the column by one, like any other code point;
  tab stops are a display concern, not a lexical one.
- `Pos` carries the **byte offset** of the first byte of the token so the source
  can be sliced; `Line` and `Col` are for humans and diagnostics.
- A `\r` inside a string literal or comment follows the context rules of
  §5 and §6; those rules decide whether it terminates the token.

## 3. Whitespace

Spaces, tabs, form feeds and vertical tabs separate tokens and produce nothing.
Outside literals and comments, whitespace carries no meaning beyond separating
tokens; indentation is never significant.

## 4. Tokens and the complete kind inventory

`Kind` is a `uint8`-backed enumeration whose zero value is `Illegal`. There are
exactly **62 kinds**:

| Group | Count | Members |
|---|---|---|
| Special / identifier | 4 | `Illegal`, `EOF`, `Newline`, `Ident` |
| Literal | 4 | `Int`, `Float`, `String`, `RawString` |
| Operator | 20 | `Plus`, `Minus`, `Star`, `Slash`, `Percent`, `StarStar`, `Eq`, `NotEq`, `Lt`, `LtEq`, `Gt`, `GtEq`, `Assign`, `PlusAssign`, `MinusAssign`, `StarAssign`, `SlashAssign`, `PercentAssign`, `DotDot`, `DotDotEq` |
| Delimiter / punctuation | 12 | `LParen`, `RParen`, `LBracket`, `RBracket`, `LBrace`, `RBrace`, `Comma`, `Colon`, `Semicolon`, `Dot`, `Question`, `FatArrow` |
| Keyword | 22 | `Const`, `Var`, `Fn`, `Return`, `If`, `Else`, `While`, `For`, `In`, `Break`, `Continue`, `Struct`, `Impl`, `Import`, `Pub`, `As`, `And`, `Or`, `Not`, `True`, `False`, `Nil` |

4 + 4 + 20 + 12 + 22 = **62**.

Clarifications that are normative:

- **`Comment` is not a kind.** Comments are preserved separately as lexical
  metadata (§6) so the formatter (M8) can reprint them; they never enter the
  parser's token stream.
- **`Newline` is produced by the lexer** (ASI, §7) and is a real kind, not a
  parser convention.
- **`Let` does not exist.** `const` replaced it; the count stayed 22
  (ADR-0004). There is no alias.
- **`main` is an ordinary identifier**, not a keyword (ADR-0005).
- `true`, `false` and `nil` are **keywords**, not literals: there is no
  distinct boolean or nil literal kind.
- The two-character operators are recognised greedily and longest-first:
  `**` before `*`, `..=` before `..`, `=>` before `=`, `==` before `=`, `!=`,
  `<=`, `>=` before their single-character prefixes, and the compound
  assignments `+= -= *= /= %=` before `+ - * / %`. There is no `===`, `!==`,
  `&&`, `||`, `!` or `++`/`--` in the language.
- **`.` is `Dot`, and `Dot` is never an error.** A `.` that does not continue a
  numeric literal lexes as exactly one `Dot` token, with no diagnostic of any
  kind. By the longest-match rule above, `..` is `DotDot` and `..=` is
  `DotDotEq`, so `...` is `DotDot` followed by `Dot`. `Dot` is the token that
  member access (`a.b`) and the §7 continuation rule refer to. Whether a `.`
  occurs in a syntactically legal position is a **parser** question answered by
  E0201; it is never a lexical one.

## 5. Identifiers

The identifier grammar is:

```
[\p{L}_][\p{L}\p{N}_]*
```

That is: a first character that is a Unicode letter or an underscore, followed
by any number of Unicode letters, Unicode digits or underscores.

Consequences, all normative:

- `_`, `_1`, `_name`, `__`, `x1`, `naïve`, `Ω`, `東京` are all **valid
  identifiers**. `_1` must **never** produce E0103 (or any other diagnostic): it
  is an identifier, not a number, because its first character is `_`.
- Identifiers are **case-sensitive**: `name`, `Name` and `NAME` are three
  distinct identifiers.
- An identifier **never** starts with a digit; a digit-initial sequence is a
  number literal (§8).
- Only `\p{L}` (letters) and `\p{N}` (numbers) are accepted. Other Unicode
  categories — marks alone, punctuation, symbols, emoji — are not identifier
  characters, so a combining mark that follows a base letter ends the token and
  the mark itself is an invalid character (E0101).
- When an identifier's text equals a keyword (§4), the keyword kind is emitted.
  The comparison is exact and ASCII-only, because every keyword is ASCII.

Convention, **not** enforced by the lexer: `snake_case` for variables and
functions, `PascalCase` for structs.

## 6. Comments

Two forms, both preserved as lexical metadata and neither producing a token:

- **Line comment:** `//` to the end of the logical line. It does not consume the
  line break, so that break still participates in ASI.
- **Block comment:** `/* … */`, **not nestable**: the first `*/` closes it, so
  `/* a /* b */` ends at the first `*/`. An unterminated block comment is
  **E0104**, reported at the opening `/*`, and the comment runs to end of file.
- A block comment may span logical line breaks and counts as whitespace for ASI
  (§7): a comment after a statement does not prevent the `Newline` from being
  inserted, and a comment that contains a line break does not itself force one.
- Comments are recorded with their positions and raw text so `atron fmt` can
  reprint them. The metadata type is an M1 implementation detail; the
  requirement that **no comment is discarded** is normative.

## 7. Automatic semicolon insertion (ASI)

The lexer inserts a `Newline` token at the end of a line **only if** the last
token emitted on that line can end a statement. That set is exactly:

```
Ident, Int, Float, String, RawString,
True, False, Nil,
Return, Break, Continue,
RParen, RBracket, RBrace, Question
```

Equivalently, `Kind.EndsStatement()` is true for exactly those kinds.
Consequently:

- A line ending in an operator, an assignment, `,`, `(`, `[`, `{`, `:`, `.` or
  `=>` continues on the next line.
- Several consecutive `Newline` tokens collapse into **one**.
- `;` also terminates a statement; see `grammar.md` for where a terminator is
  required and how a stray `;` is reported.
- If the file ends when the last emitted token can end a statement, a `Newline`
  is emitted **before** `EOF`.
- The rule is purely lexical and depends only on the previous token kind. It is
  **not** affected by what follows the line break: the lexer does **not**
  suppress a `Newline` because the next significant token is `{`.

### 7.1 A block-opening `{` on its own line

The lexer's rule is purely local and depends on the previous token kind only. A
`Newline` is emitted when an **actual** qualifying logical line break follows a
token that can end a statement:

- `fn main() {` — `)` is followed on the **same** line by `{`, so there is no
  qualifying line break and **no** `Newline` is emitted.
- `fn main()` followed by a line break and then `{` — the break is real, so
  exactly one `Newline` is emitted.
- `fn main(); {` — `;` is emitted, there is no line break before `{`, and no
  `Newline` is emitted.

Whether the emitted `Newline` may be skipped is a **parser** decision, ratified
as follows:

- The parser **may** consume exactly one `Newline` token when the next required
  token is a block-opening `{`, in a recognised syntactic header context: a
  function declaration header, and the headers of the compatible structured
  statements (`if`, `else`, `while`, `for`, `impl`). That list is closed and is
  specified in `grammar.md`.
- It **must not** consume a `Semicolon` as if it were optional whitespace.
- It **must not** ignore `Newline` tokens globally inside a block; there,
  `Newline` remains a terminator.
- It **must not** accept arbitrary statements between a header and its body. The
  only token allowed between the header and `{` is that single `Newline`, plus
  comments, which are not tokens.

Both layouts are therefore valid and equivalent:

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

The formatter's canonical output puts `{` on the header line; it does not reject
the second form as input.


## 8. Literals

### 8.1 Integer literals (`Int`)

Accepted: `42`, `1_000_000`, `0xFF`, `0b1010`, `0o755`. The value type is
`int64`.

- `_` is permitted **only between digits** of the same literal, so `1_000` and
  `0xFF_FF` are valid, while `1_`, `1__0`, `0x_FF` and `0x` are **E0103**. A
  leading `_` is not an error: it makes the sequence an identifier (§5), so
  `_1` is an `Ident`.
- There is **no** sign in the literal; `-` is always a prefix operator. Hence
  `-9223372036854775808` is unary minus applied to a literal that is itself out
  of range, and that is a documented consequence, not a special case to fix.
- Overflow of the literal's value is **E0103**.
- A decimal integer immediately followed by `.` and a digit is a float (§8.2),
  not an `Int` followed by a `Dot`. A `.` that does not meet that condition ends
  the integer and lexes as its own `Dot`, which is valid punctuation (§4, §8.2).

### 8.2 Floating-point literals (`Float`)

Accepted: `3.14`, `1e-3`, `2.5E10`, `1_000.5`. A **digit is required before and
after** the decimal point, so `.5` and `5.` are **not** float literals.

- Value semantics are **IEEE-754 binary64**.
- A well-formed literal whose value **overflows to ±Infinity**, or whose
  rounding **underflows** to zero, is **accepted**. Magnitude alone is never a
  lexical error.
- A malformed spelling is **E0103** only **once a numeric literal has begun**:
  `1e`, `1e+`, `0x`, `1_`, a misplaced `_`. A leading or standalone `.` never
  starts a literal, so it is never E0103 (see below).

**A leading dot is not a number, and it is not an error.** Scanning of a
numeric literal only ever begins at a decimal digit, so a `.` is left to the
punctuation rules of §4 unless it sits between digits, in which case the float
rule above applies:

| Source | Tokens | Lexical diagnostic |
|---|---|---|
| `1.5` | `Float`(`1.5`) | none |
| `.5` | `Dot`, `Int`(`5`) | **none** |
| `5.` | `Int`(`5`), `Dot` | **none** |
| `.` | `Dot` | **none** |
| `1..5` | `Int`(`1`), `DotDot`, `Int`(`5`) | none |
| `1e` | malformed `Int`-started literal | E0103 at `1` |

`.5` therefore fails as **an expression**, at the parser, with **E0201** — the
`Dot` is valid punctuation and expected an identifier after it. That separation
is deliberate: the lexer does not know or care where a `.` appears, so it can
never report E0103 for one. `.5` remains rejected as a float literal, as the
guide's 4.4 literal table requires, and no rule in this document permits it.

### 8.3 String literals (`String`)

Delimited by `"`. The lexer only delimits the literal and returns the raw text;
**the parser** decodes escapes and splits out interpolation parts.

- Escapes: `\n \t \r \\ \" \{ \} \0 \xNN \u{1F600}`. An unrecognised escape is
  **E0105**.
- Interpolation is `{expr}` inside the string. While scanning, the lexer tracks
  brace depth, and nested strings inside an interpolation are scanned
  recursively, so `"a {f("b")} c"` is a single token. The closing `"` ends the
  literal only at depth 0.
- An **unterminated** ordinary string — one that reaches a logical line break, or
  end of file, before the closing `"` at depth 0 — is **E0102**. Ordinary
  strings do not span lines.
- Interpolation boundaries stay inside the token's raw text. Producing the inner
  expression tokens is the parser's job, with a base offset so inner positions
  remain correct.

### 8.4 Raw string literals (`RawString`)

Delimited by backticks; **multiline, with no escapes and no interpolation**. The
token ends at the next backtick, and reaching end of file first is **E0102**.
Because there are no escapes, a backtick cannot appear inside a raw string.

Because this is the only literal that spans lines, it is where the §2 byte
invariant is most visible: a multiline raw string containing CRLF returns
exactly those two bytes per line, a raw string containing lone CR returns the one
CR byte, and `End.Offset - Pos.Offset` equals `len(Lexeme)` in every case. The
logical line count still advances once per break; only the *bookkeeping* is
normalised, never the bytes.

## 9. Source positions

```go
type Pos struct {
    Offset int // byte offset of the first byte of the token, 0-based
    Line   int // 1-based logical line
    Col    int // 1-based column in Unicode code points
}

type Token struct {
    Kind   Kind
    Lexeme string // raw source text of the token, unchanged
    Pos    Pos
    End    Pos    // position just past the token
}
```

- `End` is exclusive: it is the position of the byte after the token's last
  byte, which makes `src[Pos.Offset:End.Offset]` exactly the lexeme. This is the
  normative invariant of §2:

  ```
  Token.Lexeme == string(src[Pos.Offset:End.Offset])
  ```
- `Lexeme` is the **raw** text, before any `_` removal or escape decoding. A
  golden test compares it, so it must be **byte-faithful** — including any CR,
  CRLF or non-ASCII byte the source contains. It is never normalised.
- For a synthetic `Newline` inserted by ASI (§7), `Pos` and `End` are both the
  position where the line break occurs and `Lexeme` is empty. This is the **only**
  token excluded from the invariant above, because it has no source span of its
  own. A consumer that needs to know whether a `Newline` was explicit must
  therefore consult the source, not the token.
- CRLF increments `Line` once while `Offset` advances by two (§2), so the
  column of the first token on the next line is 1 in both LF and CRLF files.

## 10. Lexical diagnostics

The M1 subset of `internal/diag` is deliberately minimal: `Severity`, the
8-field `Diagnostic`, `Errorf`, and the E0101–E0105 constants. Rendering,
localisation, colour and the message catalogue are M6.

| Code | Condition | Position reported |
|---|---|---|
| `E0101` | invalid character | the offending code point |
| `E0102` | unterminated string or raw string literal | the opening quote or backtick |
| `E0103` | malformed number literal, or integer literal overflow; never a standalone `.` | the start of the literal |
| `E0104` | unterminated block comment | the opening `/*` |
| `E0105` | invalid escape sequence | the backslash |

Behaviour on a diagnostic:

- The lexer **never fails**: `Next()` always returns a token. On an
  unrecoverable character it emits `Illegal` with lexeme equal to the offending
  code point, records the diagnostic, and **continues at the next character**.
  Recovery must make progress, so a malformed token is never retried.
- Diagnostics are accumulated in source order and are not deduplicated; two
  identical errors at different positions are two diagnostics.
- A BOM after the start, or a `\r` that is not part of a line ending, is not
  special-cased into a diagnostic: CR is a line ending (§2), and an interior BOM
  is E0101 like any other invalid character.


## 11. Lexer API and package boundaries

The planned M1 surface:

```go
package lexer

func New(file string, src []byte) *Lexer
func (l *Lexer) Next() token.Token
func (l *Lexer) Diagnostics() []diag.Diagnostic
func (l *Lexer) Comments() []token.Comment
```

- `New` never reads the filesystem, the environment, the clock or the network.
- `Next` is total: it returns `EOF` forever once the input is exhausted.
- `Diagnostics()` and `Comments()` are read after the stream is consumed.

Dependency direction:

```
token  <-  diag  <-  lexer
```

`token` defines `Pos` and `Kind`. `diag` imports `token` for positions. `lexer`
imports both. **`token` must not import `diag`**, so the diagnostic model cannot
leak into the token package and a `Token` stays printable without it. There is
no cycle.

## 12. Token inspection command (future)

`atron tokens <file>` is **not implemented and not authorised** by this
document. Its contract is fixed here so that M1 tests and a later implementation
agree:

- exit `0` when tokenization succeeds, with the token stream on stdout;
- exit `2` when lexical diagnostics are reported;
- exit `64` for invalid command usage;
- diagnostics go to **stderr**, never stdout;
- no `main` is required in order to inspect a file (ADR-0005).

## 13. Conformance checklist for M1

The lexer is conformant when, at minimum, all of these hold:

| # | Requirement |
|---|---|
| 1 | `Kind` has exactly 62 members, and a test asserts the count and that the zero value is `Illegal`. |
| 2 | `Let` does not exist; `Const` does; the keyword group is 22. |
| 3 | `_1`, `_name`, `_`, `naïve`, `Ω`, `東京` lex as `Ident` with **no** diagnostic; a regression test asserts `_1` never yields E0103. |
| 4 | LF, CRLF and lone CR each produce exactly one logical line break; `Line`/`Col`/`Offset` agree with §2 and §9, including a CRLF in mid-file. |
| 5 | A lone `\r` produces no E0101. |
| 6 | A shebang line and a leading BOM are skipped; an interior BOM is E0101. |
| 7 | ASI fires on exactly the §7 set and collapses consecutive `Newline`s, emitting one before `EOF` when required. |
| 8 | A `Newline` is **not** suppressed when a real line break precedes `{`, and is **not** emitted when `{` follows on the same line; the parser-side skip is tested in M2, not here. |
| 9 | `.5` produces no lexical diagnostic and lexes as `Dot`, `Int`(`5`); `1e400` is a valid `Float` evaluating to +Infinity; `1_000_000` and `0xFFFF` are `Int`; `1_` is E0103. |
| 10 | `"a {f("b")} c"` is a single `String`; an unterminated `"` is E0102; a backticked raw string is `RawString`, and a multiline raw string is accepted. |
| 11 | For every real token, `Lexeme == string(src[Pos.Offset:End.Offset])`, including a multiline raw string written with CRLF and one containing a lone CR. The synthetic `Newline` is the only exception and has an empty lexeme. |
| 12 | Collected comments preserve their original source bytes exactly, including CR and CRLF inside a block comment. |
| 13 | `/* /* */` closes at the first `*/`; an unterminated `/*` is E0104. |
| 14 | Every comment is returned by `Comments()`; none is dropped. |
| 15 | `Next()` never panics and always makes progress on arbitrary bytes; `FuzzLexer` asserts this. |
| 16 | Every diagnostic carries a position that slices the source at the offending byte. |
| 17 | A standalone `.` lexes as exactly one `Dot` and produces **no** diagnostic; `..` is `DotDot`, `..=` is `DotDotEq`, and `...` is `DotDot` followed by `Dot`. |
| 18 | `a.b` yields `Ident`, `Dot`, `Ident`, so member access is available to the parser; `Dot` is never reported as E0103, and no lexical error is raised in a member-access position. |
| 19 | `.5` produces **no** lexical diagnostic and is **not** a `Float`: it lexes as `Dot`, `Int`(`5`). The rejection of `.5` as an expression is a parser E0201, asserted in M2, not here. |
| 20 | `5.` is `Int`(`5`) followed by `Dot`, and `1..5` is `Int`, `DotDot`, `Int` — a `.` never silently joins or splits a number beyond the §8.2 rule. |

## 14. What this document does not claim

- It does not claim any lexer code exists. **M1 has not started.**
- It does not define rendering, colours, message wording or localisation (M6).
- It does not define the parser, the AST, precedence or statement grammar; see
  [`grammar.md`](grammar.md).
- It does not reserve any diagnostic number above E0105.
