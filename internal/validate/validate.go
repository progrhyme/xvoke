package validate

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/santhosh-tekuri/jsonschema/v6/kind"
	"go.yaml.in/yaml/v4"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/progrhyme/xvoke/internal/schema"
)

const schemaVersionKey = "schema_version"

// Validate checks the content of the project.yaml at path.
// The problems of the content go to the result. The error is for the cases where the check itself cannot run.
func Validate(path string, src []byte) (*Result, error) {
	res := &Result{}

	root, value, diag := parse(src)
	if diag != nil {
		res.Errors = append(res.Errors, *diag)
		return res, nil
	}

	diags, err := schemaDiagnostics(root, value)
	if err != nil {
		return nil, err
	}
	res.Errors = append(res.Errors, diags...)

	errs, warns := check(path, root, value)
	res.Errors = append(res.Errors, errs...)
	res.Warnings = append(res.Warnings, warns...)
	return res, nil
}

// parse reads the single YAML document of src, and returns its root node and its value.
// An empty src gives an empty node and a nil value.
// A syntax error, or a second document, comes back as a diagnostic.
func parse(src []byte) (*yaml.Node, any, *Diagnostic) {
	dec := yaml.NewDecoder(bytes.NewReader(src))
	var doc yaml.Node
	if err := dec.Decode(&doc); err != nil {
		if errors.Is(err, io.EOF) {
			return &doc, nil, nil
		}
		return nil, nil, yamlDiagnostic(err)
	}
	var second yaml.Node
	switch err := dec.Decode(&second); {
	case err == nil:
		return nil, nil, &Diagnostic{
			Source: SourceCheck, Line: second.Line,
			Message: "invalid YAML: the file has more than one document",
		}
	case !errors.Is(err, io.EOF):
		return nil, nil, yamlDiagnostic(err)
	}
	var value any
	if err := doc.Decode(&value); err != nil {
		return nil, nil, yamlDiagnostic(err)
	}
	root := &doc
	if doc.Kind == yaml.DocumentNode && len(doc.Content) > 0 {
		root = doc.Content[0]
	}
	return root, value, nil
}

// schemaDiagnostics validates the value against the schema that schema_version selects.
// It reads schema_version without the schema, because the schema depends on it.
// When the key is omitted, it uses the latest version.
// The error is for a failure to use the schema; a wrong schema_version comes back as a diagnostic.
func schemaDiagnostics(root *yaml.Node, value any) ([]Diagnostic, error) {
	version := schema.Latest()
	fields, _ := value.(map[string]any)
	if raw, ok := fields[schemaVersionKey]; ok {
		v, isInt := raw.(int)
		if !isInt {
			return versionDiagnostic(root, schemaVersionKey+" must be an integer"), nil
		}
		version = v
	}
	s, err := schema.Compile(version)
	if errors.Is(err, schema.ErrUnknownVersion) {
		msg := fmt.Sprintf("unsupported %s %d; the latest is %d", schemaVersionKey, version, schema.Latest())
		return versionDiagnostic(root, msg), nil
	}
	if err != nil {
		return nil, err
	}
	return violations(s, root, value)
}

func versionDiagnostic(root *yaml.Node, msg string) []Diagnostic {
	return []Diagnostic{{
		Source: SourceCheck, Pointer: "/" + schemaVersionKey, Line: lineOf(root, []string{schemaVersionKey}),
		Message: msg,
	}}
}

// violations validates the value against the schema, and returns a diagnostic for each violation.
func violations(s *jsonschema.Schema, root *yaml.Node, value any) ([]Diagnostic, error) {
	verr := s.Validate(value)
	if verr == nil {
		return nil, nil
	}
	var ve *jsonschema.ValidationError
	if !errors.As(verr, &ve) {
		return nil, fmt.Errorf("validate against the schema: %w", verr)
	}
	printer := message.NewPrinter(language.English)
	var diags []Diagnostic
	collect(ve, root, printer, &diags)
	slices.SortStableFunc(diags, func(a, b Diagnostic) int {
		return cmp.Or(
			cmp.Compare(a.Line, b.Line),
			cmp.Compare(a.Pointer, b.Pointer),
			strings.Compare(a.Message, b.Message),
		)
	})
	return slices.Compact(diags), nil
}

// collect appends a diagnostic for each leaf of the validation error tree.
// The inner nodes only group the leaves, and say nothing about the cause.
func collect(ve *jsonschema.ValidationError, root *yaml.Node, printer *message.Printer, diags *[]Diagnostic) {
	if len(ve.Causes) > 0 {
		for _, c := range ve.Causes {
			collect(c, root, printer, diags)
		}
		return
	}
	// A key that is not allowed is reported at the key, which is the place to fix.
	if k, ok := ve.ErrorKind.(*kind.AdditionalProperties); ok {
		for _, prop := range k.Properties {
			location := append(slices.Clone(ve.InstanceLocation), prop)
			*diags = append(*diags, Diagnostic{
				Source: SourceSchema, Pointer: pointerOf(location), Line: lineOf(root, location),
				Message: fmt.Sprintf("unknown key %q", prop),
			})
		}
		return
	}
	*diags = append(*diags, Diagnostic{
		Source:  SourceSchema,
		Pointer: pointerOf(ve.InstanceLocation),
		Line:    lineOf(root, ve.InstanceLocation),
		Message: ve.ErrorKind.LocalizedString(printer),
	})
}
