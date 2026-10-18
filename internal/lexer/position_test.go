// SPDX-License-Identifier: Apache-2.0

package lexer_test

import (
	"testing"

	"github.com/CosmicLang/atron-lang/internal/lexer"
	"github.com/CosmicLang/atron-lang/internal/token"
)

// TestPositionsUseByteOffsetsAndRuneColumns pins the split between Offset and
// Col: offsets index the raw bytes, columns count Unicode code points.
func TestPositionsUseByteOffsetsAndRuneColumns(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		src     string
		wantPos token.Pos
		wantEnd token.Pos
	}{
		{
			name:    "ascii at start",
			src:     "abc",
			wantPos: token.Pos{Offset: 0, Line: 1, Col: 1},
			wantEnd: token.Pos{Offset: 3, Line: 1, Col: 4},
		},
		{
			name:    "after one ascii space",
			src:     " abc",
			wantPos: token.Pos{Offset: 1, Line: 1, Col: 2},
			wantEnd: token.Pos{Offset: 4, Line: 1, Col: 5},
		},
		{
			// "Ω" is two bytes and one rune, so the following identifier starts
			// at byte 3 but column 3.
			name:    "after a two-byte rune",
			src:     "Ω x",
			wantPos: token.Pos{Offset: 3, Line: 1, Col: 3},
			wantEnd: token.Pos{Offset: 4, Line: 1, Col: 4},
		},
		{
			// "東京" is six bytes and two runes, then a space, then "x".
			name:    "after three-byte runes",
			src:     "東京 x",
			wantPos: token.Pos{Offset: 7, Line: 1, Col: 4},
			wantEnd: token.Pos{Offset: 8, Line: 1, Col: 5},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			toks := collect(t, tc.src)
			// The last token before EOF is the identifier under test.
			last := toks[len(toks)-2]
			if last.Pos != tc.wantPos {
				t.Errorf("Pos = %+v, want %+v", last.Pos, tc.wantPos)
			}
			if last.End != tc.wantEnd {
				t.Errorf("End = %+v, want %+v", last.End, tc.wantEnd)
			}
		})
	}
}

// TestNonASCIIColumnsCountRunes walks a multi-byte line and checks that each
// column advances by one per rune while offsets advance by encoded width.
func TestNonASCIIColumnsCountRunes(t *testing.T) {
	t.Parallel()

	l := lexer.New("test.atrn", []byte("a Ω 東 z"))
	want := []struct {
		lexeme string
		offset int
		col    int
	}{
		{"a", 0, 1},
		{"Ω", 2, 3},
		{"東", 5, 5},
		{"z", 9, 7},
	}
	for i, w := range want {
		tok := l.Next()
		if tok.Kind != token.Ident {
			t.Fatalf("token %d = %v, want Ident", i, tok.Kind)
		}
		if tok.Lexeme != w.lexeme {
			t.Fatalf("token %d lexeme = %q, want %q", i, tok.Lexeme, w.lexeme)
		}
		if tok.Pos.Offset != w.offset {
			t.Errorf("%q Pos.Offset = %d, want %d", w.lexeme, tok.Pos.Offset, w.offset)
		}
		if tok.Pos.Col != w.col {
			t.Errorf("%q Pos.Col = %d, want %d", w.lexeme, tok.Pos.Col, w.col)
		}
	}
}

// TestLineEndingsEachCountOneBreak checks that LF, CRLF and a standalone CR are
// each exactly one logical line break, with the byte advance the contract
// specifies, and that no diagnostic is produced for a lone CR.
func TestLineEndingsEachCountOneBreak(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		src        string
		wantSecond token.Pos
	}{
		{"LF", "a\nb", token.Pos{Offset: 2, Line: 2, Col: 1}},
		{"CRLF", "a\r\nb", token.Pos{Offset: 3, Line: 2, Col: 1}},
		{"CR", "a\rb", token.Pos{Offset: 2, Line: 2, Col: 1}},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l := lexer.New("test.atrn", []byte(tc.src))

			first := l.Next()
			if first.Kind != token.Ident || first.Lexeme != "a" {
				t.Fatalf("first token = %v %q", first.Kind, first.Lexeme)
			}

			second := l.Next()
			if second.Kind != token.Ident || second.Lexeme != "b" {
				t.Fatalf("second token = %v %q", second.Kind, second.Lexeme)
			}
			if second.Pos != tc.wantSecond {
				t.Errorf("second token at %+v, want %+v", second.Pos, tc.wantSecond)
			}
			if d := l.Diagnostics(); len(d) != 0 {
				t.Errorf("%s produced diagnostics %v, want none (a lone CR is a line break, not E0101)",
					tc.name, d)
			}
		})
	}
}

// TestCRLFAdvancesTwoBytesOneLine is the explicit contract check: CRLF consumes
// two bytes and increments Line exactly once.
func TestCRLFAdvancesTwoBytesOneLine(t *testing.T) {
	t.Parallel()

	l := lexer.New("test.atrn", []byte("a\r\nb"))
	first := l.Next()
	second := l.Next()

	if got := second.Pos.Offset - first.End.Offset; got != 2 {
		t.Errorf("CRLF advanced %d bytes, want 2", got)
	}
	if got := second.Pos.Line - first.Pos.Line; got != 1 {
		t.Errorf("CRLF advanced %d lines, want 1", got)
	}
	if second.Pos.Col != 1 {
		t.Errorf("column after CRLF = %d, want 1", second.Pos.Col)
	}
}

// TestCRLFMidFileCountsOneBreak checks a CRLF that is neither at the start nor
// at the end of the file.
func TestCRLFMidFileCountsOneBreak(t *testing.T) {
	t.Parallel()

	l := lexer.New("test.atrn", []byte("aa\r\nbb\r\ncc"))
	want := []token.Pos{
		{Offset: 0, Line: 1, Col: 1},
		{Offset: 4, Line: 2, Col: 1},
		{Offset: 8, Line: 3, Col: 1},
	}
	for i, w := range want {
		tok := l.Next()
		if tok.Kind != token.Ident {
			t.Fatalf("token %d = %v, want Ident", i, tok.Kind)
		}
		if tok.Pos != w {
			t.Errorf("token %d at %+v, want %+v", i, tok.Pos, w)
		}
	}
}

// TestTabAdvancesOneColumn pins that a tab is one code point, not a tab stop.
func TestTabAdvancesOneColumn(t *testing.T) {
	t.Parallel()

	l := lexer.New("test.atrn", []byte("\tx"))
	tok := l.Next()
	if tok.Pos.Col != 2 {
		t.Fatalf("column after a tab = %d, want 2", tok.Pos.Col)
	}
}

// TestWhitespaceProducesNoTokens covers the four separating whitespace
// characters.
func TestWhitespaceProducesNoTokens(t *testing.T) {
	t.Parallel()

	for _, ws := range []string{" ", "\t", "\f", "\v", " \t\f\v"} {
		toks := collect(t, ws+"x"+ws)
		if len(toks) != 2 {
			t.Errorf("whitespace %q produced %d tokens, want 2", ws, len(toks))
			continue
		}
		if toks[0].Kind != token.Ident || toks[0].Lexeme != "x" {
			t.Errorf("whitespace %q: token = %v %q, want Ident x", ws, toks[0].Kind, toks[0].Lexeme)
		}
	}
}
