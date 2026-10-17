// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"io"
)

// Process exit codes returned by the command-line interface.
//
// The values come from chapter 4.8 of the Atron Master Implementation Guide,
// which fixes the meaning of every exit status the toolchain may produce:
//
//	0   successful execution
//	1   unrecoverable runtime error or panic
//	2   compilation diagnostics
//	64  incorrect command-line usage
//
// Codes 1 and 2 are reserved for the language implementation and are therefore
// not declared here: no code path in this milestone can produce them, and
// declaring unused constants would misrepresent how much of the contract is
// implemented. They arrive with the milestones that own them (M4 for runtime
// faults, M6 for compilation diagnostics).
const (
	// exitOK indicates that the requested command completed successfully.
	exitOK = 0
	// exitUsage indicates incorrect command-line usage, such as an unknown
	// command, an unexpected argument or no command at all. The value is 64
	// (EX_USAGE in sysexits.h), not 2, because 2 is reserved for compilation
	// diagnostics.
	exitUsage = 64
)

// usageErrorPrefix is the prefix used for messages written to standard error.
const usageErrorPrefix = "atron: "

// helpHint points the user at the top-level help text after a usage error.
const helpHint = "Run 'atron help' for usage."

// app holds the dependencies of a single command-line invocation.
//
// The app carries no package-level mutable state: every execution creates its
// own instance, which keeps concurrent invocations independent and testable.
type app struct {
	stdout io.Writer
	stderr io.Writer
	info   BuildInfo
}

// run dispatches args (the program arguments without the executable name) and
// returns the process exit code.
func (a *app) run(args []string) int {
	if len(args) == 0 {
		printUsage(a.stderr)
		return exitUsage
	}

	switch args[0] {
	case "version", "-v", "--version":
		if len(args) > 1 {
			return a.unexpectedArgs(args[0], args[1:])
		}
		a.printVersion()
		return exitOK
	case "help", "-h", "--help":
		if len(args) > 1 {
			return a.printCommandHelp(args[1])
		}
		printUsage(a.stdout)
		return exitOK
	default:
		fmt.Fprintf(a.stderr, "%sunknown command %q\n\n", usageErrorPrefix, args[0])
		fmt.Fprintln(a.stderr, helpHint)
		return exitUsage
	}
}

// unexpectedArgs reports positional arguments that a command does not accept.
func (a *app) unexpectedArgs(command string, extra []string) int {
	fmt.Fprintf(a.stderr, "%sunexpected argument %q for %q\n\n",
		usageErrorPrefix, extra[0], command)
	fmt.Fprintln(a.stderr, helpHint)
	return exitUsage
}

// printCommandHelp prints the help text of a single command.
func (a *app) printCommandHelp(name string) int {
	text, ok := commandHelp[name]
	if !ok {
		fmt.Fprintf(a.stderr, "%sunknown command %q\n\n", usageErrorPrefix, name)
		fmt.Fprintln(a.stderr, helpHint)
		return exitUsage
	}
	fmt.Fprint(a.stdout, text)
	return exitOK
}

// printUsage writes the top-level usage text to w.
func printUsage(w io.Writer) {
	fmt.Fprint(w, usageText)
}

// printVersion writes the build information to standard output.
func (a *app) printVersion() {
	fmt.Fprint(a.stdout, a.info.Report())
}

// commandHelp maps a command name to its help text.
var commandHelp = map[string]string{
	"version": `atron version — print build information

Usage:
  atron version

Prints the version, commit, build date, Go version and target platform of the
atron executable. Version, commit and build date are injected at link time by
the build (see the Makefile).

Exit status:
  0   success
  64  unexpected argument
`,
	"help": `atron help — print help for the toolchain or a command

Usage:
  atron help [command]

Without arguments, prints the top-level usage text. With a command name, prints
the help text of that command.

Exit status:
  0   success
  64  unknown command
`,
}

// usageText is the top-level usage text of the atron command.
const usageText = `atron — the Atron programming language toolchain

Usage:
  atron <command> [arguments]

Commands:
  version   Print build information (version, commit, build date, Go version)
  help      Print this help text, or help for a single command

Options:
  -h, --help      Print this help text
  -v, --version   Print build information

Exit status:
  0   success
  1   unrecoverable runtime error or panic (reserved; not produced at M0)
  2   compilation diagnostics (reserved; not produced at M0)
  64  incorrect command-line usage

Atron is at milestone M0 (repository foundation); only build information is
available. Execution commands (run, check, fmt, test, repl, disasm, tokens,
ast) arrive in later milestones. See docs/PROGRESS.md for the current status.
`

// reportFormat is the layout used by [BuildInfo.Report]. It is a constant so
// the rendered output is fixed at compile time and cannot drift at run time.
const reportFormat = "atron version %s\n" +
	"  commit:   %s\n" +
	"  built:    %s\n" +
	"  go:       %s\n" +
	"  platform: %s\n"
