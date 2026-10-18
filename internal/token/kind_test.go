// SPDX-License-Identifier: Apache-2.0

package token_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/CosmicLang/atron-lang/internal/token"
)

// TestKindInventoryIsExactly62 pins the ratified count (lexical.md §4) with a
// literal, so neither an added nor a removed kind can pass unnoticed.
func TestKindInventoryIsExactly62(t *testing.T) {
	t.Parallel()

	if got := token.KindCount; got != 62 {
		t.Fatalf("KindCount = %d, want exactly 62 kinds", got)
	}
	// The enumeration must be dense: the last kind sits at index 61.
	if int(token.Nil) != 61 {
		t.Fatalf("Nil = %d, want 61 (zero-based last of 62)", int(token.Nil))
	}
}

// TestIllegalIsZeroValue pins the requirement that the zero Kind is Illegal, so
// a zero-valued Token is never mistaken for a real token.
func TestIllegalIsZeroValue(t *testing.T) {
	t.Parallel()

	var zero token.Kind
	if zero != token.Illegal {
		t.Fatalf("zero Kind = %v, want Illegal", zero)
	}
	if token.Illegal != 0 {
		t.Fatalf("Illegal = %d, want 0", int(token.Illegal))
	}
}

// TestKindGroups pins the group boundaries of the inventory:
// 4 + 4 + 20 + 12 + 22 = 62, with each group starting where it must.
func TestKindGroups(t *testing.T) {
	t.Parallel()

	groups := []struct {
		name        string
		first, last token.Kind
		want        int
	}{
		{"special/identifier", token.Illegal, token.Ident, 4},
		{"literal", token.Int, token.RawString, 4},
		{"operator", token.Plus, token.DotDotEq, 20},
		{"delimiter/punctuation", token.LParen, token.FatArrow, 12},
		{"keyword", token.Const, token.Nil, 22},
	}

	total := 0
	for _, g := range groups {
		got := int(g.last) - int(g.first) + 1
		if got != g.want {
			t.Errorf("%s group spans %d kinds, want %d", g.name, got, g.want)
		}
		total += got
	}
	if total != 62 {
		t.Fatalf("group sizes sum to %d, want 62", total)
	}
}

// TestKindNamesAreCompleteAndUnique validates the name table mechanically:
// every kind from 0 to KindCount-1 has a non-empty, unique name.
func TestKindNamesAreCompleteAndUnique(t *testing.T) {
	t.Parallel()

	seen := make(map[string]token.Kind, token.KindCount)
	for i := 0; i < token.KindCount; i++ {
		k := token.Kind(i)
		name := k.String()
		if name == "" {
			t.Errorf("Kind(%d).String() is empty", i)
			continue
		}
		if strings.HasPrefix(name, "Kind(") {
			t.Errorf("Kind(%d).String() = %q, want a name from the inventory", i, name)
		}
		if prev, dup := seen[name]; dup {
			t.Errorf("duplicate kind name %q shared by %v and %v", name, prev, k)
		}
		seen[name] = k
	}
	if len(seen) != token.KindCount {
		t.Fatalf("name table has %d unique names, want %d", len(seen), token.KindCount)
	}
}

// TestKindStringOutOfRange documents that formatting an out-of-range kind is
// safe and diagnostic rather than a panic or an empty string.
func TestKindStringOutOfRange(t *testing.T) {
	t.Parallel()

	for _, k := range []token.Kind{token.Kind(62), token.Kind(200), token.Kind(255)} {
		if got, want := k.String(), fmt.Sprintf("Kind(%d)", uint8(k)); got != want {
			t.Errorf("Kind(%d).String() = %q, want %q", uint8(k), got, want)
		}
	}
}

// TestNoForbiddenKinds asserts by name that removed and never-ratified kinds
// cannot be reintroduced: there is no Let, no Main and no Comment kind.
func TestNoForbiddenKinds(t *testing.T) {
	t.Parallel()

	forbidden := []string{"Let", "Main", "Comment"}
	for i := 0; i < token.KindCount; i++ {
		name := token.Kind(i).String()
		for _, bad := range forbidden {
			if name == bad {
				t.Errorf("kind %d is named %q, which the ratified contract forbids", i, bad)
			}
		}
	}
}
