// SPDX-License-Identifier: Apache-2.0

// Package lexer converts Atron source text (UTF-8) into a token stream.
//
// It implements the byte and rune traversal machinery ratified in
// docs/spec/lexical.md: an immutable source buffer addressed by byte offset,
// UTF-8 decoding, rune-based columns, LF/CRLF/CR logical line breaks, Unicode
// identifier classification, exact keyword recognition, BOM and shebang
// handling, and deterministic recovery on malformed input.
//
// Source fidelity is a hard invariant. The lexer never rewrites, normalises or
// re-encodes the input, so for every real token
//
//	token.Lexeme == string(src[token.Pos.Offset:token.End.Offset])
//
// and one malformed byte is never replaced by a substitution character.
//
// Status: M1 work packages W0–W2 only. The cursor, identifier scanner and the
// New/Next/Diagnostics/Comments surface exist, but Next recognises only
// identifiers, keywords, whitespace, end of input and malformed UTF-8.
// Numbers, strings, operators, punctuation, comments and ASI are not
// implemented; see Next for exactly what is supported. This package is not a
// complete M1 lexer and M1 is not accepted.
package lexer
