// Command yup-sed is the CLI wrapper around github.com/gloo-foo/cmd-sed.
package main

import (
	clix "github.com/gloo-foo/cli"
	command "github.com/gloo-foo/cmd-sed"
)

// version is the build version. It defaults to "dev" for local builds and is
// overridden at release time via the linker: -ldflags "-X main.version=<v>".
var version = "dev"

const name = "sed"

// Error is the package's sentinel error type, so every emitted error path is
// comparable with errors.Is.
type Error string

func (e Error) Error() string { return string(e) }

// ErrMissingScript is emitted when no SCRIPT operand is supplied.
const ErrMissingScript Error = "missing SCRIPT operand"

// synopsis is the multi-line --help usage block; urfave/cli indents it three
// spaces, so the lines stay flush-left.
const synopsis = `sed SCRIPT

Apply the s/pattern/replacement/[flags] SCRIPT to standard
input, writing the result to standard output.`

// spec declares the sed wrapper: a stdin filter whose single operand is the
// substitution script.
var spec = clix.Spec{
	Name:     name,
	Summary:  "stream editor for filtering and transforming text",
	Synopsis: synopsis,
	Build:    build,
}

// build maps the invocation to sed's pipeline: standard input feeds sed, whose
// single operand is the script. A bare invocation is a usage error.
func build(inv clix.Invocation) (clix.Source, clix.Command, error) {
	operands := inv.Args.Args().Slice()
	if len(operands) == 0 {
		return nil, nil, ErrMissingScript
	}
	script := command.SedScript(operands[0])
	return clix.Stdin(inv.Stdin), command.Sed(script), nil
}

// runMain is an indirection seam so main's wiring is testable without spawning
// the process; a test swaps it and restores it.
var runMain = clix.Main

func main() { runMain(spec, version) }
