// SPDX-License-Identifier: Apache-2.0

package token

// kindNames holds the stable name of every kind, indexed by Kind.
//
// It grows with the inventory: the token-model work package (W1) extends it to
// all 62 kinds, and the tests verify completeness at that point.
var kindNames = [...]string{
	Illegal: "Illegal",
	EOF:     "EOF",
}
