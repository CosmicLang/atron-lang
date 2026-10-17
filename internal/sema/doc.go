// SPDX-License-Identifier: Apache-2.0

// Package sema performs semantic analysis of Atron programs.
//
// It will implement the resolver: scopes, symbols, hoisting, capture analysis
// (*Cell bindings) and validations E0301–E0308, as described in chapter 5.3 of
// the Atron Master Implementation Guide.
//
// Status: scaffolding only. The contents are implemented in milestone M3 and
// this package intentionally declares no behaviour yet.
package sema
