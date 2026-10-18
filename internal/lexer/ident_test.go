// SPDX-License-Identifier: Apache-2.0

package lexer_test

import (
	"testing"

	"github.com/CosmicLang/atron-lang/internal/lexer"
	"github.com/CosmicLang/atron-lang/internal/token"
)

// collect consumes the whole stream and returns every token up to and including
// the first EOF.
func collect(t *testing.T, src string) []token.Token {
	t.Helper()
	l := lexer.New("test.atrn", []byte(src))
	var toks []token.Token
	for i := 0; i < 10000; i++ {
		tok := l.Next()
		toks = append(toks, tok)
		if tok.Kind == token.EOF {
			return toks
		}
	}
	t.Fatalf("Next did not reach EOF within 10000 tokens for %q", src)
	return nil
}

// TestASCIIIdentifiers covers plain ASCII identifiers and the keyword boundary.
func TestASCIIIdentifiers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		src  string
		want token.Kind
	}{
		{"x", token.Ident},
		{"name", token.Ident},
		{"Name", token.Ident},
		{"NAME", token.Ident},
		{"snake_case_name", token.Ident},
		{"x1", token.Ident},
		{"abc123", token.Ident},
		{"main", token.Ident}, // ADR-0005
		{"let", token.Ident},  // ADR-0004: no longer reserved
		{"const", token.Const},
		{"var", token.Var},
		{"fn", token.Fn},
		{"return", token.Return},
		{"if", token.If},
		{"else", token.Else},
		{"while", token.While},
		{"for", token.For},
		{"in", token.In},
		{"break", token.Break},
		{"continue", token.Continue},
		{"struct", token.Struct},
		{"impl", token.Impl},
		{"import", token.Import},
		{"pub", token.Pub},
		{"as", token.As},
		{"and", token.And},
		{"or", token.Or},
		{"not", token.Not},
		{"true", token.True},
		{"false", token.False},
		{"nil", token.Nil},
	}

	for _, tc := range tests {
		toks := collect(t, tc.src)
		if len(toks) != 2 {
			t.Errorf("%q produced %d tokens, want 2 (token + EOF)", tc.src, len(toks))
			continue
		}
		if toks[0].Kind != tc.want {
			t.Errorf("%q lexed as %v, want %v", tc.src, toks[0].Kind, tc.want)
		}
		if toks[0].Lexeme != tc.src {
			t.Errorf("%q lexeme = %q, want the source text", tc.src, toks[0].Lexeme)
		}
	}
}

// TestKeywordRecognitionIsCaseSensitive proves that only lower-case ASCII is
// recognised and that a keyword prefix does not split the identifier.
func TestKeywordRecognitionIsCaseSensitive(t *testing.T) {
	t.Parallel()

	for _, src := range []string{"IF", "If", "CONST", "True", "Nil", "consts", "iff", "fnord"} {
		toks := collect(t, src)
		if toks[0].Kind != token.Ident {
			t.Errorf("%q lexed as %v, want Ident", src, toks[0].Kind)
		}
		if toks[0].Lexeme != src {
			t.Errorf("%q lexeme = %q, want one whole identifier", src, toks[0].Lexeme)
		}
	}
}

// TestUnderscoreIdentifiers pins the checklist requirement that `_`, `_1` and
// `_name` are identifiers and that `_1` never yields E0103 or any diagnostic.
func TestUnderscoreIdentifiers(t *testing.T) {
	t.Parallel()

	for _, src := range []string{"_", "_1", "_name", "__", "_x9"} {
		l := lexer.New("test.atrn", []byte(src))
		tok := l.Next()
		if tok.Kind != token.Ident {
			t.Errorf("%q lexed as %v, want Ident", src, tok.Kind)
		}
		if tok.Lexeme != src {
			t.Errorf("%q lexeme = %q, want %q", src, tok.Lexeme, src)
		}
		if d := l.Diagnostics(); len(d) != 0 {
			t.Errorf("%q produced %d diagnostics (%v), want none", src, len(d), d)
		}
	}
}

// TestIdentifiersEndingInDigits covers the identifier continuation rule for
// ASCII digits.
func TestIdentifiersEndingInDigits(t *testing.T) {
	t.Parallel()

	for _, src := range []string{"a1", "_1", "x9", "value2026", "_0_1"} {
		toks := collect(t, src)
		if toks[0].Kind != token.Ident || toks[0].Lexeme != src {
			t.Errorf("%q lexed as %v %q, want Ident %q", src, toks[0].Kind, toks[0].Lexeme, src)
		}
	}
}

// TestUnicodeIdentifiers covers the ratified [\p{L}_][\p{L}\p{N}_]* grammar,
// including non-Latin scripts and combining-mark rejection.
func TestUnicodeIdentifiers(t *testing.T) {
	t.Parallel()

	valid := []string{"naïve", "Ω", "東京", "Ωmega", "café1", "Ω2"}
	for _, src := range valid {
		l := lexer.New("test.atrn", []byte(src))
		tok := l.Next()
		if tok.Kind != token.Ident {
			t.Errorf("%q lexed as %v, want Ident", src, tok.Kind)
		}
		if tok.Lexeme != src {
			t.Errorf("%q lexeme = %q, want %q", src, tok.Lexeme, src)
		}
		if d := l.Diagnostics(); len(d) != 0 {
			t.Errorf("%q produced diagnostics %v, want none", src, d)
		}
	}

	// A combining mark alone is not an identifier character. It ends the
	// preceding identifier and is then an invalid character in its own right;
	// W2 does not yet report E0101 for it because that verdict belongs with the
	// full character classification of a later work package. It must at least
	// never be absorbed into the identifier.
	l := lexer.New("test.atrn", []byte("a\u0301"))
	first := l.Next()
	if first.Kind != token.Ident || first.Lexeme != "a" {
		t.Fatalf("first token = %v %q, want Ident \"a\"", first.Kind, first.Lexeme)
	}
	if second := l.Next(); second.Lexeme == "" || second.Lexeme[0] == 'a' {
		t.Fatalf("combining mark was absorbed into the identifier: %q", second.Lexeme)
	}
}
