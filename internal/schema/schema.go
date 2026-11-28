// Package schema holds the versioned JSON Schemas of project.yaml.
package schema

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"regexp"
	"strconv"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"go.yaml.in/yaml/v4"
)

// ErrUnknownVersion means that no schema exists for the requested version.
var ErrUnknownVersion = errors.New("unknown schema version")

//go:embed project.schema.v*.yaml
var files embed.FS

var fileNamePattern = regexp.MustCompile(`^project\.schema\.v([0-9]+)\.yaml$`)

func fileName(version int) string {
	return fmt.Sprintf("project.schema.v%d.yaml", version)
}

// Get returns the schema of the version, in the YAML form.
func Get(version int) ([]byte, error) {
	b, err := files.ReadFile(fileName(version))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("schema version %d: %w", version, ErrUnknownVersion)
	}
	if err != nil {
		return nil, fmt.Errorf("read schema version %d: %w", version, err)
	}
	return b, nil
}

// Latest returns the newest version among the embedded schemas.
func Latest() int {
	entries, err := fs.ReadDir(files, ".")
	if err != nil {
		panic(err) // reading "." of an embedded file system does not fail
	}
	latest := 0
	for _, e := range entries {
		m := fileNamePattern.FindStringSubmatch(e.Name())
		if m == nil {
			continue
		}
		v, err := strconv.Atoi(m[1])
		if err != nil {
			continue // an overflowing number is not a version
		}
		latest = max(latest, v)
	}
	return latest
}

// Compile returns the compiled schema of the version, or an error wrapping ErrUnknownVersion if none exists.
func Compile(version int) (*jsonschema.Schema, error) {
	b, err := Get(version)
	if err != nil {
		return nil, err
	}
	var doc any
	if err := yaml.Unmarshal(b, &doc); err != nil {
		return nil, fmt.Errorf("parse schema version %d: %w", version, err)
	}
	url := "xvoke://" + fileName(version)
	c := jsonschema.NewCompiler()
	if err := c.AddResource(url, doc); err != nil {
		return nil, fmt.Errorf("add schema version %d: %w", version, err)
	}
	s, err := c.Compile(url)
	if err != nil {
		return nil, fmt.Errorf("compile schema version %d: %w", version, err)
	}
	return s, nil
}
