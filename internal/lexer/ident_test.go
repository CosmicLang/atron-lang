// SPDX-License-Identifier: Apache-2.0

package lexer_test

import (
	"testing"
	"unicode"
	"unicode/utf8"

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

// unicodeNumberSamples names one representative of each Unicode numeric
// category that the ratified grammar accepts through \p{N}:
//
//	Nd  U+0663 ARABIC-INDIC DIGIT THREE    decimal digit
//	Nl  U+216B ROMAN NUMERAL TWELVE        letter-like number
//	No  U+00B2 SUPERSCRIPT TWO             other number
//
// isDigit records unicode.IsDigit for the rune. It is true only for Nd, which
// is precisely the defect these tests guard against: an implementation built on
// unicode.IsDigit silently truncates identifiers at Nl and No characters.
var unicodeNumberSamples = []struct {
	r        rune
	category string
	name     string
	isDigit  bool
}{
	{'\u0663', "Nd", "U+0663 ARABIC-INDIC DIGIT THREE", true},
	{'\u216B', "Nl", "U+216B ROMAN NUMERAL TWELVE", false},
	{'\u00B2', "No", "U+00B2 SUPERSCRIPT TWO", false},
	{'\u00BE', "No", "U+00BE VULGAR FRACTION THREE QUARTERS", false},
}

// TestUnicodeNumberCategoriesAreIdentifiers is the conformance regression test
// for the ratified continuation class [\p{L}\p{N}_]*.
//
// Each case must produce exactly one Ident token followed by EOF, whose lexeme
// is the whole original source, with no lexical diagnostic, and whose exclusive
// end position is reported in bytes and runes respectively.
func TestUnicodeNumberCategoriesAreIdentifiers(t *testing.T) {
	t.Parallel()

	cases := []string{
		"a9",      // ASCII, the pre-existing behaviour
		"a\u0663", // Nd
		"a\u00B2", // No
		"a\u216B", // Nl
		"a\u00BE", // No
		"_\u00B2", // underscore followed by a number
		"valor\u00B2",
		"item\u216B",
	}

	for _, src := range cases {
		src := src
		t.Run(src, func(t *testing.T) {
			t.Parallel()

			l := lexer.New("test.atrn", []byte(src))
			tok := l.Next()

			if tok.Kind != token.Ident {
				t.Fatalf("first token kind = %v (%q), want Ident: a Unicode number must continue an identifier", tok.Kind, tok.Lexeme)
			}
			if tok.Lexeme != src {
				t.Fatalf("lexeme = %q, want the complete source %q", tok.Lexeme, src)
			}
			if d := l.Diagnostics(); len(d) != 0 {
				t.Fatalf("diagnostics = %v, want none", d)
			}

			// Exactly one token, then EOF.
			next := l.Next()
			if next.Kind != token.EOF {
				t.Fatalf("second token kind = %v (%q), want EOF: the identifier was split", next.Kind, next.Lexeme)
			}
			if next.Lexeme != "" {
				t.Fatalf("EOF lexeme = %q, want empty", next.Lexeme)
			}

			// Offsets count bytes, columns count runes (lexical.md §2, §9).
			wantBytes := len(src)
			wantRunes := utf8.RuneCountInString(src)
			if tok.Pos != (token.Pos{Offset: 0, Line: 1, Col: 1}) {
				t.Errorf("start = %+v, want offset 0 line 1 col 1", tok.Pos)
			}
			if tok.End.Offset != wantBytes {
				t.Errorf("end offset = %d, want %d byte(s)", tok.End.Offset, wantBytes)
			}
			if tok.End.Line != 1 {
				t.Errorf("end line = %d, want 1", tok.End.Line)
			}
			if tok.End.Col != 1+wantRunes {
				t.Errorf("end column = %d, want %d rune(s) after column 1", tok.End.Col, 1+wantRunes)
			}
			if next.Pos != tok.End {
				t.Errorf("EOF position = %+v, want the token end %+v", next.Pos, tok.End)
			}

			// Byte fidelity: the lexeme is a verbatim slice of the source.
			if tok.Lexeme != string([]byte(src)[tok.Pos.Offset:tok.End.Offset]) {
				t.Errorf("lexeme %q is not src[%d:%d]", tok.Lexeme, tok.Pos.Offset, tok.End.Offset)
			}
		})
	}
}

// TestUnicodeNumbersAreNotIdentifierStarts proves the first-character rule is
// unchanged: \p{N} characters continue an identifier but never begin one, so a
// number-initial sequence still falls outside the identifier family.
func TestUnicodeNumbersAreNotIdentifierStarts(t *testing.T) {
	t.Parallel()

	for _, sample := range unicodeNumberSamples {
		src := string(sample.r)

		if !unicode.IsNumber(sample.r) {
			t.Fatalf("%s: sample is not in category N", sample.name)
		}

		l := lexer.New("test.atrn", []byte(src))
		tok := l.Next()
		if tok.Kind == token.Ident {
			t.Errorf("%q (%s) lexed as Ident; a Unicode number must not start an identifier", src, sample.name)
		}
		if tok.Lexeme != src {
			t.Errorf("%q lexeme = %q, want the source rune", src, tok.Lexeme)
		}
		// W2 placeholder behaviour is unchanged: no diagnostic.
		if d := l.Diagnostics(); len(d) != 0 {
			t.Errorf("%q produced diagnostics %v, want none for an unsupported family", src, d)
		}
		if tok.End.Offset != len(src) || tok.End.Col != 2 {
			t.Errorf("%q ended at byte %d col %d, want byte %d col 2",
				src, tok.End.Offset, tok.End.Col, len(src))
		}
	}
}

// TestMarksEmojiAndSymbolsAreNotIdentifierContinuation pins the negative side of
// the continuation class: only \p{L}, \p{N} and underscore continue a name.
func TestMarksEmojiAndSymbolsAreNotIdentifierContinuation(t *testing.T) {
	t.Parallel()

	// A combining acute accent (Mn), an emoji (So), a math symbol (Sm) and
	// punctuation (Po) must each terminate the identifier.
	for _, suffix := range []string{"\u0301", "\U0001F600", "\u00D7", "\u00A1"} {
		src := "a" + suffix

		l := lexer.New("test.atrn", []byte(src))
		first := l.Next()
		if first.Kind != token.Ident || first.Lexeme != "a" {
			t.Errorf("%q: first token = %v %q, want Ident \"a\"", src, first.Kind, first.Lexeme)
		}
		if first.End.Offset != 1 || first.End.Col != 2 {
			t.Errorf("%q: identifier ended at byte %d col %d, want byte 1 col 2",
				src, first.End.Offset, first.End.Col)
		}
		if d := l.Diagnostics(); len(d) != 0 {
			t.Errorf("%q produced diagnostics %v, want none", src, d)
		}

		second := l.Next()
		if second.Kind == token.Ident {
			t.Fatalf("%q: non-identifier rune was absorbed: %v %q", src, second.Kind, second.Lexeme)
		}
		if second.Lexeme != suffix {
			t.Errorf("%q: second lexeme = %q, want %q", src, second.Lexeme, suffix)
		}
	}

	// The same rule holds after a Unicode number, not just after ASCII.
	l := lexer.New("test.atrn", []byte("a\u00B2\u0301"))
	first := l.Next()
	if first.Kind != token.Ident || first.Lexeme != "a\u00B2" {
		t.Fatalf("first token = %v %q, want Ident %q", first.Kind, first.Lexeme, "a\u00B2")
	}
	if second := l.Next(); second.Kind == token.Ident {
		t.Fatalf("combining mark was absorbed after a Unicode number: %q", second.Lexeme)
	}
}

// TestKeywordClassificationUnchangedByUnicodeContinuation proves the correction
// affects only the continuation class: exact keyword matching still wins, and a
// keyword extended by a Unicode number is an ordinary identifier.
func TestKeywordClassificationUnchangedByUnicodeContinuation(t *testing.T) {
	t.Parallel()

	keywords := map[string]token.Kind{
		"const":  token.Const,
		"fn":     token.Fn,
		"return": token.Return,
		"true":   token.True,
	}
	for src, want := range keywords {
		l := lexer.New("test.atrn", []byte(src))
		if got := l.Next().Kind; got != want {
			t.Errorf("%q lexed as %v, want %v", src, got, want)
		}
	}

	// A keyword prefix extended by a Unicode number is not a keyword.
	for _, src := range []string{"const\u00B2", "fn\u216B", "return\u0663", "true\u00BE"} {
		l := lexer.New("test.atrn", []byte(src))
		tok := l.Next()
		if tok.Kind != token.Ident {
			t.Errorf("%q lexed as %v, want Ident", src, tok.Kind)
		}
		if tok.Lexeme != src {
			t.Errorf("%q lexeme = %q, want %q", src, tok.Lexeme, src)
		}
	}

	// Two identifiers separated by a space still tokenise separately.
	toks := collect(t, "a\u00B2 b\u216B")
	if len(toks) != 3 {
		t.Fatalf("got %d tokens, want 3 (two identifiers and EOF)", len(toks))
	}
	if toks[0].Lexeme != "a\u00B2" || toks[1].Lexeme != "b\u216B" {
		t.Fatalf("lexemes = %q, %q, want %q, %q", toks[0].Lexeme, toks[1].Lexeme, "a\u00B2", "b\u216B")
	}
	if toks[1].Pos != (token.Pos{Offset: 4, Line: 1, Col: 4}) {
		t.Errorf("second identifier position = %+v, want offset 4 col 4", toks[1].Pos)
	}
}

// TestUnicodeConformanceIsNotSatisfiableByIsDigit keeps the suite honest: if
// every accepted number character were also an unicode.IsDigit character, the
// tests above would pass against the defective implementation and prove nothing.
func TestUnicodeConformanceIsNotSatisfiableByIsDigit(t *testing.T) {
	t.Parallel()

	nonDigit := 0
	for _, sample := range unicodeNumberSamples {
		if got := unicode.IsDigit(sample.r); got != sample.isDigit {
			t.Errorf("%s: unicode.IsDigit = %v, recorded %v", sample.name, got, sample.isDigit)
		}
		if !unicode.IsNumber(sample.r) {
			t.Errorf("%s: unicode.IsNumber = false, want true", sample.name)
		}
		if !sample.isDigit {
			nonDigit++
		}
	}
	if nonDigit == 0 {
		t.Fatal("no accepted sample lies outside unicode.IsDigit; the suite cannot detect the defect")
	}
}
