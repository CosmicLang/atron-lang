// SPDX-License-Identifier: Apache-2.0

package token

import "strconv"

// Kind identifies the lexical class of a token.
//
// Kind is a uint8-backed enumeration whose zero value is Illegal, so a
// zero-valued Token is always recognisably invalid. The complete inventory is
// fixed by the ratified lexical contract (lexical.md §4) and is declared in
// kinds.go: exactly 62 kinds.
type Kind uint8

// String returns the stable name of the kind, or "Kind(n)" for a value outside
// the enumeration. It never panics, which keeps token formatting safe on tokens
// produced by a buggy or future scanner.
func (k Kind) String() string {
	if int(k) < len(kindNames) && kindNames[k] != "" {
		return kindNames[k]
	}
	return "Kind(" + strconv.Itoa(int(k)) + ")"
}

// Token is a single lexical token.
//
// Lexeme is the raw source text of the token, unchanged: no underscore removal,
// no escape decoding and no normalisation. For every real, source-backed token
//
//	Lexeme == string(src[Pos.Offset:End.Offset])
//
// holds. End is exclusive, which is what makes that slicing exact. The only
// exception is the synthetic Newline inserted by automatic semicolon insertion
// (lexical.md §7), which has an empty Lexeme and no source span of its own.
type Token struct {
	// Kind is the lexical class of the token.
	Kind Kind
	// Lexeme is the raw source text, verbatim.
	Lexeme string
	// Pos is the position of the first byte of the token.
	Pos Pos
	// End is the position just past the last byte of the token, exclusive.
	End Pos
}

// Span returns the byte range the token occupies, as src[start:end].
//
// For a synthetic Newline the range is empty and both bounds equal Pos.Offset,
// which is exactly what the lexical contract requires of a token with no source
// span of its own.
func (t Token) Span() (start, end int) {
	return t.Pos.Offset, t.End.Offset
}

// Comment is the lexical metadata retained for one source comment.
//
// Comments are not tokens: there is no Comment kind, because comments never
// enter the parser's token stream (lexical.md §4). They are kept separately so
// the formatter can reprint them exactly, which is why Text holds the original
// source bytes including the delimiters and any CR or CRLF inside a block
// comment.
type Comment struct {
	// Text is the raw source text of the comment, delimiters included.
	Text string
	// Pos is the position of the opening delimiter.
	Pos Pos
	// End is the position just past the closing delimiter, or the end of the
	// line for a line comment, exclusive.
	End Pos
	// Block reports whether this is a /* ... */ comment rather than a //
	// line comment.
	Block bool
}
