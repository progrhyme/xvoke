package schema_test

import (
	"errors"
	"os"
	"strings"
	"testing"

	"go.yaml.in/yaml/v4"

	"github.com/progrhyme/xvoke/internal/schema"
)

const examplePath = "../../plugin/skills/xv-run/references/project.example.yaml"

func TestGet(t *testing.T) {
	tests := []struct {
		name    string
		version int
		wantErr error
	}{
		{name: "v1", version: 1},
		{name: "unknown version", version: 99, wantErr: schema.ErrUnknownVersion},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := schema.Get(tt.version)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && len(b) == 0 {
				t.Error("the schema is empty")
			}
		})
	}
}

func TestLatest(t *testing.T) {
	if got := schema.Latest(); got < 1 {
		t.Fatalf("Latest = %d, want 1 or more", got)
	}
	if _, err := schema.Get(schema.Latest()); err != nil {
		t.Errorf("Get(Latest()) = %v", err)
	}
}

func TestCompile(t *testing.T) {
	if _, err := schema.Compile(schema.Latest()); err != nil {
		t.Fatalf("Compile(Latest()) = %v", err)
	}
	if _, err := schema.Compile(99); !errors.Is(err, schema.ErrUnknownVersion) {
		t.Errorf("Compile(99) = %v, want ErrUnknownVersion", err)
	}
}

func decode(t *testing.T, src string) any {
	t.Helper()
	var doc any
	if err := yaml.Unmarshal([]byte(src), &doc); err != nil {
		t.Fatalf("parse YAML: %v", err)
	}
	return doc
}

func validate(t *testing.T, src string) error {
	t.Helper()
	s, err := schema.Compile(schema.Latest())
	if err != nil {
		t.Fatal(err)
	}
	return s.Validate(decode(t, src))
}

// TestExampleIsValid checks that project.example.yaml passes the latest schema.
func TestExampleIsValid(t *testing.T) {
	b, err := os.ReadFile(examplePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := validate(t, string(b)); err != nil {
		t.Errorf("project.example.yaml does not pass the schema: %v", err)
	}
}

const base = "name: app\nroot: ~/repos/app\ndefault_branch: main\n"

func TestValid(t *testing.T) {
	tests := []struct {
		name string
		src  string
	}{
		{name: "required keys only", src: base},
		{name: "schema_version 1", src: "schema_version: 1\n" + base},
		{name: "prefer_pull_request false", src: base + "workflow_style:\n  prefer_pull_request: false\n"},
		{name: "review by an agent, authorizer omitted", src: base + "quality_gates:\n  - kind: review\n    skill: code-review\n"},
		{name: "review by an agent, authorizer agent", src: base + "quality_gates:\n  - kind: review\n    authorizer: agent\n    skill: code-review\n"},
		{name: "review by the user", src: base + "quality_gates:\n  - kind: review\n    authorizer: user\n    description: Check it\n"},
		{name: "ci without options", src: base + "quality_gates:\n  - kind: ci\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validate(t, tt.src); err != nil {
				t.Errorf("want valid, got: %v", err)
			}
		})
	}
}

func TestInvalid(t *testing.T) {
	tests := []struct {
		name string
		src  string
		// want is a part of the error, so that the input fails for the reason it names.
		want string
	}{
		{name: "unknown top-level key", src: base + "defualt_branch: main\n", want: "defualt_branch"},
		{name: "missing name", src: "root: ~/repos/app\ndefault_branch: main\n", want: "missing property 'name'"},
		{name: "missing root", src: "name: app\ndefault_branch: main\n", want: "missing property 'root'"},
		{name: "missing default_branch", src: "name: app\nroot: ~/repos/app\n", want: "missing property 'default_branch'"},
		{name: "schema_version is not supported", src: "schema_version: 2\n" + base, want: "/schema_version"},
		{name: "schema_version is not an integer", src: "schema_version: one\n" + base, want: "/schema_version"},
		{name: "unknown key in workflow_style", src: base + "workflow_style:\n  prefer_branch: true\n", want: "prefer_branch"},
		{name: "prefer_pull_request is not a boolean", src: base + "workflow_style:\n  prefer_pull_request: \"false\"\n", want: "/workflow_style/prefer_pull_request"},
		{name: "unknown kind", src: base + "quality_gates:\n  - kind: lint\n", want: "/quality_gates/0/kind"},
		{name: "gate without kind", src: base + "quality_gates:\n  - command: make test\n", want: "missing property 'kind'"},
		{name: "test without command", src: base + "quality_gates:\n  - kind: test\n", want: "missing property 'command'"},
		{name: "consistency without paths", src: base + "quality_gates:\n  - kind: consistency\n    description: Check it\n", want: "missing property 'paths'"},
		{name: "consistency without description", src: base + "quality_gates:\n  - kind: consistency\n    paths: [\"docs/**\"]\n", want: "missing property 'description'"},
		{name: "review by an agent without skill", src: base + "quality_gates:\n  - kind: review\n", want: "missing property 'skill'"},
		{name: "review by the user without description", src: base + "quality_gates:\n  - kind: review\n    authorizer: user\n", want: "missing property 'description'"},
		{name: "unknown authorizer", src: base + "quality_gates:\n  - kind: review\n    authorizer: robot\n    skill: code-review\n", want: "/quality_gates/0/authorizer"},
		{name: "unknown key in a gate", src: base + "quality_gates:\n  - kind: test\n    command: make test\n    skill: code-review\n", want: "additional properties 'skill'"},
		{name: "paths is not a list", src: base + "quality_gates:\n  - kind: ci\n    paths: \"docs/**\"\n", want: "/quality_gates/0/paths"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate(t, tt.src)
			if err == nil {
				t.Fatal("want invalid, got valid")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("the error does not contain %q: %v", tt.want, err)
			}
		})
	}
}
