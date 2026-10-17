// SPDX-License-Identifier: Apache-2.0

// Package stdlib implements the native modules of the Atron standard library.
//
// It will provide the math, strings, os, fs, time, json, rand and testing
// modules described in chapter 5.7 of the Atron Master Implementation Guide.
// Modules are constructed per interpreter; the package must not keep mutable
// global state.
//
// Status: scaffolding only. The contents are implemented in milestone M7 and
// this package intentionally declares no behaviour yet.
package stdlib
