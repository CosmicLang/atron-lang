// SPDX-License-Identifier: Apache-2.0

// Package module loads, caches and resolves Atron source modules.
//
// It will implement import resolution, the module cache, cycle detection
// (E0307) and ATRON_PATH handling, as described in chapter 5.7 of the Atron
// Master Implementation Guide. Paths are normalised with filepath.Clean so that
// behaviour is consistent across platforms.
//
// Status: scaffolding only. The contents are implemented in milestone M7 and
// this package intentionally declares no behaviour yet.
package module
