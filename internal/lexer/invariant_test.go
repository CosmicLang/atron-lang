// SPDX-License-Identifier: Apache-2.0

package lexer_test

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/CosmicLang/atron-lang/internal/lexer"
	"github.com/CosmicLang/atron-lang/internal/token"
)

// TestLexemeMatchesSourceBytes enforces the §9 invariant on a wide range of
// inputs: for every real token, Lexeme is exactly the bytes it spans. The only
// permitted exception is a synthetic Newline, which W2 does not yet emit.
func TestLexemeMatchesSourceBytes(t *testing.T) {
	t.Parallel()

	sources := []string{
		"",
		"a",
		"main",
		"_1",
		"naïve Ω 東京",
		"a\nb\r\nc\rd",
		"\xef\xbb\xbfmain",
		"#!atron\nmain",
		"a\x80b\xffc",
		"a\uFFFD\uFEFFb",
		"a b\tc\fd\ve",
		strings.Repeat("identifier ", 50),
		strings.Repeat("Ω\n", 40),
	}

	for _, src := range sources {
		l := lexer.New("test.atrn", []byte(src))
		for i := 0; i < 10000; i++ {
			tok := l.Next()
			if tok.Kind == token.EOF {
				break
			}
			if tok.Kind == token.Newline {
				if tok.Lexeme != "" {
					t.Fatalf("%q: synthetic Newline has lexeme %q", src, tok.Lexeme)
				}
				continue
			}
			start, end := tok.Pos.Offset, tok.End.Offset
			if start < 0 || end > len(src) || start > end {
				t.Fatalf("%q: token span [%d,%d) out of range", src, start, end)
			}
			if got, want := tok.Lexeme, src[start:end]; got != want {
				t.Fatalf("%q: lexeme %q != source slice %q", src, got, want)
			}
		}
	}
}

// TestDeterministicProgression checks that offsets never move backwards and that
// every emitted token advances the cursor (or is EOF at the end).
func TestDeterministicProgression(t *testing.T) {
	t.Parallel()

	src := "a\x80Ω\r\nb\rc d\uFEFFe"
	l := lexer.New("test.atrn", []byte(src))
	prev := 0
	for i := 0; i < 1000; i++ {
		tok := l.Next()
		if tok.Kind == token.EOF {
			break
		}
		if tok.Pos.Offset < prev {
			t.Fatalf("cursor moved backwards from %d to %d", prev, tok.Pos.Offset)
		}
		if tok.End.Offset < tok.Pos.Offset {
			t.Fatalf("token End %d precedes Pos %d", tok.End.Offset, tok.Pos.Offset)
		}
		prev = tok.End.Offset
	}
}

// TestRandomInputAlwaysTerminates is a bounded property test: arbitrary bytes
// must always produce a finite token stream that ends in EOF, without panics or
// infinite loops. It is a smoke test for the traversal machinery, not the M1
// FuzzLexer acceptance gate, which requires the complete scanner.
func TestRandomInputAlwaysTerminates(t *testing.T) {
	t.Parallel()

	rng := rand.New(rand.NewSource(20261018))
	for iter := 0; iter < 200; iter++ {
		n := rng.Intn(64)
		buf := make([]byte, n)
		for i := range buf {
			buf[i] = byte(rng.Intn(256))
		}

		l := lexer.New("random.atrn", buf)
		tokens := 0
		for {
			tok := l.Next()
			tokens++
			if tok.Kind == token.EOF {
				break
			}
			if tokens > 10*len(buf)+64 {
				t.Fatalf("iteration %d (%d bytes) produced over %d tokens: no progress", iter, n, tokens)
			}
			if tok.Pos.Offset < 0 || tok.End.Offset > len(buf) || tok.Pos.Offset > tok.End.Offset {
				t.Fatalf("iteration %d: token span out of range: %+v..%+v", iter, tok.Pos, tok.End)
			}
		}
		// Every diagnostic must carry a position that slices the source.
		for _, d := range l.Diagnostics() {
			if d.Pos.Offset < 0 || d.End.Offset > len(buf) || d.Pos.Offset > d.End.Offset {
				t.Fatalf("iteration %d: diagnostic span out of range: %+v..%+v", iter, d.Pos, d.End)
			}
			if d.File != "random.atrn" {
				t.Fatalf("iteration %d: diagnostic file = %q", iter, d.File)
			}
		}
	}
}

// TestDiagnosticsAreCopied checks that callers cannot mutate the lexer's
// accumulated diagnostics.
func TestDiagnosticsAreCopied(t *testing.T) {
	t.Parallel()

	l := lexer.New("test.atrn", []byte("\x80"))
	if tok := l.Next(); tok.Kind != token.Illegal {
		t.Fatalf("token = %v, want Illegal", tok.Kind)
	}
	first := l.Diagnostics()
	if len(first) != 1 {
		t.Fatalf("got %d diagnostics, want 1", len(first))
	}
	first[0].Code = "E9999"

	second := l.Diagnostics()
	if len(second) != 1 || second[0].Code != "E0101" {
		t.Fatalf("mutating the returned slice changed the lexer state: %v", second)
	}
}

// TestFileIsReportedVerbatim checks that New stores the file name as given and
// that diagnostics carry it unchanged, without touching the filesystem.
func TestFileIsReportedVerbatim(t *testing.T) {
	t.Parallel()

	const name = "does/not/exist.atrn"
	l := lexer.New(name, []byte("\x80"))
	if got := l.File(); got != name {
		t.Fatalf("File() = %q, want %q", got, name)
	}
	if tok := l.Next(); tok.Kind != token.Illegal {
		t.Fatalf("token = %v, want Illegal", tok.Kind)
	}
	diags := l.Diagnostics()
	if len(diags) != 1 {
		t.Fatalf("got %d diagnostics, want 1", len(diags))
	}
	if diags[0].File != name {
		t.Fatalf("diagnostic file = %q, want %q", diags[0].File, name)
	}
}

// TestNoDiagnosticsForValidIdentifiers pins that valid input never produces a
// diagnostic, in particular the `_1` regression the checklist requires.
func TestNoDiagnosticsForValidIdentifiers(t *testing.T) {
	t.Parallel()

	for _, src := range []string{"_1", "_name", "naïve", "Ω", "東京", "main", "const x"} {
		l := lexer.New("test.atrn", []byte(src))
		for {
			if tok := l.Next(); tok.Kind == token.EOF {
				break
			}
		}
		if d := l.Diagnostics(); len(d) != 0 {
			t.Errorf("%q produced diagnostics %v, want none", src, d)
		}
	}
}

// TestSourceIsNotModified pins the immutable-input contract.
func TestSourceIsNotModified(t *testing.T) {
	t.Parallel()

	src := []byte("a\r\nb\x80 Ω")
	original := string(src)
	l := lexer.New("test.atrn", src)
	for {
		if tok := l.Next(); tok.Kind == token.EOF {
			break
		}
	}
	if got := string(src); got != original {
		t.Fatalf("source buffer was modified: %q -> %q", original, got)
	}
}
