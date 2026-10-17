// SPDX-License-Identifier: Apache-2.0

// Package lexer converts Atron source text (UTF-8) into a token stream.
//
// It will implement automatic semicolon insertion (Newline tokens), number and
// string literals, raw strings, comment collection, shebang handling and CRLF
// normalisation, as described in chapter 5.1 of the Atron Master Implementation
// Guide. Lexical errors are reported as E01xx diagnostics.
//
// Status: scaffolding only. The contents are implemented in milestone M1 and
// this package intentionally declares no behaviour yet.
package lexer
