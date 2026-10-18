// SPDX-License-Identifier: Apache-2.0

package token_test

import (
	"strings"
	"testing"

	"github.com/CosmicLang/atron-lang/internal/token"
)

// TestKeywordInventoryIsExactly22 enumerates the keyword range and checks each
// spelling, so the keyword table and the enumeration agree in both directions.
func TestKeywordInventoryIsExactly22(t *testing.T) {
	t.Parallel()

	want := map[string]token.Kind{
		"const":    token.Const,
		"var":      token.Var,
		"fn":       token.Fn,
		"return":   token.Return,
		"if":       token.If,
		"else":     token.Else,
		"while":    token.While,
		"for":      token.For,
		"in":       token.In,
		"break":    token.Break,
		"continue": token.Continue,
		"struct":   token.Struct,
		"impl":     token.Impl,
		"import":   token.Import,
		"pub":      token.Pub,
		"as":       token.As,
		"and":      token.And,
		"or":       token.Or,
		"not":      token.Not,
		"true":     token.True,
		"false":    token.False,
		"nil":      token.Nil,
	}

	if token.KeywordCount != 22 {
		t.Fatalf("KeywordCount = %d, want 22", token.KeywordCount)
	}
	got := token.Keywords()
	if len(got) != token.KeywordCount {
		t.Fatalf("Keywords() has %d entries, want %d", len(got), token.KeywordCount)
	}
	for spelling, kind := range want {
		if got[spelling] != kind {
			t.Errorf("Keywords()[%q] = %v, want %v", spelling, got[spelling], kind)
		}
	}
	for spelling, kind := range got {
		if want[spelling] != kind {
			t.Errorf("unexpected keyword %q -> %v", spelling, kind)
		}
	}

	// The enumeration's keyword range must be exactly the 22 kinds in the table.
	if n := int(token.Nil) - int(token.Const) + 1; n != 22 {
		t.Fatalf("keyword range spans %d kinds, want 22", n)
	}
	for i := token.Const; i <= token.Nil; i++ {
		if !i.IsKeyword() {
			t.Errorf("%v is inside the keyword range but IsKeyword() is false", i)
		}
	}
	if token.Ident.IsKeyword() || token.FatArrow.IsKeyword() || token.Newline.IsKeyword() {
		t.Error("a non-keyword kind reports IsKeyword() == true")
	}
}

// TestLookupIdent checks exact, case-sensitive keyword recognition and the
// classification of the identifiers the contract singles out.
func TestLookupIdent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		text string
		want token.Kind
	}{
		{"const", token.Const},
		{"var", token.Var},
		{"fn", token.Fn},
		{"nil", token.Nil},
		{"true", token.True},
		{"false", token.False},
		// ADR-0004: `let` is no longer reserved and there is no alias.
		{"let", token.Ident},
		{"Let", token.Ident},
		// ADR-0005: `main` is an ordinary identifier.
		{"main", token.Ident},
		{"Main", token.Ident},
		// Identifier forms that must never be mistaken for keywords.
		{"_", token.Ident},
		{"_1", token.Ident},
		{"_name", token.Ident},
		{"x1", token.Ident},
		{"naïve", token.Ident},
		{"Ω", token.Ident},
		{"東京", token.Ident},
		// Case sensitivity: keywords are lower-case ASCII only.
		{"IF", token.Ident},
		{"If", token.Ident},
		{"CONST", token.Ident},
		{"True", token.Ident},
		// Keyword plus trailing text is one identifier.
		{"consts", token.Ident},
		{"iff", token.Ident},
		{"", token.Ident},
	}

	for _, tc := range tests {
		if got := token.LookupIdent(tc.text); got != tc.want {
			t.Errorf("LookupIdent(%q) = %v, want %v", tc.text, got, tc.want)
		}
	}
}

// TestKeywordsReturnsCopy proves that callers cannot mutate the contract.
func TestKeywordsReturnsCopy(t *testing.T) {
	t.Parallel()

	m := token.Keywords()
	m["let"] = token.Const
	delete(m, "const")

	if token.LookupIdent("let") != token.Ident {
		t.Fatal("mutating the Keywords() result changed keyword recognition")
	}
	if token.LookupIdent("const") != token.Const {
		t.Fatal("deleting from the Keywords() result removed a keyword")
	}
}

// TestCommentMetadataPreservesBytes checks that the metadata type can carry raw
// source text verbatim, including CRLF, which is what the formatter needs.
func TestCommentMetadataPreservesBytes(t *testing.T) {
	t.Parallel()

	c := token.Comment{
		Text:  "/* a\r\nb */",
		Pos:   token.Pos{Offset: 0, Line: 1, Col: 1},
		End:   token.Pos{Offset: 10, Line: 2, Col: 5},
		Block: true,
	}
	if !strings.Contains(c.Text, "\r\n") {
		t.Fatal("comment text lost its CRLF bytes")
	}
	if !c.Block {
		t.Fatal("block comment not marked as block")
	}
	if got, want := c.End.Offset-c.Pos.Offset, len(c.Text); got != want {
		t.Fatalf("comment span %d bytes, lexeme %d bytes", got, want)
	}
}
