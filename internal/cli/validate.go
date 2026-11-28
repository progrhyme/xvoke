package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/progrhyme/xvoke/internal/validate"
)

func validateCommand() command {
	return command{
		Name:    "validate",
		Summary: "Validate a project.yaml. The exit status is 0 if it is valid, 1 if it is invalid, and 2 if the command could not run.",
		Usage:   "xvoke validate <project.yaml>",
		Run:     runValidate,
	}
}

type diagnosticOutput struct {
	Source  string `json:"source"`
	Pointer string `json:"pointer"`
	Line    int    `json:"line,omitempty"`
	Message string `json:"message"`
}

type validateOutput struct {
	Valid    bool               `json:"valid"`
	Path     string             `json:"path"`
	Errors   []diagnosticOutput `json:"errors"`
	Warnings []diagnosticOutput `json:"warnings"`
}

func runValidate(args []string, e env) int {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) { // "-help", which the flag package takes as a request for help
			return e.printCommandHelp(validateCommand())
		}
		return e.fail(fmt.Sprintf("%v; %s", err, hint))
	}
	if fs.NArg() != 1 {
		return e.fail("validate takes one argument: the path of a project.yaml; " + hint)
	}
	path := fs.Arg(0)

	src, err := os.ReadFile(path)
	if err != nil {
		return e.fail(fmt.Sprintf("cannot read the file: %v", err))
	}
	res, err := validate.Validate(path, src)
	if err != nil {
		return e.fail(fmt.Sprintf("cannot validate: %v", err))
	}

	out := validateOutput{
		Valid:    res.Valid(),
		Path:     path,
		Errors:   toDiagnosticOutputs(res.Errors),
		Warnings: toDiagnosticOutputs(res.Warnings),
	}
	status := exitOK
	if !out.Valid {
		status = exitInvalid
	}
	return e.print(out, status)
}

// toDiagnosticOutputs converts diagnostics to the output form. It never returns nil, so that the output has an array.
func toDiagnosticOutputs(diags []validate.Diagnostic) []diagnosticOutput {
	out := make([]diagnosticOutput, len(diags))
	for i, d := range diags {
		out[i] = diagnosticOutput{Source: d.Source, Pointer: d.Pointer, Line: d.Line, Message: d.Message}
	}
	return out
}
