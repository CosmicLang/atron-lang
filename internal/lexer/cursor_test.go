// SPDX-License-Identifier: Apache-2.0

package lexer

import (
	"math/rand"
	"testing"

	"github.com/CosmicLang/atron-lang/internal/token"
)

// TestPeekAtEOF covers the end-of-input branch of peek.
func TestPeekAtEOF(t *testing.T) {
	t.Parallel()

	for _, src := range []string{"", "a"} {
		l := New("test.atrn", []byte(src))
		for !l.atEOF() {
			if _, _, ok := l.peek(); !ok {
				t.Fatal("peek reported EOF before atEOF did")
			}
			l.advance()
		}
		if r, size, ok := l.peek(); ok || r != 0 || size != 0 {
			t.Fatalf("%q: peek at EOF = (%q,%d,%t), want (0,0,false)", src, r, size, ok)
		}
	}
}

// TestAdvanceAtEOF covers the end-of-input branch of advance, which must report
// failure rather than move or wrap.
func TestAdvanceAtEOF(t *testing.T) {
	t.Parallel()

	l := New("test.atrn", []byte("a"))
	l.advance()
	before := l.pos()
	if r, size, ok := l.advance(); ok || r != 0 || size != 0 {
		t.Fatalf("advance at EOF = (%q,%d,%t), want (0,0,false)", r, size, ok)
	}
	if after := l.pos(); after != before {
		t.Fatalf("advance at EOF moved the cursor from %+v to %+v", before, after)
	}
}

// TestAdvanceConsumesCRLFFirst covers the CRLF branch directly, including the
// byte width it reports.
func TestAdvanceConsumesCRLFFirst(t *testing.T) {
	t.Parallel()

	l := New("test.atrn", []byte("\r\nx"))
	r, size, ok := l.advance()
	if !ok || r != '\r' || size != 1 {
		t.Fatalf("advance = (%q,%d,%t), want (CR,1,true)", r, size, ok)
	}
	if l.off != 2 || l.line != 2 || l.col != 1 {
		t.Fatalf("cursor after CRLF = off %d line %d col %d, want 2/2/1", l.off, l.line, l.col)
	}
}

// TestAdvanceStandaloneCR covers a CR that is not followed by LF, including a CR
// as the very last byte.
func TestAdvanceStandaloneCR(t *testing.T) {
	t.Parallel()

	for _, src := range []string{"\rx", "\r"} {
		l := New("test.atrn", []byte(src))
		l.advance()
		if l.off != 1 || l.line != 2 || l.col != 1 {
			t.Errorf("%q: cursor = off %d line %d col %d, want 1/2/1", src, l.off, l.line, l.col)
		}
	}
}

// TestDigitPredicates covers the ASCII digit classification used by the numeric
// work package, including its false cases.
func TestDigitPredicates(t *testing.T) {
	t.Parallel()

	for r := '0'; r <= '9'; r++ {
		if !isDigit(r) {
			t.Errorf("isDigit(%q) = false, want true", r)
		}
	}
	for _, r := range []rune{'/', ':', 'a', '_', 'Ω', '²', 0} {
		if isDigit(r) {
			t.Errorf("isDigit(%q) = true, want false", r)
		}
	}

	l := New("test.atrn", []byte("7"))
	if !l.digitStart() {
		t.Error("digitStart() = false on a digit")
	}
	l = New("test.atrn", []byte("x"))
	if l.digitStart() {
		t.Error("digitStart() = true on a letter")
	}
	l = New("test.atrn", nil)
	if l.digitStart() {
		t.Error("digitStart() = true at EOF")
	}
}

// TestIdentPredicates covers the identifier character classes, including the
// categories the contract excludes.
func TestIdentPredicates(t *testing.T) {
	t.Parallel()

	starts := []rune{'_', 'a', 'Z', 'Ω', '東', 'ñ'}
	for _, r := range starts {
		if !isIdentStart(r) {
			t.Errorf("isIdentStart(%q) = false, want true", r)
		}
	}
	notStarts := []rune{'0', '9', '-', '.', ' ', '\u0301', '\uFFFD', 0}
	for _, r := range notStarts {
		if isIdentStart(r) {
			t.Errorf("isIdentStart(%q) = true, want false", r)
		}
	}

	continues := []rune{'_', 'a', '0', 'Ω', '東', '\u0660'}
	for _, r := range continues {
		if !isIdentContinue(r) {
			t.Errorf("isIdentContinue(%q) = false, want true", r)
		}
	}
	notContinues := []rune{'-', '.', ' ', '\u0301', '!', 0}
	for _, r := range notContinues {
		if isIdentContinue(r) {
			t.Errorf("isIdentContinue(%q) = true, want false", r)
		}
	}
}

// TestSkipSpaceStopsAtNonSpace covers the early return of skipSpace.
func TestSkipSpaceStopsAtNonSpace(t *testing.T) {
	t.Parallel()

	l := New("test.atrn", []byte("  \t\n\r\n\rx"))
	l.skipSpace()
	if r, _, _ := l.peek(); r != 'x' {
		t.Fatalf("skipSpace stopped at %q, want 'x'", r)
	}
	if l.line != 4 || l.col != 1 {
		t.Fatalf("cursor = line %d col %d, want 4/1", l.line, l.col)
	}

	l = New("test.atrn", []byte("   "))
	l.skipSpace()
	if !l.atEOF() {
		t.Fatal("skipSpace did not consume trailing whitespace to EOF")
	}
}

// TestSkipBOMOnlyAtStart covers the size guard and the non-BOM path.
func TestSkipBOMOnlyAtStart(t *testing.T) {
	t.Parallel()

	l := New("test.atrn", []byte("\xef\xbb"))
	l.skipBOM()
	if l.off != 0 {
		t.Fatalf("a truncated BOM was skipped: offset %d", l.off)
	}

	l = New("test.atrn", []byte("ab"))
	l.skipBOM()
	if l.off != 0 {
		t.Fatalf("a non-BOM prefix was skipped: offset %d", l.off)
	}
}

// TestSkipShebangPrefixes covers both early returns and the run-to-EOF case.
func TestSkipShebangPrefixes(t *testing.T) {
	t.Parallel()

	for _, src := range []string{"", "#", "a"} {
		l := New("test.atrn", []byte(src))
		l.skipShebang()
		if l.off != 0 {
			t.Errorf("%q: cursor moved to %d, want 0", src, l.off)
		}
	}

	l := New("test.atrn", []byte("#!/bin/atron"))
	l.skipShebang()
	if !l.atEOF() {
		t.Fatal("a shebang without a line break was not skipped to EOF")
	}
}

// TestCursorAlwaysProgresses is a bounded property test over arbitrary bytes:
// repeated advance must terminate and every offset must stay inside the buffer.
func TestCursorAlwaysProgresses(t *testing.T) {
	t.Parallel()

	rng := rand.New(rand.NewSource(7))
	for iter := 0; iter < 100; iter++ {
		buf := make([]byte, rng.Intn(48))
		for i := range buf {
			buf[i] = byte(rng.Intn(256))
		}
		l := New("random.atrn", buf)
		steps := 0
		last := l.off
		for !l.atEOF() {
			if _, _, ok := l.advance(); !ok {
				t.Fatalf("iteration %d: advance failed before EOF", iter)
			}
			if l.off <= last {
				t.Fatalf("iteration %d: no progress at offset %d", iter, last)
			}
			if l.off > len(buf) {
				t.Fatalf("iteration %d: offset %d past buffer length %d", iter, l.off, len(buf))
			}
			last = l.off

			steps++
			if steps > len(buf) {
				t.Fatalf("iteration %d: more steps (%d) than bytes (%d)", iter, steps, len(buf))
			}
		}
	}
}

// TestPosMatchesCursorFields pins that pos() is derived from the cursor fields
// rather than from any cached copy.
func TestPosMatchesCursorFields(t *testing.T) {
	t.Parallel()

	l := New("test.atrn", []byte("ab"))
	l.advance()
	want := token.Pos{Offset: 1, Line: 1, Col: 2}
	if got := l.pos(); got != want {
		t.Fatalf("pos() = %+v, want %+v", got, want)
	}
}

// TestInvalidUTF8AtEOF covers the end-of-input guard of invalidUTF8: exhausted
// input is not malformed, it is simply finished.
func TestInvalidUTF8AtEOF(t *testing.T) {
	t.Parallel()

	l := New("test.atrn", nil)
	if l.invalidUTF8() {
		t.Fatal("invalidUTF8() = true at EOF")
	}

	l = New("test.atrn", []byte("ok"))
	l.advance()
	l.advance()
	if !l.atEOF() {
		t.Fatal("cursor not at EOF")
	}
	if l.invalidUTF8() {
		t.Fatal("invalidUTF8() = true after consuming valid input")
	}
}

// TestScanInvalidRuneAtEOFYieldsEOF covers the defensive end-of-input branch of
// scanInvalidRune. Next never calls it at EOF, so the test drives it directly to
// prove the branch returns a well-formed EOF token instead of reading past the
// buffer.
func TestScanInvalidRuneAtEOFYieldsEOF(t *testing.T) {
	t.Parallel()

	l := New("test.atrn", nil)
	tok := l.scanInvalidRune()
	if tok.Kind != token.EOF {
		t.Fatalf("scanInvalidRune at EOF = %v, want EOF", tok.Kind)
	}
	if tok.Pos != tok.End || tok.Pos != (token.Pos{Offset: 0, Line: 1, Col: 1}) {
		t.Fatalf("scanInvalidRune at EOF = %+v..%+v, want an empty span at 1:1", tok.Pos, tok.End)
	}
	if tok.Lexeme != "" {
		t.Fatalf("EOF token lexeme = %q, want empty", tok.Lexeme)
	}
}
