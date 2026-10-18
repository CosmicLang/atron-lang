// SPDX-License-Identifier: Apache-2.0

package token

// KindCount is the number of kinds in the ratified inventory.
//
// It is derived from the enumeration rather than written by hand, so the tests
// compare it against the literal 62 required by lexical.md §4.
const KindCount = int(Nil) + 1

// kindNames holds the stable name of every kind, indexed by Kind.
//
// The literal is written in enumeration order and the tests verify that it is
// complete, unique and consistent with String, so the inventory cannot drift
// silently. There is deliberately no entry for a Comment kind or a Let kind:
// neither exists in the ratified contract.
var kindNames = [...]string{
	Illegal: "Illegal",
	EOF:     "EOF",
	Newline: "Newline",
	Ident:   "Ident",

	Int:       "Int",
	Float:     "Float",
	String:    "String",
	RawString: "RawString",

	Plus:          "Plus",
	Minus:         "Minus",
	Star:          "Star",
	Slash:         "Slash",
	Percent:       "Percent",
	StarStar:      "StarStar",
	Eq:            "Eq",
	NotEq:         "NotEq",
	Lt:            "Lt",
	LtEq:          "LtEq",
	Gt:            "Gt",
	GtEq:          "GtEq",
	Assign:        "Assign",
	PlusAssign:    "PlusAssign",
	MinusAssign:   "MinusAssign",
	StarAssign:    "StarAssign",
	SlashAssign:   "SlashAssign",
	PercentAssign: "PercentAssign",
	DotDot:        "DotDot",
	DotDotEq:      "DotDotEq",

	LParen:    "LParen",
	RParen:    "RParen",
	LBracket:  "LBracket",
	RBracket:  "RBracket",
	LBrace:    "LBrace",
	RBrace:    "RBrace",
	Comma:     "Comma",
	Colon:     "Colon",
	Semicolon: "Semicolon",
	Dot:       "Dot",
	Question:  "Question",
	FatArrow:  "FatArrow",

	Const:    "Const",
	Var:      "Var",
	Fn:       "Fn",
	Return:   "Return",
	If:       "If",
	Else:     "Else",
	While:    "While",
	For:      "For",
	In:       "In",
	Break:    "Break",
	Continue: "Continue",
	Struct:   "Struct",
	Impl:     "Impl",
	Import:   "Import",
	Pub:      "Pub",
	As:       "As",
	And:      "And",
	Or:       "Or",
	Not:      "Not",
	True:     "True",
	False:    "False",
	Nil:      "Nil",
}

// IsKeyword reports whether k is one of the 22 keyword kinds.
func (k Kind) IsKeyword() bool {
	return k >= Const && k <= Nil
}

// EndsStatement reports whether a line ending after a token of kind k causes a
// Newline to be inserted after it by automatic semicolon insertion
// (lexical.md §7).
//
// The eligible set is closed and exact:
//
//	Ident, Int, Float, String, RawString,
//	True, False, Nil,
//	Return, Break, Continue,
//	RParen, RBracket, RBrace, Question
//
// Every other kind, including Newline itself and EOF, returns false. The tests
// enumerate all 62 kinds, so a kind added later cannot silently become
// statement-ending.
func (k Kind) EndsStatement() bool {
	switch k {
	case Ident, Int, Float, String, RawString,
		True, False, Nil,
		Return, Break, Continue,
		RParen, RBracket, RBrace, Question:
		return true
	default:
		return false
	}
}
