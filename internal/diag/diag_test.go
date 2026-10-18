// SPDX-License-Identifier: Apache-2.0

package diag_test

import (
	"testing"

	"github.com/CosmicLang/atron-lang/internal/diag"
	"github.com/CosmicLang/atron-lang/internal/token"
)

// TestLexicalCodesAreStable pins the five ratified code strings as literals.
// If a code is renamed, this test fails rather than silently changing the
// contract users and golden files depend on.
func TestLexicalCodesAreStable(t *testing.T) {
	t.Parallel()

	want := []string{"E0101", "E0102", "E0103", "E0104", "E0105"}
	if diag.E0101 != "E0101" || diag.E0102 != "E0102" || diag.E0103 != "E0103" ||
		diag.E0104 != "E0104" || diag.E0105 != "E0105" {
		t.Fatalf("lexical code constants changed: %q %q %q %q %q",
			diag.E0101, diag.E0102, diag.E0103, diag.E0104, diag.E0105)
	}

	got := diag.LexicalCodes()
	if len(got) != len(want) {
		t.Fatalf("LexicalCodes() has %d entries, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("LexicalCodes()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
	if len(got) != 5 {
		t.Fatalf("M1 owns %d lexical codes, want exactly 5", len(got))
	}
}

// TestLexicalCodesAreDistinct guards against two conditions sharing a code.
func TestLexicalCodesAreDistinct(t *testing.T) {
	t.Parallel()

	seen := make(map[string]bool)
	for _, code := range diag.LexicalCodes() {
		if seen[code] {
			t.Errorf("duplicate lexical code %q", code)
		}
		seen[code] = true
	}
}

// TestLexicalCodesAreMutable guards the copy contract: mutating the returned
// slice must not affect a later call.
func TestLexicalCodesAreMutable(t *testing.T) {
	t.Parallel()

	codes := diag.LexicalCodes()
	codes[0] = "E9999"
	if got := diag.LexicalCodes()[0]; got != diag.E0101 {
		t.Fatalf("LexicalCodes()[0] = %q after mutation, want %q", got, diag.E0101)
	}
}

// TestSeverityOrdering pins the ordering and names of the severities, which the
// renderer and any future "worst severity" summary rely on.
func TestSeverityOrdering(t *testing.T) {
	t.Parallel()

	// Severities must be strictly ordered so a summary can take the maximum.
	if diag.Note >= diag.Warning || diag.Warning >= diag.Error {
		t.Fatal("severities are not ordered Note < Warning < Error")
	}
	tests := []struct {
		sev  diag.Severity
		want string
	}{
		{diag.Note, "note"},
		{diag.Warning, "warning"},
		{diag.Error, "error"},
		{diag.Severity(200), "Severity(200)"},
	}
	for _, tc := range tests {
		if got := tc.sev.String(); got != tc.want {
			t.Errorf("Severity(%d).String() = %q, want %q", uint8(tc.sev), got, tc.want)
		}
	}
}

// TestNewPreservesEveryField checks that construction copies all eight fields
// and does not derive any of them from hidden state.
func TestNewPreservesEveryField(t *testing.T) {
	t.Parallel()

	pos := token.Pos{Offset: 4, Line: 1, Col: 5}
	end := token.Pos{Offset: 7, Line: 1, Col: 8}
	d := diag.New(diag.E0101, diag.Error, "examples/demo.atrn", pos, end, "invalid character U+00A0")

	if d.Code != diag.E0101 {
		t.Errorf("Code = %q, want %q", d.Code, diag.E0101)
	}
	if d.Severity != diag.Error {
		t.Errorf("Severity = %v, want error", d.Severity)
	}
	if d.Message != "invalid character U+00A0" {
		t.Errorf("Message = %q", d.Message)
	}
	if d.File != "examples/demo.atrn" {
		t.Errorf("File = %q", d.File)
	}
	if d.Pos != pos {
		t.Errorf("Pos = %+v, want %+v", d.Pos, pos)
	}
	if d.End != end {
		t.Errorf("End = %+v, want %+v", d.End, end)
	}
	if d.Label != "" || d.Help != "" {
		t.Errorf("Label/Help = %q/%q, want both empty", d.Label, d.Help)
	}

	// The value must be comparable and equal to an identically built one.
	if d != diag.New(diag.E0101, diag.Error, "examples/demo.atrn", pos, end, "invalid character U+00A0") {
		t.Error("two identical constructions are not equal; the diagnostic carries hidden state")
	}
}

// TestErrorfFormatsMessageOnly confirms that only the message is formatted and
// that positions are not interpolated by the format string.
func TestErrorfFormatsMessageOnly(t *testing.T) {
	t.Parallel()

	pos := token.Pos{Offset: 0, Line: 2, Col: 3}
	d := diag.Errorf(diag.E0103, "a.atrn", pos, pos, "malformed number %q at %d", "1e", 7)

	if d.Message != `malformed number "1e" at 7` {
		t.Fatalf("Message = %q", d.Message)
	}
	if d.Severity != diag.Error {
		t.Fatalf("Severity = %v, want error", d.Severity)
	}
	if d.Pos != pos || d.End != pos {
		t.Fatalf("positions not preserved: %+v %+v", d.Pos, d.End)
	}
	if d.Code != diag.E0103 || d.File != "a.atrn" {
		t.Fatalf("code/file not preserved: %q %q", d.Code, d.File)
	}
}

// TestLabelAndHelpAreSettable checks the two optional field values are part of
// the inspectable struct, so callers set them directly rather than through
// copy-on-write helpers.
func TestLabelAndHelpAreSettable(t *testing.T) {
	t.Parallel()

	pos := token.Pos{Offset: 0, Line: 1, Col: 1}
	d := diag.New(diag.E0102, diag.Error, "a.atrn", pos, pos, "unterminated string")
	d.Label = "string starts here"
	d.Help = "add a closing quote"

	if d.Label != "string starts here" {
		t.Fatalf("Label = %q", d.Label)
	}
	if d.Help != "add a closing quote" {
		t.Fatalf("Help = %q", d.Help)
	}
	if d.Error() != "a.atrn:1:1: error[E0102]: unterminated string" {
		t.Fatalf("Label/Help leaked into Error(): %q", d.Error())
	}
}
