// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"runtime"
	"strings"
)

// Build metadata injected with -ldflags by the build (see the Makefile).
//
// These are the only package-level variables in the command. They are written
// exclusively by the Go linker and are never mutated at run time, so they do not
// constitute mutable global execution state.
var (
	// version is the semantic version of the build, for example "0.1.0".
	version = "devel"
	// commit is the Git commit the binary was built from.
	commit = "unknown"
	// date is the RFC 3339 build timestamp.
	date = "unknown"
)

// unknownValue replaces empty build metadata so that output stays stable when
// the binary is produced by a plain `go build`.
const unknownValue = "unknown"

// BuildInfo describes the build of the atron executable.
//
// A BuildInfo is a plain value; it is safe to copy and has no hidden state.
type BuildInfo struct {
	// Version is the semantic version injected at link time.
	Version string
	// Commit is the Git commit hash injected at link time.
	Commit string
	// Date is the build timestamp injected at link time.
	Date string
	// GoVersion is the version of the Go toolchain used to build the binary.
	GoVersion string
	// GOOS is the target operating system.
	GOOS string
	// GOARCH is the target architecture.
	GOARCH string
}

// currentBuildInfo returns the build information of the running executable.
func currentBuildInfo() BuildInfo {
	return BuildInfo{
		Version:   orUnknown(version),
		Commit:    orUnknown(commit),
		Date:      orUnknown(date),
		GoVersion: runtime.Version(),
		GOOS:      runtime.GOOS,
		GOARCH:    runtime.GOARCH,
	}
}

// Platform returns the target platform in the form "os/arch".
func (b *BuildInfo) Platform() string {
	return b.GOOS + "/" + b.GOARCH
}

// Report renders the build information as human-readable text.
//
// The output is deterministic for a given binary: it contains no timestamps of
// the current time, absolute paths or memory addresses.
func (b *BuildInfo) Report() string {
	return fmt.Sprintf(reportFormat,
		orUnknown(b.Version),
		orUnknown(b.Commit),
		orUnknown(b.Date),
		orUnknown(b.GoVersion),
		b.Platform(),
	)
}

// orUnknown returns s, or [unknownValue] when s is empty.
func orUnknown(s string) string {
	if strings.TrimSpace(s) == "" {
		return unknownValue
	}
	return s
}
