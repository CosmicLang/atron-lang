// SPDX-License-Identifier: Apache-2.0

package token

// The first kinds of the inventory. The foundation commit establishes only the
// special kinds that the position and token primitives need; the complete
// 62-kind inventory is added by the token-model work package (W1).
const (
	// Illegal is the zero value and marks input the lexer could not classify.
	Illegal Kind = iota
	// EOF marks the end of the input.
	EOF
)
