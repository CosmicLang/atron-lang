// SPDX-License-Identifier: Apache-2.0

package diag

// Lexical diagnostic codes, fixed by lexical.md §10.
//
// Each code denotes exactly one condition. They are exported so the lexer, the
// tests and any future renderer agree on the identifiers without duplicating
// string literals, and each is documented with the position the lexer reports.
const (
	// E0101 reports an invalid character, such as an interior BOM or a code
	// point outside the identifier classes. It is reported at the offending
	// code point.
	E0101 = "E0101"
	// E0102 reports an unterminated string or raw string literal. It is
	// reported at the opening quote or backtick.
	E0102 = "E0102"
	// E0103 reports a malformed number literal or an integer literal overflow.
	// It is reported at the start of the literal. It is never reported for a
	// standalone '.', which is always valid punctuation.
	E0103 = "E0103"
	// E0104 reports an unterminated block comment. It is reported at the
	// opening "/*".
	E0104 = "E0104"
	// E0105 reports an invalid escape sequence. It is reported at the
	// backslash.
	E0105 = "E0105"
)

// LexicalCodes lists the M1 lexical codes in ascending order.
//
// It is a copy, so callers cannot change the contract through it. The tests use
// it to assert that the five codes are distinct and that no code above E0105 is
// claimed by M1.
func LexicalCodes() []string {
	return []string{E0101, E0102, E0103, E0104, E0105}
}
