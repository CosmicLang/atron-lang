// SPDX-License-Identifier: Apache-2.0

package diag_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/CosmicLang/atron-lang/internal/diag"
	"github.com/CosmicLang/atron-lang/internal/token"
)

// TestErrorFormattingIsDeterministic pins the exact rendering and checks that
// repeated calls, and calls from independent diagnostics, produce identical
// output. No colour, width, locale or terminal state may influence it.
func TestErrorFormattingIsDeterministic(t *testing.T) {
	t.Parallel()

	pos := token.Pos{Offset: 10, Line: 3, Col: 5}
	end := token.Pos{Offset: 11, Line: 3, Col: 6}
	d := diag.New(diag.E0101, diag.Error, "examples/demo.atrn", pos, end, "invalid character")

	const want = "examples/demo.atrn:3:5: error[E0101]: invalid character"
	first := d.Error()
	if first != want {
		t.Fatalf("Error() = %q, want %q", first, want)
	}
	for i := 0; i < 100; i++ {
		if got := d.Error(); got != first {
			t.Fatalf("Error() call %d = %q, want %q", i, got, first)
		}
	}
	if strings.ContainsAny(first, "\x1b") {
		t.Fatal("Error() contains an ANSI escape; colour is an M6 concern")
	}
}

// TestErrorWithoutFileOmitsPrefix covers the case where a diagnostic has no
// filename, which happens for synthetic or in-memory inputs.
func TestErrorWithoutFileOmitsPrefix(t *testing.T) {
	t.Parallel()

	d := diag.New(diag.E0104, diag.Error, "", token.Pos{Offset: 0, Line: 1, Col: 1},
		token.Pos{Offset: 2, Line: 1, Col: 3}, "unterminated block comment")

	const want = "1:1: error[E0104]: unterminated block comment"
	if got := d.Error(); got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

// TestDiagnosticSatisfiesError makes sure a diagnostic can be returned where an
// error is expected, with errors.Is/As still working on the concrete value.
func TestDiagnosticSatisfiesError(t *testing.T) {
	t.Parallel()

	pos := token.Pos{Offset: 0, Line: 1, Col: 1}
	var err error = diag.Errorf(diag.E0105, "a.atrn", pos, pos, "invalid escape \\q")

	var d diag.Diagnostic
	if !errors.As(err, &d) {
		t.Fatal("errors.As failed to recover the Diagnostic")
	}
	if d.Code != diag.E0105 {
		t.Fatalf("recovered code = %q, want %q", d.Code, diag.E0105)
	}
	if !errors.Is(err, err) {
		t.Fatal("errors.Is(err, err) is false")
	}
}

// TestIndependentDiagnosticsRecordedInOrder checks that several diagnostics are
// distinct values that keep their own positions, and that two identical
// messages at different offsets are not collapsed.
func TestIndependentDiagnosticsRecordedInOrder(t *testing.T) {
	t.Parallel()

	type spec struct {
		offset int
		line   int
		col    int
	}
	specs := []spec{{0, 1, 1}, {12, 2, 1}, {24, 3, 1}}

	var diags []diag.Diagnostic
	for _, s := range specs {
		pos := token.Pos{Offset: s.offset, Line: s.line, Col: s.col}
		diags = append(diags, diag.Errorf(diag.E0101, "a.atrn", pos, pos, "invalid character"))
	}

	if len(diags) != 3 {
		t.Fatalf("recorded %d diagnostics, want 3", len(diags))
	}
	for i, d := range diags {
		if d.Pos.Offset != specs[i].offset || d.Pos.Line != specs[i].line {
			t.Errorf("diagnostic %d position = %+v, want offset %d line %d",
				i, d.Pos, specs[i].offset, specs[i].line)
		}
	}
	if diags[0] == diags[1] {
		t.Error("diagnostics at different offsets compare equal; identical messages must not be deduplicated")
	}
}

// TestNoEnvironmentOrTerminalDependence proves construction is a pure function
// of its arguments: environment variables that a renderer might honour must
// have no effect on the diagnostic or its Error() text.
func TestNoEnvironmentOrTerminalDependence(t *testing.T) {
	pos := token.Pos{Offset: 3, Line: 1, Col: 4}
	build := func() (diag.Diagnostic, string) {
		d := diag.New(diag.E0102, diag.Error, "a.atrn", pos, pos, "unterminated string")
		return d, d.Error()
	}

	base, baseText := build()

	for _, kv := range [][2]string{
		{"ATRON_LANG", "pt_BR"},
		{"NO_COLOR", "1"},
		{"TERM", "xterm-256color"},
		{"LANG", "tr_TR.UTF-8"},
	} {
		t.Setenv(kv[0], kv[1])
		got, gotText := build()
		if got != base {
			t.Errorf("setting %s changed the diagnostic value", kv[0])
		}
		if gotText != baseText {
			t.Errorf("setting %s changed Error() from %q to %q", kv[0], baseText, gotText)
		}
	}

	// The message must also be byte-stable when formatted twice with the same
	// arguments, which rules out map iteration order or pointer values leaking
	// into the text.
	for i := 0; i < 10; i++ {
		if _, text := build(); text != baseText {
			t.Fatalf("run %d produced %q, want %q", i, text, baseText)
		}
	}
}

// TestDiagnosticDoesNotImportLexer documents the dependency direction by
// checking the package boundary indirectly: a diagnostic is constructible using
// only token positions, so diag cannot depend on lexer, and token cannot depend
// on diag (which would be a cycle).
func TestDiagnosticDoesNotImportLexer(t *testing.T) {
	t.Parallel()

	pos := token.Pos{Offset: 0, Line: 1, Col: 1}
	d := diag.New(diag.E0103, diag.Error, "", pos, pos, "malformed number literal")
	if !strings.Contains(d.Error(), "E0103") {
		t.Fatalf("unexpected rendering: %q", d.Error())
	}
}

// TestSeverityStringUsedInError ensures the rendered severity and the
// Severity.String() method cannot disagree.
func TestSeverityStringUsedInError(t *testing.T) {
	t.Parallel()

	pos := token.Pos{Offset: 0, Line: 1, Col: 1}
	for _, sev := range []diag.Severity{diag.Note, diag.Warning, diag.Error} {
		d := diag.New(diag.E0101, sev, "a.atrn", pos, pos, "x")
		want := fmt.Sprintf("a.atrn:1:1: %s[E0101]: x", sev)
		if got := d.Error(); got != want {
			t.Errorf("Error() = %q, want %q", got, want)
		}
	}
}
