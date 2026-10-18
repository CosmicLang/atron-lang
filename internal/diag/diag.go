// SPDX-License-Identifier: Apache-2.0

package diag

import (
	"fmt"

	"github.com/CosmicLang/atron-lang/internal/token"
)

// Severity classifies how serious a diagnostic is.
//
// The values are ordered from least to most severe so that a collection can be
// summarised by taking the maximum severity. Only Error is produced during M1:
// lexical problems are always fatal to compilation, even though the lexer
// itself recovers and keeps scanning.
type Severity uint8

const (
	// Note is informational and does not affect compilation.
	Note Severity = iota
	// Warning reports a probable mistake but does not stop compilation.
	Warning
	// Error reports a defect that prevents successful compilation.
	Error
)

// String returns the lower-case name of the severity as used in rendered
// diagnostics, or "Severity(n)" for a value outside the enumeration.
func (s Severity) String() string {
	switch s {
	case Note:
		return "note"
	case Warning:
		return "warning"
	case Error:
		return "error"
	default:
		return fmt.Sprintf("Severity(%d)", uint8(s))
	}
}

// Diagnostic is a single reported problem.
//
// It is a concrete, inspectable value: every field is exported and nothing is
// hidden behind global mutable state, so a diagnostic can be constructed,
// compared and formatted deterministically in tests and by callers.
//
// Pos is the start of the offending span and End is exclusive, so a renderer
// can highlight exactly the bytes in [Pos.Offset, End.Offset). Both are in the
// file named by File, which the lexer passes through unchanged; it is never
// resolved against the filesystem.
type Diagnostic struct {
	// Code is the stable diagnostic identifier, such as "E0103".
	Code string
	// Severity is how serious the problem is.
	Severity Severity
	// Message is the primary, one-line description of the problem.
	Message string
	// File is the source file the diagnostic refers to.
	File string
	// Pos is the position of the first byte of the offending span.
	Pos token.Pos
	// End is the position just past the offending span, exclusive.
	End token.Pos
	// Label is the short annotation attached to the highlighted span.
	Label string
	// Help is an optional suggestion for fixing the problem.
	Help string
}

// Message formatting is fixed here so that diagnostics are reproducible across
// platforms, locales and terminal configurations: no width, colour or
// translation is involved at this layer.

// New constructs a diagnostic for the named file and span.
//
// The file name and both positions are explicit parameters, so construction
// depends on nothing but its arguments: no environment variable, terminal
// capability, clock or filesystem is consulted. Label and Help start empty and
// can be set by the caller or through WithLabel and WithHelp.
func New(code string, sev Severity, file string, pos, end token.Pos, msg string) Diagnostic {
	return Diagnostic{
		Code:     code,
		Severity: sev,
		Message:  msg,
		File:     file,
		Pos:      pos,
		End:      end,
	}
}

// Errorf constructs an error-severity diagnostic with a formatted message.
//
// The format string follows fmt conventions; the arguments fill only the
// message. Positions and the file name are always supplied explicitly, which
// keeps the result independent of any hidden state.
func Errorf(code, file string, pos, end token.Pos, format string, args ...any) Diagnostic {
	return New(code, Error, file, pos, end, fmt.Sprintf(format, args...))
}

// WithLabel and WithHelp are intentionally not provided: the eight contract
// fields are exported and can be set directly, so no copy-on-write helpers are
// needed, and fewer methods keep the M1 diagnostic model minimal.

// Error implements the error interface so a diagnostic can be returned
// directly. The rendering is deterministic and stable: it is
//
//	file:line:col: error[CODE]: message
//
// with the "file:" part omitted when File is empty. It deliberately contains no
// colour, no snippet and no locale-dependent wording; those belong to the M6
// renderer.
//
// The receiver is a value rather than a pointer because Diagnostic itself must
// satisfy error: a pointer receiver would make only *Diagnostic match, and the
// lexer stores plain values. The value is copied once per rendering, which is
// the intended cost of keeping diagnostics inspectable.
//
//nolint:gocritic // hugeParam: a value receiver is required for Diagnostic to satisfy error.
func (d Diagnostic) Error() string {
	var prefix string
	if d.File != "" {
		prefix = d.File + ":" + d.Pos.String() + ": "
	} else {
		prefix = d.Pos.String() + ": "
	}
	return prefix + d.Severity.String() + "[" + d.Code + "]: " + d.Message
}
