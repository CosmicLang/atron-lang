// SPDX-License-Identifier: Apache-2.0

package lexer

import (
	"github.com/CosmicLang/atron-lang/internal/diag"
	"github.com/CosmicLang/atron-lang/internal/token"
)

// Next returns the next token.
//
// Next is total: it never panics and it returns EOF forever once the input is
// exhausted. It also always makes forward progress on arbitrary bytes, so a
// loop that consumes tokens terminates even for malformed input.
//
// # Scope in work packages W0–W2
//
// W2 covers the traversal machinery, not the whole scanner. Next currently
// recognises exactly these families:
//
//   - identifiers and keywords (token.Ident, or the 22 keyword kinds), including
//     Unicode identifiers, with exact case-sensitive keyword matching;
//   - whitespace and logical line breaks, which separate tokens and produce
//     nothing;
//   - end of input, reported as token.EOF;
//   - malformed UTF-8, reported as token.Illegal with one E0101 diagnostic at
//     the offending byte;
//   - an interior U+FEFF, reported as token.Illegal with one E0101 diagnostic,
//     because a byte order mark is never part of any Atron token.
//
// Every other family — numbers, strings, raw strings, operators, punctuation,
// comments and automatic semicolon insertion — is deliberately not implemented
// yet. Such input yields a token.Illegal whose Lexeme is the single source rune
// and, crucially, **no diagnostic**: the lexer must never invent a normative
// lexical error (E0101–E0105) for syntax this work package does not implement.
// A later work package replaces these placeholders with real tokens and, where
// the contract requires it, real diagnostics.
//
// Callers can therefore tell the two Illegal cases apart: a ratified lexical
// error is Illegal *and* accompanied by a diagnostic, while a
// not-yet-implemented family is Illegal with no diagnostic. Next is not a
// complete M1 scanner and must not be described as one.
func (l *Lexer) Next() token.Token {
	if l.done {
		p := l.pos()
		return token.Token{Kind: token.EOF, Pos: p, End: p}
	}

	l.skipSpace()

	if l.atEOF() {
		l.done = true
		p := l.pos()
		return token.Token{Kind: token.EOF, Pos: p, End: p}
	}

	if l.invalidUTF8() {
		start := l.pos()
		tok := l.scanInvalidRune()
		l.report(diag.E0101, start, tok.End, "invalid UTF-8 byte sequence")
		return tok
	}

	if l.identStart() {
		return l.scanIdentifier()
	}

	// A U+FEFF that was not skipped as the leading BOM is invalid everywhere
	// else, and that verdict is independent of which token families exist: the
	// character is never part of any Atron token, so E0101 is the ratified
	// answer rather than a guess (lexical.md §1, §10).
	if r, _, ok := l.peek(); ok && r == '\uFEFF' {
		start := l.pos()
		tok := l.scanUnsupported()
		l.report(diag.E0101, start, tok.End, "invalid character U+FEFF (byte order mark)")
		return tok
	}

	return l.scanUnsupported()
}

// scanUnsupported consumes one rune of a token family that W2 does not
// implement and returns a placeholder Illegal token for it.
//
// The rune is published verbatim as the lexeme, so the source-fidelity
// invariant still holds and recovery is exact. No diagnostic is recorded: an
// unimplemented family is a known limitation of this work package, not a
// normative lexical error, and reporting one would fabricate a false finding
// about valid Atron source.
func (l *Lexer) scanUnsupported() token.Token {
	start := l.pos()
	startOff := l.off

	if !l.atEOF() {
		l.advance()
	}

	return token.Token{
		Kind:   token.Illegal,
		Lexeme: string(l.src[startOff:l.off]),
		Pos:    start,
		End:    l.pos(),
	}
}
