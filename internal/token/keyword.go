// SPDX-License-Identifier: Apache-2.0

package token

// keywords maps every keyword spelling to its kind.
//
// The map is the single source of truth for keyword recognition and is
// deliberately exact and case-sensitive: every keyword is ASCII, and `let`,
// `main`, `_1` and `Name` are all absent so they lex as Ident (ADR-0004,
// ADR-0005, lexical.md §5). There are exactly 22 entries and the tests assert
// that count by enumerating the keyword range, not by trusting this comment.
var keywords = map[string]Kind{
	"const":    Const,
	"var":      Var,
	"fn":       Fn,
	"return":   Return,
	"if":       If,
	"else":     Else,
	"while":    While,
	"for":      For,
	"in":       In,
	"break":    Break,
	"continue": Continue,
	"struct":   Struct,
	"impl":     Impl,
	"import":   Import,
	"pub":      Pub,
	"as":       As,
	"and":      And,
	"or":       Or,
	"not":      Not,
	"true":     True,
	"false":    False,
	"nil":      Nil,
}

// LookupIdent classifies an identifier's text.
//
// It returns the keyword kind when text is exactly a keyword spelling, and
// Ident otherwise. The comparison is exact and case-sensitive, so "if" is the
// If keyword while "If" and "IF" are ordinary identifiers.
func LookupIdent(text string) Kind {
	if k, ok := keywords[text]; ok {
		return k
	}
	return Ident
}

// KeywordCount is the number of keywords in the ratified inventory.
const KeywordCount = 22

// Keywords returns every keyword spelling and its kind, as a fresh map.
//
// The result is a copy, so callers such as the `atron tokens` documentation and
// completion generators cannot mutate the lexical contract by accident. It is
// used by the tests to enumerate the inventory mechanically.
func Keywords() map[string]Kind {
	out := make(map[string]Kind, len(keywords))
	for spelling, kind := range keywords {
		out[spelling] = kind
	}
	return out
}
