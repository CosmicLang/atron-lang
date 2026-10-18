// SPDX-License-Identifier: Apache-2.0

package token_test

import (
	"testing"

	"github.com/CosmicLang/atron-lang/internal/token"
)

// TestEndsStatementForEveryKind evaluates the predicate for all 62 kinds and
// compares it against the ratified set spelled out here independently of the
// implementation, so the two cannot drift together.
func TestEndsStatementForEveryKind(t *testing.T) {
	t.Parallel()

	wantTrue := map[token.Kind]bool{
		token.Ident:     true,
		token.Int:       true,
		token.Float:     true,
		token.String:    true,
		token.RawString: true,
		token.True:      true,
		token.False:     true,
		token.Nil:       true,
		token.Return:    true,
		token.Break:     true,
		token.Continue:  true,
		token.RParen:    true,
		token.RBracket:  true,
		token.RBrace:    true,
		token.Question:  true,
	}

	trueCount := 0
	for i := 0; i < token.KindCount; i++ {
		k := token.Kind(i)
		want := wantTrue[k]
		if got := k.EndsStatement(); got != want {
			t.Errorf("%v.EndsStatement() = %t, want %t", k, got, want)
		}
		if want {
			trueCount++
		}
	}
	if trueCount != 15 {
		t.Fatalf("ratified eligible set has %d kinds, want 15", trueCount)
	}

	// Spot-check the exclusions the contract calls out explicitly.
	excluded := []token.Kind{
		token.Illegal, token.EOF, token.Newline,
		token.Plus, token.Assign, token.Comma, token.Colon, token.Dot,
		token.LParen, token.LBracket, token.LBrace, token.FatArrow,
		token.DotDot, token.DotDotEq, token.StarStar,
	}
	for _, k := range excluded {
		if k.EndsStatement() {
			t.Errorf("%v.EndsStatement() = true, want false", k)
		}
	}
}

// TestTokenSpanIsExclusive verifies the span helper agrees with the exclusive
// End semantics, including the zero-width synthetic Newline.
func TestTokenSpanIsExclusive(t *testing.T) {
	t.Parallel()

	src := []byte("hello")
	tok := token.Token{
		Kind:   token.Ident,
		Lexeme: "hello",
		Pos:    token.Pos{Offset: 0, Line: 1, Col: 1},
		End:    token.Pos{Offset: 5, Line: 1, Col: 6},
	}
	start, end := tok.Span()
	if got := string(src[start:end]); got != tok.Lexeme {
		t.Fatalf("src[Pos:End] = %q, want %q", got, tok.Lexeme)
	}

	synthetic := token.Token{
		Kind: token.Newline,
		Pos:  token.Pos{Offset: 5, Line: 1, Col: 6},
		End:  token.Pos{Offset: 5, Line: 1, Col: 6},
	}
	start, end = synthetic.Span()
	if start != 5 || end != 5 || synthetic.Lexeme != "" {
		t.Fatalf("synthetic Newline span = [%d,%d) lexeme %q, want empty span at 5",
			start, end, synthetic.Lexeme)
	}
}

// TestPosString covers the "line:column" rendering and confirms Offset is not
// part of it: offsets address bytes, not readers.
func TestPosString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		pos  token.Pos
		want string
	}{
		{token.Pos{Offset: 0, Line: 1, Col: 1}, "1:1"},
		{token.Pos{Offset: 40, Line: 3, Col: 7}, "3:7"},
		// A token after three multi-byte runes starts at column 4 but byte 9,
		// which is why Offset is excluded from the human-readable form.
		{token.Pos{Offset: 9, Line: 2, Col: 4}, "2:4"},
	}
	for _, tc := range tests {
		if got := tc.pos.String(); got != tc.want {
			t.Errorf("Pos%+v.String() = %q, want %q", tc.pos, got, tc.want)
		}
	}
}
