// SPDX-License-Identifier: Apache-2.0

// Package format implements the canonical Atron source formatter.
//
// It will provide the idempotent formatter used by "atron fmt", preserving
// comments and layout conventions described in chapter 5.9 of the Atron Master
// Implementation Guide. The property fmt(fmt(x)) == fmt(x) is a hard test
// requirement.
//
// Status: scaffolding only. The contents are implemented in milestone M8 and
// this package intentionally declares no behaviour yet.
package format
