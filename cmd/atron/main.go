// SPDX-License-Identifier: Apache-2.0

// Command atron is the command-line interface of the Atron programming
// language.
//
// At milestone M0 the command exposes only build information and help output.
// Language execution commands (run, check, fmt, test, repl, disasm, tokens and
// ast) are specified in the Atron Master Implementation Guide, chapter 5.9, and
// will be added in later milestones.
package main

import (
	"os"
)

func main() {
	app := &app{
		stdout: os.Stdout,
		stderr: os.Stderr,
		info:   currentBuildInfo(),
	}
	os.Exit(app.run(os.Args[1:]))
}
