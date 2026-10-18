// SPDX-License-Identifier: Apache-2.0

package lexer_test

import (
	"testing"

	"github.com/CosmicLang/atron-lang/internal/lexer"
	"github.com/CosmicLang/atron-lang/internal/token"
)

// TestEmptySourceReturnsEOFatOneOne pins the EOF position for empty input.
func TestEmptySourceReturnsEOFatOneOne(t *testing.T) {
	t.Parallel()

	for _, src := range []string{"", "\n", "\r\n"} {
		l := lexer.New("empty.atrn", []byte(src))
		tok := l.Next()
		if tok.Kind != token.EOF {
			t.Fatalf("%q: token = %v, want EOF", src, tok.Kind)
		}
		if tok.Pos != tok.End {
			t.Fatalf("%q: EOF span is not empty: %+v..%+v", src, tok.Pos, tok.End)
		}
		if tok.Lexeme != "" {
			t.Fatalf("%q: EOF lexeme = %q, want empty", src, tok.Lexeme)
		}
	}
}

// TestEOFPositionAfterTrailingText pins the EOF position after content,
// including a trailing line break.
func TestEOFPositionAfterTrailingText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		src  string
		want token.Pos
	}{
		{"abc", token.Pos{Offset: 3, Line: 1, Col: 4}},
		{"abc\n", token.Pos{Offset: 4, Line: 2, Col: 1}},
		{"abc\r\n", token.Pos{Offset: 5, Line: 2, Col: 1}},
		{"abc\r", token.Pos{Offset: 4, Line: 2, Col: 1}},
	}

	for _, tc := range tests {
		toks := collect(t, tc.src)
		eof := toks[len(toks)-1]
		if eof.Kind != token.EOF {
			t.Fatalf("%q: last token = %v, want EOF", tc.src, eof.Kind)
		}
		if eof.Pos != tc.want {
			t.Errorf("%q: EOF at %+v, want %+v", tc.src, eof.Pos, tc.want)
		}
		if eof.Pos != eof.End {
			t.Errorf("%q: EOF span is not empty: %+v..%+v", tc.src, eof.Pos, eof.End)
		}
	}
}

// TestEOFForever pins that Next keeps returning EOF without moving.
func TestEOFForever(t *testing.T) {
	t.Parallel()

	l := lexer.New("test.atrn", []byte("x"))
	first := l.Next()
	if first.Kind != token.Ident {
		t.Fatalf("first token = %v", first.Kind)
	}
	for i := 0; i < 5; i++ {
		tok := l.Next()
		if tok.Kind != token.EOF {
			t.Fatalf("call %d = %v, want EOF", i, tok.Kind)
		}
		if tok.Pos.Offset != 1 {
			t.Fatalf("call %d EOF offset = %d, want 1", i, tok.Pos.Offset)
		}
	}
	if d := l.Diagnostics(); len(d) != 0 {
		t.Fatalf("diagnostics on exhausted input: %v", d)
	}
}

// TestLeadingBOMIsSkipped checks that the BOM is not part of any token and does
// not shift the first column, and that no diagnostic is produced.
func TestLeadingBOMIsSkipped(t *testing.T) {
	t.Parallel()

	src := "\xef\xbb\xbfmain"
	l := lexer.New("test.atrn", []byte(src))
	tok := l.Next()

	if tok.Kind != token.Ident {
		t.Fatalf("token = %v, want Ident", tok.Kind)
	}
	if tok.Lexeme != "main" {
		t.Fatalf("lexeme = %q, want \"main\" with no BOM bytes", tok.Lexeme)
	}
	if tok.Pos != (token.Pos{Offset: 3, Line: 1, Col: 1}) {
		t.Fatalf("Pos = %+v, want offset 3, line 1, column 1", tok.Pos)
	}
	if d := l.Diagnostics(); len(d) != 0 {
		t.Fatalf("leading BOM produced diagnostics %v, want none", d)
	}
}

// TestBOMOnlySourceIsEmpty checks that a BOM alone leaves only EOF.
func TestBOMOnlySourceIsEmpty(t *testing.T) {
	t.Parallel()

	l := lexer.New("test.atrn", []byte("\xef\xbb\xbf"))
	tok := l.Next()
	if tok.Kind != token.EOF {
		t.Fatalf("token = %v, want EOF", tok.Kind)
	}
	if tok.Pos.Offset != 3 {
		t.Fatalf("EOF offset = %d, want 3", tok.Pos.Offset)
	}
	if d := l.Diagnostics(); len(d) != 0 {
		t.Fatalf("BOM-only source produced diagnostics %v", d)
	}
}

// TestInteriorBOMIsE0101 covers the ratified rule that a BOM after the start is
// an invalid character at its own position.
func TestInteriorBOMIsE0101(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		src    string
		offset int
		col    int
	}{
		{"after a token", "a\xef\xbb\xbfb", 1, 2},
		{"after a space", "a \xef\xbb\xbfb", 2, 3},
		{"after a line break", "a\n\xef\xbb\xbfb", 2, 1},
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

			bad := l.Next()
			if bad.Kind != token.Illegal {
				t.Fatalf("BOM token = %v, want Illegal", bad.Kind)
			}
			if bad.Lexeme != "\ufeff" {
				t.Fatalf("BOM lexeme = %q, want U+FEFF", bad.Lexeme)
			}
			if bad.Pos.Offset != tc.offset || bad.Pos.Col != tc.col {
				t.Fatalf("BOM at %+v, want offset %d column %d", bad.Pos, tc.offset, tc.col)
			}

			diags := l.Diagnostics()
			if len(diags) != 1 {
				t.Fatalf("got %d diagnostics, want exactly 1: %v", len(diags), diags)
			}
			if diags[0].Code != "E0101" {
				t.Fatalf("code = %q, want E0101", diags[0].Code)
			}
			if diags[0].Pos != bad.Pos || diags[0].End != bad.End {
				t.Fatalf("diagnostic span %+v..%+v, want the offending code point %+v..%+v",
					diags[0].Pos, diags[0].End, bad.Pos, bad.End)
			}

			// Recovery must continue normally past the BOM.
			last := l.Next()
			if last.Kind != token.Ident || last.Lexeme != "b" {
				t.Fatalf("token after the BOM = %v %q, want Ident b", last.Kind, last.Lexeme)
			}
		})
	}
}
