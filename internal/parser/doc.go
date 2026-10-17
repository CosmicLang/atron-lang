// SPDX-License-Identifier: Apache-2.0

// Package parser builds the Atron abstract syntax tree from tokens.
//
// It will implement recursive descent with a Pratt expression parser, error
// recovery with synchronisation and a maximum of ten reported errors, as
// described in chapter 5.2 of the Atron Master Implementation Guide. Syntax
// errors are reported as E02xx diagnostics.
//
// Status: scaffolding only. The contents are implemented in milestone M2 and
// this package intentionally declares no behaviour yet.
package parser
