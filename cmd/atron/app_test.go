// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestExitCodeContract pins the process exit codes to the values fixed by
// chapter 4.8 of the Atron Master Implementation Guide. The expectations are
// written as literals on purpose: if someone changes a constant, this test must
// fail rather than follow the change silently.
func TestExitCodeContract(t *testing.T) {
	t.Parallel()

	if exitOK != 0 {
		t.Errorf("exitOK = %d, want 0 (successful execution)", exitOK)
	}
	if exitUsage != 64 {
		t.Errorf("exitUsage = %d, want 64 (incorrect CLI usage, EX_USAGE)", exitUsage)
	}
}

// TestAppRun covers the public behavior of the command-line interface: the
// selected output stream, the exit code and the key content of each message.
func TestAppRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		args       []string
		wantCode   int
		wantStdout []string
		wantStderr []string
		noStdout   bool
		noStderr   bool
	}{
		{
			name:       "no arguments prints usage to stderr",
			args:       nil,
			wantCode:   64,
			wantStderr: []string{"Usage:", "atron <command> [arguments]", "version", "help"},
			noStdout:   true,
		},
		{
			name:       "help prints usage to stdout",
			args:       []string{"help"},
			wantCode:   0,
			wantStdout: []string{"Usage:", "atron <command> [arguments]", "version", "help"},
			noStderr:   true,
		},
		{
			name:       "long help flag behaves like help",
			args:       []string{"--help"},
			wantCode:   0,
			wantStdout: []string{"Usage:"},
			noStderr:   true,
		},
		{
			name:       "short help flag behaves like help",
			args:       []string{"-h"},
			wantCode:   0,
			wantStdout: []string{"Usage:"},
			noStderr:   true,
		},
		{
			name:       "help for version command",
			args:       []string{"help", "version"},
			wantCode:   0,
			wantStdout: []string{"atron version", "Usage:", "atron version"},
			noStderr:   true,
		},
		{
			name:       "help for unknown command fails",
			args:       []string{"help", "bogus"},
			wantCode:   64,
			wantStderr: []string{`unknown command "bogus"`, helpHint},
			noStdout:   true,
		},
		{
			name:       "version prints build information",
			args:       []string{"version"},
			wantCode:   0,
			wantStdout: []string{"atron version ", "commit:", "built:", "go:", "platform:"},
			noStderr:   true,
		},
		{
			name:       "version accepts long flag",
			args:       []string{"--version"},
			wantCode:   0,
			wantStdout: []string{"atron version "},
			noStderr:   true,
		},
		{
			name:       "version accepts short flag",
			args:       []string{"-v"},
			wantCode:   0,
			wantStdout: []string{"atron version "},
			noStderr:   true,
		},
		{
			name:       "version rejects extra arguments",
			args:       []string{"version", "extra"},
			wantCode:   64,
			wantStderr: []string{`unexpected argument "extra"`, helpHint},
			noStdout:   true,
		},
		{
			name:       "unknown command fails",
			args:       []string{"bogus"},
			wantCode:   64,
			wantStderr: []string{`unknown command "bogus"`, helpHint},
			noStdout:   true,
		},
		{
			name:       "reserved language command is not implemented in M0",
			args:       []string{"run", "hello.atrn"},
			wantCode:   64,
			wantStderr: []string{`unknown command "run"`, helpHint},
			noStdout:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var stdout, stderr bytes.Buffer
			a := &app{
				stdout: &stdout,
				stderr: &stderr,
				info: BuildInfo{
					Version:   "1.2.3",
					Commit:    "abc1234",
					Date:      "2026-10-17T00:00:00Z",
					GoVersion: "go1.23.0",
					GOOS:      "linux",
					GOARCH:    "amd64",
				},
			}

			got := a.run(tc.args)

			if got != tc.wantCode {
				t.Errorf("run(%q) exit code = %d, want %d", tc.args, got, tc.wantCode)
			}
			assertStream(t, "stdout", stdout.String(), tc.wantStdout, tc.noStdout)
			assertStream(t, "stderr", stderr.String(), tc.wantStderr, tc.noStderr)
		})
	}
}

// assertStream checks that got contains every want substring and, when
// wantEmpty is set, that it is empty.
func assertStream(t *testing.T, name, got string, want []string, wantEmpty bool) {
	t.Helper()

	if wantEmpty {
		if got != "" {
			t.Errorf("%s = %q, want empty", name, got)
		}
		return
	}
	for _, sub := range want {
		if !strings.Contains(got, sub) {
			t.Errorf("%s = %q, want it to contain %q", name, got, sub)
		}
	}
}

// TestUsageMentionsOnlyImplementedCommands guards the M0 milestone boundary:
// the help text must not advertise language execution commands that do not
// exist yet.
func TestUsageMentionsOnlyImplementedCommands(t *testing.T) {
	t.Parallel()

	listCommands := []string{"run", "check", "fmt", "test", "repl", "disasm", "tokens", "ast"}
	for _, cmd := range listCommands {
		if hasCommandLine(usageText, cmd) {
			t.Errorf("usage text advertises unimplemented command %q", cmd)
		}
	}
	if !hasCommandLine(usageText, "version") {
		t.Error("usage text does not list the version command")
	}
	if !hasCommandLine(usageText, "help") {
		t.Error("usage text does not list the help command")
	}
}

// hasCommandLine reports whether text contains a command table entry that
// starts with name.
func hasCommandLine(text, name string) bool {
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), name+" ") {
			return true
		}
	}
	return false
}

// TestCompiledBinaryExitCodes builds the real command and checks the process
// exit status. The tests above verify what run returns; this one verifies that
// os.Exit propagates that value to the operating system, which is the
// observable contract a shell script depends on.
func TestCompiledBinaryExitCodes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping build-and-exec test in short mode")
	}

	bin := buildTestBinary(t, "")

	tests := []struct {
		name     string
		args     []string
		wantCode int
	}{
		{name: "no arguments", args: nil, wantCode: 64},
		{name: "unknown command", args: []string{"bogus"}, wantCode: 64},
		{name: "unknown command from help", args: []string{"help", "bogus"}, wantCode: 64},
		{name: "unexpected argument", args: []string{"version", "extra"}, wantCode: 64},
		{name: "reserved run command", args: []string{"run", "x.atrn"}, wantCode: 64},
		{name: "help", args: []string{"help"}, wantCode: 0},
		{name: "help flag", args: []string{"--help"}, wantCode: 0},
		{name: "version", args: []string{"version"}, wantCode: 0},
		{name: "version flag", args: []string{"--version"}, wantCode: 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(bin, tc.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()

			got := 0
			if err != nil {
				var exitErr *exec.ExitError
				if !errors.As(err, &exitErr) {
					t.Fatalf("running the binary with %q failed: %v", tc.args, err)
				}
				got = exitErr.ExitCode()
			}
			if got != tc.wantCode {
				t.Errorf("args %q: exit status = %d, want %d (stderr: %q)",
					tc.args, got, tc.wantCode, stderr.String())
			}
			if tc.wantCode == exitOK && stderr.Len() != 0 {
				t.Errorf("args %q: wrote %q to stderr on success", tc.args, stderr.String())
			}
		})
	}
}

// TestCompiledBinaryVersionMetadata checks that metadata injected with -ldflags
// reaches the binary. It also guards the linker symbol name: package main is
// recorded as "main", so "main.version" works while the module path does not.
func TestCompiledBinaryVersionMetadata(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping build-and-exec test in short mode")
	}

	const (
		wantVersion = "9.9.9-test"
		wantCommit  = "0123456789abcdef0123456789abcdef01234567"
		wantDate    = "2026-10-17T00:00:00Z"
	)

	ldflags := "-X 'main.version=" + wantVersion + "'" +
		" -X 'main.commit=" + wantCommit + "'" +
		" -X 'main.date=" + wantDate + "'"
	bin := buildTestBinary(t, ldflags)

	out, err := exec.Command(bin, "version").Output()
	if err != nil {
		t.Fatalf("running the version command failed: %v", err)
	}
	for _, want := range []string{wantVersion, wantCommit, wantDate} {
		if !strings.Contains(string(out), want) {
			t.Errorf("version output %q does not contain %q", out, want)
		}
	}
}

// buildTestBinary compiles the command into a temporary directory and returns
// its path. When ldflags is not empty it is passed to the build, which lets
// tests exercise the metadata injection path.
func buildTestBinary(t *testing.T, ldflags string) string {
	t.Helper()

	bin := filepath.Join(t.TempDir(), "atron")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}

	args := []string{"build", "-trimpath"}
	if ldflags != "" {
		args = append(args, "-ldflags", ldflags)
	}
	args = append(args, "-o", bin, ".")

	build := exec.Command("go", args...)
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		t.Fatalf("building the atron command failed: %v", err)
	}
	return bin
}
