// SPDX-License-Identifier: Apache-2.0

package lexer_test

import (
	"testing"

	"github.com/CosmicLang/atron-lang/internal/lexer"
	"github.com/CosmicLang/atron-lang/internal/token"
)

// TestShebangIsIgnored checks that a shebang line is skipped, including its line
// break, and leaves no token and no metadata behind.
func TestShebangIsIgnored(t *testing.T) {
	t.Parallel()

	l := lexer.New("script.atrn", []byte("#!/usr/bin/env atron\nmain"))
	tok := l.Next()

	if tok.Kind != token.Ident || tok.Lexeme != "main" {
		t.Fatalf("first token = %v %q, want Ident main", tok.Kind, tok.Lexeme)
	}
	if tok.Pos != (token.Pos{Offset: 21, Line: 2, Col: 1}) {
		t.Fatalf("Pos = %+v, want offset 21, line 2, column 1", tok.Pos)
	}
	if d := l.Diagnostics(); len(d) != 0 {
		t.Fatalf("shebang produced diagnostics %v", d)
	}
	// The shebang is not a comment, so comment metadata must stay empty.
	if c := l.Comments(); len(c) != 0 {
		t.Fatalf("shebang produced %d comment records, want 0", len(c))
	}
}

// TestShebangVariants covers CRLF, a shebang with no line break, and a BOM
// before the shebang.
func TestShebangVariants(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		src     string
		wantTok string
		wantOff int
		wantLn  int
	}{
		{"CRLF terminator", "#!atron\r\nx", "x", 9, 2},
		{"standalone CR terminator", "#!atron\rx", "x", 8, 2},
		{"no terminator", "#!atron", "", 7, 1},
		{"BOM then shebang", "\xef\xbb\xbf#!atron\nx", "x", 11, 2},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l := lexer.New("test.atrn", []byte(tc.src))
			tok := l.Next()

			if tc.wantTok == "" {
				if tok.Kind != token.EOF {
					t.Fatalf("token = %v %q, want EOF", tok.Kind, tok.Lexeme)
				}
			} else if tok.Kind != token.Ident || tok.Lexeme != tc.wantTok {
				t.Fatalf("token = %v %q, want Ident %q", tok.Kind, tok.Lexeme, tc.wantTok)
			}
			if tok.Pos.Offset != tc.wantOff {
				t.Errorf("offset = %d, want %d", tok.Pos.Offset, tc.wantOff)
			}
			if tok.Pos.Line != tc.wantLn {
				t.Errorf("line = %d, want %d", tok.Pos.Line, tc.wantLn)
			}
			if d := l.Diagnostics(); len(d) != 0 {
				t.Errorf("diagnostics = %v, want none", d)
			}
		})
	}
}

// TestHashBangIsNotSpecialElsewhere pins that "#!" is only special at the start
// of the file: later in the file it is ordinary unsupported input, with no
// diagnostic invented for it by W2.
func TestHashBangIsNotSpecialElsewhere(t *testing.T) {
	t.Parallel()

	l := lexer.New("test.atrn", []byte("a #! b"))
	if tok := l.Next(); tok.Kind != token.Ident {
		t.Fatalf("first token = %v, want Ident", tok.Kind)
	}
	hash := l.Next()
	if hash.Kind != token.Illegal || hash.Lexeme != "#" {
		t.Fatalf("token = %v %q, want Illegal #", hash.Kind, hash.Lexeme)
	}
	if d := l.Diagnostics(); len(d) != 0 {
		t.Fatalf("W2 invented diagnostics for unimplemented punctuation: %v", d)
	}
}

// TestMalformedUTF8Recovery covers malformed bytes, truncated sequences and
// malformed bytes between valid identifiers.
//
// Recovery is defined per byte and deterministically: a malformed sequence is
// republished as one Illegal token per undecodable byte, each with its own
// E0101, and scanning resumes at the next byte. This test therefore asserts the
// complete token stream, not just the first Illegal token, so the recovery
// policy cannot drift unnoticed.
func TestMalformedUTF8Recovery(t *testing.T) {
	t.Parallel()

	type step struct {
		kind   token.Kind
		lexeme string
	}

	tests := []struct {
		name string
		src  []byte
		want []step
	}{
		{
			name: "lone continuation byte",
			src:  []byte("a\x80b"),
			want: []step{{token.Ident, "a"}, {token.Illegal, "\x80"}, {token.Ident, "b"}},
		},
		{
			name: "invalid start byte",
			src:  []byte("a\xffb"),
			want: []step{{token.Ident, "a"}, {token.Illegal, "\xff"}, {token.Ident, "b"}},
		},
		{
			// \xe6\x9d is a truncated three-byte sequence: neither byte decodes,
			// so both are reported separately and "b" is still reached.
			name: "truncated three-byte",
			src:  []byte("a\xe6\x9db"),
			want: []step{
				{token.Ident, "a"},
				{token.Illegal, "\xe6"},
				{token.Illegal, "\x9d"},
				{token.Ident, "b"},
			},
		},
		{
			// An overlong encoding is invalid in both bytes.
			name: "overlong encoding",
			src:  []byte("a\xc0\xafb"),
			want: []step{
				{token.Ident, "a"},
				{token.Illegal, "\xc0"},
				{token.Illegal, "\xaf"},
				{token.Ident, "b"},
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l := lexer.New("test.atrn", tc.src)

			illegal := 0
			for i, w := range tc.want {
				tok := l.Next()
				if tok.Kind != w.kind || tok.Lexeme != w.lexeme {
					t.Fatalf("token %d = %v %q, want %v %q", i, tok.Kind, tok.Lexeme, w.kind, w.lexeme)
				}
				if tok.Kind == token.Illegal {
					illegal++
					// Source fidelity: the lexeme is exactly the source bytes,
					// never a substitution character.
					if got, want := tok.Lexeme, string(tc.src[tok.Pos.Offset:tok.End.Offset]); got != want {
						t.Fatalf("token %d lexeme = %q, want source bytes %q", i, got, want)
					}
					if tok.Lexeme == "\uFFFD" {
						t.Fatalf("token %d: lexer substituted U+FFFD for malformed input", i)
					}
				}
			}
			if tok := l.Next(); tok.Kind != token.EOF {
				t.Fatalf("stream continued past the expected tokens: %v %q", tok.Kind, tok.Lexeme)
			}

			diags := l.Diagnostics()
			if len(diags) != illegal {
				t.Fatalf("got %d diagnostics for %d illegal tokens: %v", len(diags), illegal, diags)
			}
			for i, d := range diags {
				if d.Code != "E0101" {
					t.Errorf("diagnostic %d code = %q, want E0101", i, d.Code)
				}
			}
		})
	}
}

// TestValidReplacementCharacterIsNotMalformed pins the size==1 distinction: a
// genuine U+FFFD in the source is valid UTF-8 text, not malformed input, so it
// must not produce E0101.
func TestValidReplacementCharacterIsNotMalformed(t *testing.T) {
	t.Parallel()

	l := lexer.New("test.atrn", []byte("a\uFFFDb"))
	if tok := l.Next(); tok.Kind != token.Ident || tok.Lexeme != "a" {
		t.Fatalf("first token = %v %q", tok.Kind, tok.Lexeme)
	}
	mid := l.Next()
	if mid.Lexeme != "\uFFFD" {
		t.Fatalf("middle token lexeme = %q, want U+FFFD", mid.Lexeme)
	}
	// "a" occupies column 1, so U+FFFD starts at column 2.
	if mid.Pos.Col != 2 {
		t.Fatalf("U+FFFD at column %d, want 2", mid.Pos.Col)
	}
	if d := l.Diagnostics(); len(d) != 0 {
		t.Fatalf("a correctly encoded U+FFFD produced diagnostics %v, want none", d)
	}
	if tok := l.Next(); tok.Kind != token.Ident || tok.Lexeme != "b" {
		t.Fatalf("token after U+FFFD = %v %q, want Ident b", tok.Kind, tok.Lexeme)
	}
}

// TestNoDiagnosticForUnimplementedFamilies is the guard against fabricating
// findings: operators, numbers, strings and punctuation are not implemented in
// W2, so they must yield Illegal placeholders with no diagnostic at all.
func TestNoDiagnosticForUnimplementedFamilies(t *testing.T) {
	t.Parallel()

	for _, src := range []string{"1", "1.5", "\"s\"", "`r`", "+", ".", "/* c */", "// c", "{", ";"} {
		l := lexer.New("test.atrn", []byte(src))
		for i := 0; i < 100; i++ {
			if tok := l.Next(); tok.Kind == token.EOF {
				break
			}
		}
		if d := l.Diagnostics(); len(d) != 0 {
			t.Errorf("%q: W2 produced diagnostics %v for an unimplemented family; want none", src, d)
		}
	}
}
