// SPDX-License-Identifier: Apache-2.0

package token

// The complete kind inventory, in the order fixed by lexical.md §4.
//
// There are exactly 62 kinds: 4 special/identifier, 4 literal, 20 operator,
// 12 delimiter/punctuation and 22 keyword. The zero value is Illegal, and the
// block is written as one enumeration so the numeric values are contiguous from
// 0 to 61 and the count can be checked mechanically.
const (
	// Special and identifier kinds.

	// Illegal is the zero value and marks input the lexer could not classify.
	Illegal Kind = iota
	// EOF marks the end of the input.
	EOF
	// Newline is the synthetic line terminator inserted by ASI (lexical.md §7).
	// It is a real kind, not a parser convention, and it has an empty Lexeme.
	Newline
	// Ident is an identifier (lexical.md §5). `main` and `_1` are Ident.
	Ident

	// Literal kinds.

	// Int is an integer literal such as 42 or 0xFF.
	Int
	// Float is a floating-point literal such as 3.14 or 1e-3.
	Float
	// String is a double-quoted string literal; escapes are decoded by the
	// parser, not by the lexer.
	String
	// RawString is a backtick-delimited, multiline, escape-free literal.
	RawString

	// Operator kinds, recognised longest-first.

	// Plus is `+`.
	Plus
	// Minus is `-`.
	Minus
	// Star is `*`.
	Star
	// Slash is `/`.
	Slash
	// Percent is `%`.
	Percent
	// StarStar is `**`.
	StarStar
	// Eq is `==`.
	Eq
	// NotEq is `!=`.
	NotEq
	// Lt is `<`.
	Lt
	// LtEq is `<=`.
	LtEq
	// Gt is `>`.
	Gt
	// GtEq is `>=`.
	GtEq
	// Assign is `=`.
	Assign
	// PlusAssign is `+=`.
	PlusAssign
	// MinusAssign is `-=`.
	MinusAssign
	// StarAssign is `*=`.
	StarAssign
	// SlashAssign is `/=`.
	SlashAssign
	// PercentAssign is `%=`.
	PercentAssign
	// DotDot is `..`; it is distinct from DotDotEq and from Dot.
	DotDot
	// DotDotEq is `..=`.
	DotDotEq

	// Delimiter and punctuation kinds.

	// LParen is `(`.
	LParen
	// RParen is `)`.
	RParen
	// LBracket is `[`.
	LBracket
	// RBracket is `]`.
	RBracket
	// LBrace is `{`.
	LBrace
	// RBrace is `}`.
	RBrace
	// Comma is `,`.
	Comma
	// Colon is `:`.
	Colon
	// Semicolon is `;`.
	Semicolon
	// Dot is `.`. It is always valid punctuation and never a lexical error:
	// a standalone dot is never E0103 (lexical.md §4, §8.2).
	Dot
	// Question is `?`.
	Question
	// FatArrow is `=>`.
	FatArrow

	// Keyword kinds. `let` is not among them: it was replaced by `const`
	// (ADR-0004) and is an ordinary identifier. `main` is an identifier too
	// (ADR-0005).

	// Const is the `const` keyword, which requires an initializer.
	Const
	// Var is the `var` keyword, which declares a mutable binding.
	Var
	// Fn is the `fn` keyword.
	Fn
	// Return is the `return` keyword.
	Return
	// If is the `if` keyword.
	If
	// Else is the `else` keyword.
	Else
	// While is the `while` keyword.
	While
	// For is the `for` keyword.
	For
	// In is the `in` keyword.
	In
	// Break is the `break` keyword.
	Break
	// Continue is the `continue` keyword.
	Continue
	// Struct is the `struct` keyword.
	Struct
	// Impl is the `impl` keyword.
	Impl
	// Import is the `import` keyword.
	Import
	// Pub is the `pub` keyword.
	Pub
	// As is the `as` keyword.
	As
	// And is the `and` keyword.
	And
	// Or is the `or` keyword.
	Or
	// Not is the `not` keyword.
	Not
	// True is the `true` keyword; there is no separate boolean literal kind.
	True
	// False is the `false` keyword.
	False
	// Nil is the `nil` keyword.
	Nil
)
