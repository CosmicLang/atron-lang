// SPDX-License-Identifier: Apache-2.0

package main

import (
	"strings"
	"testing"
)

// TestBuildInfoReport checks the rendered build information for determinism and
// complete field coverage.
func TestBuildInfoReport(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		info BuildInfo
		want []string
	}{
		{
			name: "all fields populated",
			info: BuildInfo{
				Version:   "0.1.0",
				Commit:    "0123456789abcdef",
				Date:      "2026-10-17T12:00:00Z",
				GoVersion: "go1.23.4",
				GOOS:      "darwin",
				GOARCH:    "arm64",
			},
			want: []string{
				"atron version 0.1.0",
				"0123456789abcdef",
				"2026-10-17T12:00:00Z",
				"go1.23.4",
				"darwin/arm64",
			},
		},
		{
			name: "empty fields fall back to unknown",
			info: BuildInfo{},
			want: []string{
				"atron version " + unknownValue,
				"commit:   " + unknownValue,
				"built:    " + unknownValue,
				"go:       " + unknownValue,
				"platform: /",
			},
		},
		{
			name: "blank fields fall back to unknown",
			info: BuildInfo{Version: "   ", Commit: "\t"},
			want: []string{
				"atron version " + unknownValue,
				"commit:   " + unknownValue,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.info.Report()
			for _, sub := range tc.want {
				if !strings.Contains(got, sub) {
					t.Errorf("Report() = %q, want it to contain %q", got, sub)
				}
			}

			if again := tc.info.Report(); again != got {
				t.Errorf("Report() is not deterministic: %q != %q", got, again)
			}
		})
	}
}

// TestBuildInfoReportLines locks the output shape so that scripted consumers
// relying on line prefixes keep working.
func TestBuildInfoReportLines(t *testing.T) {
	t.Parallel()

	info := BuildInfo{
		Version:   "0.0.0",
		Commit:    "deadbeef",
		Date:      "2026-01-01T00:00:00Z",
		GoVersion: "go1.23.0",
		GOOS:      "linux",
		GOARCH:    "amd64",
	}
	report := info.Report()

	lines := strings.Split(strings.TrimSuffix(report, "\n"), "\n")
	wantPrefixes := []string{
		"atron version ",
		"  commit:   ",
		"  built:    ",
		"  go:       ",
		"  platform: ",
	}
	if len(lines) != len(wantPrefixes) {
		t.Fatalf("Report() has %d lines, want %d:\n%s", len(lines), len(wantPrefixes), report)
	}
	for i, prefix := range wantPrefixes {
		if !strings.HasPrefix(lines[i], prefix) {
			t.Errorf("line %d = %q, want prefix %q", i, lines[i], prefix)
		}
	}
	if !strings.HasSuffix(report, "\n") {
		t.Error("Report() must end with a newline")
	}
}

// TestPlatform checks the os/arch rendering.
func TestPlatform(t *testing.T) {
	t.Parallel()

	tests := []struct {
		info BuildInfo
		want string
	}{
		{BuildInfo{GOOS: "linux", GOARCH: "amd64"}, "linux/amd64"},
		{BuildInfo{GOOS: "windows", GOARCH: "arm64"}, "windows/arm64"},
		{BuildInfo{}, "/"},
	}
	for _, tc := range tests {
		if got := tc.info.Platform(); got != tc.want {
			t.Errorf("Platform() = %q, want %q", got, tc.want)
		}
	}
}

// TestCurrentBuildInfo checks that the build information of the test binary is
// populated from the runtime and that the injected defaults are non-empty.
func TestCurrentBuildInfo(t *testing.T) {
	t.Parallel()

	info := currentBuildInfo()
	if info.GoVersion == "" || !strings.HasPrefix(info.GoVersion, "go") {
		t.Errorf("GoVersion = %q, want a toolchain version", info.GoVersion)
	}
	if info.GOOS == "" || info.GOARCH == "" {
		t.Errorf("platform = %q, want a non-empty os/arch", info.Platform())
	}
	if info.Version == "" || info.Commit == "" || info.Date == "" {
		t.Errorf("build metadata = %+v, want non-empty defaults", info)
	}
}

// TestOrUnknown checks the fallback helper used when metadata is missing.
func TestOrUnknown(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{"", unknownValue},
		{"  ", unknownValue},
		{"\t\n", unknownValue},
		{"0.1.0", "0.1.0"},
		{" v0.1.0 ", " v0.1.0 "},
	}
	for _, tc := range tests {
		if got := orUnknown(tc.in); got != tc.want {
			t.Errorf("orUnknown(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
