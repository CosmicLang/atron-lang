// SPDX-License-Identifier: Apache-2.0

// Package diag defines the Atron diagnostic model.
//
// It provides the severity enumeration, the eight-field Diagnostic value, a
// deterministic constructor and the stable lexical codes E0101–E0105, as
// ratified in docs/spec/lexical.md §10. Diagnostics are plain values with
// exported fields: they can be inspected, compared and formatted in tests
// without touching a terminal, an environment variable, the clock or the
// filesystem.
//
// The package depends on internal/token for source positions. The dependency is
// one-way: internal/token never imports this package, so no cycle exists and a
// Token stays printable without the diagnostic model.
//
// Status: M1 work packages W0–W2. Severity, Diagnostic, New, Errorf and the
// lexical codes are implemented and tested. Rendering, ANSI colour,
// localisation, message catalogues, source snippets and the ATRON_LANG and
// NO_COLOR behaviour belong to M6 and are deliberately absent.
package diag
