// SPDX-License-Identifier: Apache-2.0

package lexer

import (
	"github.com/CosmicLang/atron-lang/internal/diag"
	"github.com/CosmicLang/atron-lang/internal/token"
)

// Lexer turns Atron source text into a stream of tokens.
//
// A Lexer owns a reference to the caller's byte slice and treats it as
// immutable: it never writes to it, never takes a normalised copy and never
// caches rewritten text. Every Lexeme is therefore a slice of exactly the bytes
// that were supplied, which is what makes the contract invariant
//
//	token.Lexeme == string(src[token.Pos.Offset:token.End.Offset])
//
// hold for every real, source-backed token.
//
// A Lexer is not safe for concurrent use. Diagnostics and comments are
// accumulated during scanning and read afterwards with Diagnostics and
// Comments.
type Lexer struct {
	// file is the name reported in diagnostics. It is stored verbatim and is
	// never opened, resolved or validated against the filesystem.
	file string
	// src is the caller's input, treated as immutable.
	src []byte

	// off, line and col are the cursor: a byte offset into src plus the logical
	// line and rune-based column it addresses.
	off  int
	line int
	col  int

	// diags holds diagnostics in the order they were produced. They are not
	// deduplicated: two identical problems at different offsets are two
	// diagnostics.
	diags []diag.Diagnostic

	// done records that EOF has been returned, after which Next keeps returning
	// EOF without consulting the input again.
	done bool
}

// New returns a Lexer over src, reporting diagnostics against the given file
// name.
//
// New performs the whole-file preamble once: it skips a leading UTF-8 BOM and
// then a shebang line, in that order, because a BOM may precede the shebang
// (lexical.md §1). Neither is part of any token.
//
// New never reads the filesystem, the environment, the clock or the network. It
// does not copy src, so the caller must not modify the slice while the Lexer is
// in use.
func New(file string, src []byte) *Lexer {
	l := &Lexer{
		file: file,
		src:  src,
		line: 1,
		col:  1,
	}
	l.skipBOM()
	l.skipShebang()
	return l
}

// File returns the file name used in diagnostics.
func (l *Lexer) File() string { return l.file }

// Diagnostics returns the diagnostics produced so far, in source order.
//
// The result is a copy, so a caller cannot append to or reorder the lexer's
// internal slice. Diagnostics are complete only after the stream has been
// consumed.
func (l *Lexer) Diagnostics() []diag.Diagnostic {
	if len(l.diags) == 0 {
		return nil
	}
	out := make([]diag.Diagnostic, len(l.diags))
	copy(out, l.diags)
	return out
}

// Comments returns the comment metadata collected so far, in source order.
//
// Comment collection is not part of work packages W0–W2, so this always returns
// nil for now. The accessor carries its final signature so callers can be
// written against the ratified API. When collection lands, the returned slice
// must be a copy, so callers cannot mutate the lexer's records.
func (*Lexer) Comments() []token.Comment {
	return nil
}

// report records a diagnostic at the given span. It is the only place
// diagnostics are appended, so ordering is exactly production order.
func (l *Lexer) report(code string, pos, end token.Pos, format string, args ...any) {
	l.diags = append(l.diags, diag.Errorf(code, l.file, pos, end, format, args...))
}
