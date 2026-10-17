// SPDX-License-Identifier: Apache-2.0

// Package compiler translates the Atron abstract syntax tree into bytecode.
//
// It will consume the AST plus the resolver symbol table and emit
// *object.Function values with chunks, as described in chapter 5.5 of the Atron
// Master Implementation Guide. It is not a tree-walking interpreter.
//
// Status: scaffolding only. The contents are implemented in milestone M4 and
// this package intentionally declares no behaviour yet.
package compiler
