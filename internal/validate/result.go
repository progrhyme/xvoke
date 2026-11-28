// Package validate checks the content of a project.yaml and reports what is wrong with it.
//
// It takes the content as bytes and returns diagnostics. Reading the file, the form of the output,
// and the exit status belong to the caller.
package validate

// The sources of a diagnostic.
const (
	// SourceSchema marks a violation of the JSON Schema.
	SourceSchema = "schema"
	// SourceCheck marks a violation found by the code, which the schema cannot express.
	SourceCheck = "check"
)

// Diagnostic is one problem found in a project.yaml.
type Diagnostic struct {
	Source string
	// Pointer is a JSON Pointer to the place of the problem. It is empty for the whole document.
	Pointer string
	// Line is the line number in the YAML. It is 0 when the place is unknown.
	Line    int
	Message string
}

// Result is the outcome of validating a project.yaml.
type Result struct {
	Errors   []Diagnostic
	Warnings []Diagnostic
}

// Valid reports whether the project.yaml has no errors. Warnings do not affect it.
func (r *Result) Valid() bool {
	return len(r.Errors) == 0
}
