// SPDX-License-Identifier: Apache-2.0

package lexer

import (
	"unicode/utf8"

	"github.com/CosmicLang/atron-lang/internal/token"
)

// bomBytes is the UTF-8 encoding of U+FEFF, which is only legal at the very
// start of a source file (lexical.md §1).
const bomBytes = "\xef\xbb\xbf"

// atEOF reports whether every byte of the source has been consumed.
func (l *Lexer) atEOF() bool { return l.off >= len(l.src) }

// pos returns the current position. The zero value of the cursor corresponds to
// offset 0, line 1, column 1, which New establishes explicitly.
func (l *Lexer) pos() token.Pos {
	return token.Pos{Offset: l.off, Line: l.line, Col: l.col}
}

// peek decodes the rune at the cursor without consuming it.
//
// When the cursor is at end of input it returns ok == false. Invalid UTF-8 is
// reported as (utf8.RuneError, 1, true): a size of 1 for RuneError is what
// distinguishes a malformed byte from a correctly encoded U+FFFD, which always
// occupies three bytes. That distinction keeps recovery deterministic without
// ever substituting replacement characters into the source.
func (l *Lexer) peek() (r rune, size int, ok bool) {
	if l.atEOF() {
		return 0, 0, false
	}
	r, size = utf8.DecodeRune(l.src[l.off:])
	return r, size, true
}

// peekByteAt returns the byte at the given offset from the cursor, if any. It is
// used for two-byte lookahead that never crosses a UTF-8 sequence boundary,
// such as distinguishing CRLF from a standalone CR.
func (l *Lexer) peekByteAt(delta int) (b byte, ok bool) {
	i := l.off + delta
	if i < 0 || i >= len(l.src) {
		return 0, false
	}
	return l.src[i], true
}

// advance consumes one logical input unit and returns it.
//
// A unit is one rune, except that CRLF is consumed as a single logical line
// break advancing two bytes while incrementing Line only once (lexical.md §2).
// Columns count runes and are reset to 1 after any line break; a tab advances
// the column by one because tab stops are a display concern.
//
// advance always moves the cursor forward, which is what guarantees that Next
// terminates on arbitrary input.
func (l *Lexer) advance() (r rune, size int, ok bool) {
	r, size, ok = l.peek()
	if !ok {
		return 0, 0, false
	}

	switch r {
	case '\r':
		// CRLF and a standalone CR are each exactly one logical line break.
		if next, has := l.peekByteAt(1); has && next == '\n' {
			l.off += 2
		} else {
			l.off++
		}
		l.line++
		l.col = 1
	case '\n':
		l.off++
		l.line++
		l.col = 1
	default:
		// The offset advances by the true encoded width, so every offset stays a
		// valid index into the original source; only the column counts runes.
		l.off += size
		l.col++
	}
	return r, size, true
}

// skipBOM consumes a UTF-8 BOM at offset 0.
//
// The BOM is not part of any token, so it is not counted as a column: the first
// token of the file starts at column 1. A BOM anywhere else is not consumed
// here; it is reported as E0101 at its own position.
func (l *Lexer) skipBOM() {
	if len(l.src) >= len(bomBytes) && string(l.src[:len(bomBytes)]) == bomBytes {
		l.off += len(bomBytes)
	}
}

// skipShebang consumes a shebang line, including its line break.
//
// A shebang is recognised only at the logical start of the file, which is after
// a leading BOM has been skipped. It is not a comment and is not preserved as
// lexical metadata (lexical.md §1), so nothing is recorded for it; only the
// cursor moves. A shebang without a final line break runs to end of input.
func (l *Lexer) skipShebang() {
	if b, ok := l.peekByteAt(0); !ok || b != '#' {
		return
	}
	if b, ok := l.peekByteAt(1); !ok || b != '!' {
		return
	}
	for !l.atEOF() {
		r, _, _ := l.advance()
		if r == '\n' || r == '\r' {
			return
		}
	}
}

// skipSpace consumes whitespace and logical line breaks.
//
// Whitespace separates tokens and produces nothing (lexical.md §3). Line breaks
// are consumed here and only advance the position: the synthetic Newline token
// belongs to automatic semicolon insertion, which this work package does not
// implement, so no token is emitted for them.
func (l *Lexer) skipSpace() {
	for !l.atEOF() {
		r, _, _ := l.peek()
		if isSpace(r) || isLineBreak(r) {
			l.advance()
			continue
		}
		return
	}
}

// isSpace reports whether r is whitespace that separates tokens.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\f' || r == '\v'
}

// isLineBreak reports whether r begins a logical line break.
//
// Three representations are recognised: LF, CRLF and a standalone CR. A lone CR
// is a line break, not an invalid character, so it never produces E0101
// (lexical.md §2).
func isLineBreak(r rune) bool { return r == '\n' || r == '\r' }
