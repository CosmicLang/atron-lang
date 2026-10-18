// SPDX-License-Identifier: Apache-2.0

package lexer

import (
	"unicode"
	"unicode/utf8"

	"github.com/CosmicLang/atron-lang/internal/token"
)

// isIdentStart reports whether r may begin an identifier.
//
// The ratified grammar is [\p{L}_][\p{L}\p{N}_]*. Only Unicode letters and the
// underscore start an identifier, so a digit-initial sequence is a number and
// `_1` is an identifier rather than a malformed number (lexical.md §5).
func isIdentStart(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

// isIdentContinue reports whether r may continue an identifier.
//
// The ratified continuation class is \p{N} — the whole Unicode number
// category, which spans Nd (decimal digits), Nl (letter-like numbers such as
// the Roman numerals) and No (other numbers such as superscripts and vulgar
// fractions). unicode.IsNumber is the predicate for that category;
// unicode.IsDigit would accept only Nd and would therefore split `valor²` and
// `itemⅫ` into an identifier plus a stray rune (lexical.md §5).
//
// Marks, punctuation, symbols and emoji are not in \p{N}, so a combining mark
// after a base letter ends the token and is then reported as invalid by the
// caller.
func isIdentContinue(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsNumber(r)
}

// isDigit reports whether r is an ASCII decimal digit.
//
// Number literals are ASCII-only, so this predicate is deliberately narrower
// than the identifier continuation class: '²' and 'Ⅻ' continue an identifier
// because they are in \p{N}, but they never begin a number literal. W2 does not
// scan numbers, so this predicate and digitStart exist as the cursor-level
// classification the numeric work package will build on; they are exercised by
// the in-package cursor tests.
func isDigit(r rune) bool { return r >= '0' && r <= '9' }

// scanIdentifier consumes an identifier starting at the cursor and returns it.
//
// The scan uses peek to test continuation without consuming, so the cursor stops
// exactly at the first byte that cannot continue the identifier; Line and Col
// therefore need no rewinding. The returned token's Lexeme is a slice of the
// untouched source, so it is byte-faithful by construction, and its Kind comes
// from the exact, case-sensitive keyword lookup.
func (l *Lexer) scanIdentifier() token.Token {
	start := l.pos()
	startOff := l.off

	for !l.atEOF() {
		r, _, _ := l.peek()
		if !isIdentContinue(r) {
			break
		}
		l.advance()
	}

	text := string(l.src[startOff:l.off])
	return token.Token{
		Kind:   token.LookupIdent(text),
		Lexeme: text,
		Pos:    start,
		End:    l.pos(),
	}
}

// scanInvalidRune consumes exactly one malformed input unit and returns the
// Illegal token that represents it.
//
// A malformed UTF-8 sequence is invalid input, not text to be repaired
// (lexical.md §1, §10). The offending unit is therefore republished byte for
// byte as the token's lexeme, with no replacement character substituted, so
// Lexeme still equals src[Pos.Offset:End.Offset]. Recovery consumes exactly the
// bytes reported by the decoder, which guarantees forward progress.
func (l *Lexer) scanInvalidRune() token.Token {
	start := l.pos()
	startOff := l.off

	_, size, ok := l.peek()
	if !ok {
		return token.Token{Kind: token.EOF, Pos: start, End: start}
	}
	for i := 0; i < size && l.off < len(l.src); i++ {
		l.off++
		l.col++
	}

	raw := string(l.src[startOff:l.off])
	return token.Token{Kind: token.Illegal, Lexeme: raw, Pos: start, End: l.pos()}
}

// invalidUTF8 reports whether the byte at the cursor begins a malformed UTF-8
// sequence.
//
// utf8.DecodeRune returns RuneError with size 1 only for malformed input; a
// genuine U+FFFD in the source decodes with size 3 and is valid text. This
// predicate is the single place that decision is made.
func (l *Lexer) invalidUTF8() bool {
	if l.atEOF() {
		return false
	}
	r, size := utf8.DecodeRune(l.src[l.off:])
	return r == utf8.RuneError && size == 1
}

// identStart reports whether the cursor sits on the start of an identifier.
func (l *Lexer) identStart() bool {
	r, _, ok := l.peek()
	return ok && isIdentStart(r)
}

// digitStart reports whether the cursor sits on an ASCII decimal digit.
func (l *Lexer) digitStart() bool {
	r, _, ok := l.peek()
	return ok && isDigit(r)
}
