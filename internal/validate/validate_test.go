package validate_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/progrhyme/xvoke/internal/schema"
	"github.com/progrhyme/xvoke/internal/validate"
)

// rootToken stands for the root directory in a source, which each test creates.
const rootToken = "@root@"

// base holds the required keys. The line numbers in the tests depend on it: it takes 3 lines.
const base = "name: app\nroot: " + rootToken + "\ndefault_branch: main\n"

// projectFile returns the path of a project.yaml in a directory named name, and creates an empty directory to use as root.
func projectFile(t *testing.T, name string) (path, root string) {
	t.Helper()
	dir := t.TempDir()
	root = filepath.Join(dir, "work")
	if err := os.Mkdir(root, 0o755); err != nil {
		t.Fatal(err)
	}
	projectDir := filepath.Join(dir, "xvokedir", name)
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatal(err)
	}
	return filepath.Join(projectDir, "project.yaml"), root
}

// run validates src, with rootToken replaced by root.
func run(t *testing.T, path, root, src string) *validate.Result {
	t.Helper()
	res, err := validate.Validate(path, []byte(strings.ReplaceAll(src, rootToken, root)))
	if err != nil {
		t.Fatal(err)
	}
	return res
}

// want is the part of a diagnostic that the tests compare. Message is matched by a substring.
type want struct {
	source, pointer string
	line            int
	message         string
}

func assertDiagnostics(t *testing.T, kind string, got []validate.Diagnostic, wants []want) {
	t.Helper()
	if len(got) != len(wants) {
		t.Fatalf("%s: got %d, want %d: %+v", kind, len(got), len(wants), got)
	}
	for i, w := range wants {
		g := got[i]
		if g.Source != w.source || g.Pointer != w.pointer || g.Line != w.line || !strings.Contains(g.Message, w.message) {
			t.Errorf("%s[%d] = %+v, want %+v", kind, i, g, w)
		}
	}
}

func TestValid(t *testing.T) {
	path, root := projectFile(t, "app")
	res := run(t, path, root, base)
	if !res.Valid() {
		t.Errorf("want valid, got %+v", res)
	}
	assertDiagnostics(t, "errors", res.Errors, nil)
	assertDiagnostics(t, "warnings", res.Warnings, nil)
}

func TestSchemaErrors(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want []want
	}{
		{
			name: "missing name",
			src:  "root: " + rootToken + "\ndefault_branch: main\n",
			want: []want{{"schema", "", 1, "name"}},
		},
		{
			name: "unknown top-level key",
			src:  base + "defualt_branch: main\n",
			want: []want{{"schema", "/defualt_branch", 4, "unknown key"}},
		},
		{
			name: "wrong type, with the line of a nested key",
			src:  base + "workflow_style:\n  prefer_pull_request: \"no\"\n",
			want: []want{{"schema", "/workflow_style/prefer_pull_request", 5, "want boolean"}},
		},
		{
			name: "missing key in a gate",
			src:  base + "quality_gates:\n  - kind: ci\n  - kind: test\n",
			want: []want{{"schema", "/quality_gates/1", 6, "command"}},
		},
		{
			name: "unknown kind",
			src:  base + "quality_gates:\n  - kind: lint\n",
			want: []want{{"schema", "/quality_gates/0/kind", 5, "must be one of"}},
		},
		{
			name: "several errors are sorted by line",
			src:  base + "quality_gates:\n  - kind: review\n  - kind: lint\nworkflow_style:\n  pull_request: false\n",
			want: []want{
				{"schema", "/quality_gates/0", 5, "skill"},
				{"schema", "/quality_gates/1/kind", 6, "must be one of"},
				{"schema", "/workflow_style/pull_request", 8, "unknown key"},
			},
		},
		{
			name: "not a mapping",
			src:  "- a\n- b\n",
			want: []want{{"schema", "", 1, "object"}},
		},
		{
			name: "empty",
			src:  "",
			want: []want{{"schema", "", 0, "object"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, root := projectFile(t, "app")
			res := run(t, path, root, tt.src)
			if res.Valid() {
				t.Error("want invalid")
			}
			assertDiagnostics(t, "errors", res.Errors, tt.want)
			assertDiagnostics(t, "warnings", res.Warnings, nil)
		})
	}
}

func TestBrokenYAML(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want []want
	}{
		{name: "syntax error", src: "name: a\nfoo: [x\n", want: []want{{"check", "", 3, "invalid YAML"}}},
		{name: "more than one document", src: "name: app\n---\nname: other\n", want: []want{{"check", "", 2, "more than one document"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, root := projectFile(t, "app")
			res := run(t, path, root, tt.src)
			if res.Valid() {
				t.Error("want invalid")
			}
			assertDiagnostics(t, "errors", res.Errors, tt.want)
		})
	}
}

func TestSchemaVersion(t *testing.T) {
	if schema.Latest() != 1 {
		t.Skip("the cases assume a single schema version; update them when a new one is added")
	}
	tests := []struct {
		name    string
		version string
		want    []want
	}{
		{name: "omitted", version: ""},
		{name: "latest", version: "schema_version: 1\n"},
		{name: "not an integer", version: "schema_version: one\n", want: []want{{"check", "/schema_version", 1, "must be an integer"}}},
		{name: "not supported", version: "schema_version: 99\n", want: []want{{"check", "/schema_version", 1, "unsupported"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, root := projectFile(t, "app")
			res := run(t, path, root, tt.version+base)
			assertDiagnostics(t, "errors", res.Errors, tt.want)
		})
	}
}

func TestChecks(t *testing.T) {
	t.Run("name does not match the directory", func(t *testing.T) {
		path, root := projectFile(t, "bar")
		res := run(t, path, root, strings.Replace(base, "name: app", "name: foo", 1))
		if res.Valid() {
			t.Error("want invalid")
		}
		assertDiagnostics(t, "errors", res.Errors, []want{{"check", "/name", 1, `"foo" does not match the directory name "bar"`}})
	})
	t.Run("root does not exist is a warning", func(t *testing.T) {
		path, root := projectFile(t, "app")
		res := run(t, path, root, strings.Replace(base, rootToken, rootToken+"/missing", 1))
		if !res.Valid() {
			t.Errorf("a warning must not make the result invalid: %+v", res)
		}
		assertDiagnostics(t, "warnings", res.Warnings, []want{{"check", "/root", 2, "does not exist"}})
	})
	t.Run("root is a file", func(t *testing.T) {
		path, root := projectFile(t, "app")
		file := filepath.Join(root, "file")
		if err := os.WriteFile(file, nil, 0o644); err != nil {
			t.Fatal(err)
		}
		res := run(t, path, root, strings.Replace(base, rootToken, file, 1))
		assertDiagnostics(t, "warnings", res.Warnings, []want{{"check", "/root", 2, "not a directory"}})
	})
	t.Run("root starting with ~ is the home directory", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)
		t.Setenv("USERPROFILE", home)
		if err := os.Mkdir(filepath.Join(home, "work"), 0o755); err != nil {
			t.Fatal(err)
		}
		path, root := projectFile(t, "app")
		res := run(t, path, root, strings.Replace(base, rootToken, "~/work", 1))
		assertDiagnostics(t, "warnings", res.Warnings, nil)
		res = run(t, path, root, strings.Replace(base, rootToken, "~/other", 1))
		assertDiagnostics(t, "warnings", res.Warnings, []want{{"check", "/root", 2, "does not exist"}})
	})
	t.Run("schema errors and check results come together", func(t *testing.T) {
		path, root := projectFile(t, "bar")
		res := run(t, path, root, base+"extra: 1\n")
		assertDiagnostics(t, "errors", res.Errors, []want{
			{"schema", "/extra", 4, "unknown key"},
			{"check", "/name", 1, "does not match"},
		})
	})
}
