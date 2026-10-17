// SPDX-License-Identifier: Apache-2.0

// Package vm executes Atron bytecode on a stack virtual machine.
//
// It will provide the execution loop, call frames, reentrant calls for native
// functions and context-aware cancellation, as described in chapter 5.6 of the
// Atron Master Implementation Guide. Execution errors are reported as R00xx
// runtime errors with a stack trace.
//
// The VM is designed to be instantiated per execution: the package must not
// introduce mutable global execution state, so that several VMs can coexist in
// one process.
//
// Status: scaffolding only. The contents are implemented from milestone M4 and
// this package intentionally declares no behaviour yet.
package vm
