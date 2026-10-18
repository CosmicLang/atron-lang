// SPDX-License-Identifier: Apache-2.0

package token

import "strconv"

// Pos is a position in Atron source text.
//
// Offset is the 0-based byte offset of the first byte of the token, so that
// src[Pos.Offset:End.Offset] is exactly the token's raw text. Line is the
// 1-based logical line and Col is the 1-based column counted in Unicode code
// points, never bytes.
//
// LF, CRLF and a standalone CR each advance Line by exactly one and reset Col
// to 1, while Offset advances by the true number of bytes the source contains
// (lexical.md §2, §9). Normalisation is therefore bookkeeping only: no byte of
// the source is ever rewritten.
type Pos struct {
	// Offset is the byte offset of the first byte of the token, 0-based.
	Offset int
	// Line is the logical line, 1-based.
	Line int
	// Col is the column counted in Unicode code points, 1-based.
	Col int
}

// String renders the position as "line:column" for diagnostics and tests.
//
// Offset is deliberately omitted: it addresses the source buffer for slicing
// and is not a location a reader of a diagnostic recognises.
func (p Pos) String() string {
	return strconv.Itoa(p.Line) + ":" + strconv.Itoa(p.Col)
}
