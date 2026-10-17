// SPDX-License-Identifier: Apache-2.0

// Package atron exposes the public embedding API of the Atron language.
//
// It is the only public Go API of the module: everything under internal/ is
// private to the module. The package will provide Interpreter, functional
// options, Eval and Register for calling Atron from Go, as described in chapter
// 5.10 of the Atron Master Implementation Guide.
//
// Compatibility: the package will follow strict semantic compatibility from
// version 1.0 onward.
//
// Status: scaffolding only. The contents are implemented in milestone M7 and
// this package intentionally declares no behaviour yet.
package atron
