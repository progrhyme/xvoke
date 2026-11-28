// Package cli runs the xvoke command: it parses the arguments, runs the command,
// and writes the result to the standard output as one JSON document.
//
// The standard error output is not used. A failure to run is reported as {"error":{"message":...}}.
package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// The exit statuses.
const (
	exitOK      = 0 // the command succeeded, and the content is valid
	exitInvalid = 1 // the command ran, and found the content invalid
	exitError   = 2 // the command could not run: a wrong usage, or an unreadable file
)

const hint = "run `xvoke help` for the usage"

// Deps is what the caller provides to Run.
type Deps struct {
	// Version is the version of xvoke.
	Version string
}

// env is what a command gets when it runs.
type env struct {
	stdout   io.Writer
	deps     Deps
	commands []command
}

// Run runs the command that args (without the program name) names, and returns the exit status.
func Run(args []string, stdout io.Writer, deps Deps) int {
	e := env{stdout: stdout, deps: deps, commands: commands()}
	if len(args) == 0 {
		return e.fail("no command given; " + hint)
	}
	name := args[0]
	if isHelpFlag(name) {
		return e.printIndexHelp()
	}
	c, ok := lookup(e.commands, name)
	if !ok {
		return e.fail(unknownMessage(name))
	}
	if wantsHelp(args[1:]) {
		return e.printCommandHelp(c)
	}
	return c.Run(args[1:], e)
}

// print writes v as a single-line JSON document, and returns status. A failure to write is an error.
func (e env) print(v any, status int) int {
	enc := json.NewEncoder(e.stdout)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		// The output cannot be written, and nothing may go to the standard error output, so only the status tells it.
		return exitError
	}
	return status
}

type errorOutput struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Message string `json:"message"`
}

// fail writes the error document, and returns exitError.
func (e env) fail(message string) int {
	return e.print(errorOutput{Error: errorBody{Message: message}}, exitError)
}

// noArgs returns the error message if a command that takes no argument got some, or "" if it got none.
func noArgs(name string, args []string) string {
	if len(args) == 0 {
		return ""
	}
	return name + " takes no argument; " + hint
}

func unknownMessage(name string) string {
	kind := "command"
	if strings.HasPrefix(name, "-") {
		kind = "flag"
	}
	return fmt.Sprintf("unknown %s %q; %s", kind, name, hint)
}
