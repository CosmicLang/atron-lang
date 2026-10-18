// SPDX-License-Identifier: Apache-2.0

// Package token defines the lexical tokens produced by the Atron lexer.
//
// It provides the token kind enumeration, the source position and token
// representations, the keyword table and the automatic-semicolon-insertion
// eligibility predicate, as ratified in docs/spec/lexical.md (§4, §5, §7, §9).
//
// The package is deliberately dependency-free: it imports nothing from the
// repository, so token and diag never form a cycle. In particular this package
// must never import internal/diag.
//
// Status: M1 work package W1. The complete 62-kind inventory, the 22-keyword
// table, keyword lookup, EndsStatement and the position model are implemented
// and tested. Comment metadata is defined, but comment collection is not
// implemented yet.
package token
