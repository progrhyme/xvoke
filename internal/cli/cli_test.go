package cli_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/progrhyme/xvoke/internal/cli"
)

type result struct {
	status int
	stdout string
}

func run(t *testing.T, args ...string) result {
	t.Helper()
	var out bytes.Buffer
	status := cli.Run(args, &out, cli.Deps{Version: "1.2.3"})
	return result{status: status, stdout: out.String()}
}

// wantStatus fails the test unless the exit status is want.
func (r result) wantStatus(t *testing.T, want int) {
	t.Helper()
	if r.status != want {
		t.Fatalf("status = %d, want %d\n%s", r.status, want, r.stdout)
	}
}

// decode parses the output, and fails unless it is exactly one JSON document.
func (r result) decode(t *testing.T) map[string]any {
	t.Helper()
	dec := json.NewDecoder(strings.NewReader(r.stdout))
	var doc map[string]any
	if err := dec.Decode(&doc); err != nil {
		t.Fatalf("the output is not a JSON object: %v\n%s", err, r.stdout)
	}
	if _, err := dec.Token(); err != io.EOF {
		t.Fatalf("the output has more than one JSON document:\n%s", r.stdout)
	}
	return doc
}

func (r result) errorMessage(t *testing.T) string {
	t.Helper()
	body, ok := r.decode(t)["error"].(map[string]any)
	if !ok {
		t.Fatalf("the output has no error: %s", r.stdout)
	}
	msg, _ := body["message"].(string)
	return msg
}

// projectFile writes a project.yaml into a directory named name, and returns its path.
func projectFile(t *testing.T, name, content string) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "xvokedir", name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "project.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestHelp(t *testing.T) {
	tests := []struct {
		name string
		args []string
		// command is the name in the output of a command help. It is empty for the index.
		command string
	}{
		{name: "--help", args: []string{"--help"}},
		{name: "-h", args: []string{"-h"}},
		{name: "help", args: []string{"help"}},
		{name: "help <command>", args: []string{"help", "validate"}, command: "validate"},
		{name: "<command> --help", args: []string{"validate", "--help"}, command: "validate"},
		{name: "<command> -h", args: []string{"schema", "-h"}, command: "schema"},
		{name: "-help, which the flag package takes as help", args: []string{"validate", "-help"}, command: "validate"},
		{name: "help help", args: []string{"help", "help"}, command: "help"},
		{name: "-h after other arguments", args: []string{"validate", "project.yaml", "-h"}, command: "validate"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := run(t, tt.args...)
			r.wantStatus(t, 0)
			doc := r.decode(t)
			if tt.command == "" {
				cmds, ok := doc["commands"].([]any)
				if !ok || len(cmds) == 0 {
					t.Fatalf("the index help has no commands: %s", r.stdout)
				}
				if doc["usage"] == nil {
					t.Errorf("the index help has no usage: %s", r.stdout)
				}
				return
			}
			c, ok := doc["command"].(map[string]any)
			if !ok || c["name"] != tt.command {
				t.Fatalf("command = %v, want %q", doc["command"], tt.command)
			}
			flags, _ := c["flags"].([]any)
			if len(flags) == 0 {
				t.Errorf("the command help lists no flags, not even -h: %s", r.stdout)
			}
		})
	}
}

func TestHelpListsEveryCommand(t *testing.T) {
	doc := run(t, "help").decode(t)
	cmds, _ := doc["commands"].([]any)
	var names []string
	for _, c := range cmds {
		fields, _ := c.(map[string]any)
		name, _ := fields["name"].(string)
		names = append(names, name)
	}
	for _, want := range []string{"validate", "schema", "version", "help"} {
		if !slices.Contains(names, want) {
			t.Errorf("the index help does not list %q: %v", want, names)
		}
	}
}

func TestDashDashEndsHelpFlags(t *testing.T) {
	// A file named -h comes after "--", so it is not a request for help.
	r := run(t, "validate", "--", "-h")
	if r.status != 2 || !strings.Contains(r.errorMessage(t), "cannot read the file") {
		t.Errorf("status = %d, output = %s", r.status, r.stdout)
	}
}

func TestUsageErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "no command", args: nil, want: "no command"},
		{name: "unknown command", args: []string{"bogus"}, want: `unknown command "bogus"`},
		{name: "unknown flag as a command", args: []string{"--bogus"}, want: `unknown flag "--bogus"`},
		{name: "help of an unknown command", args: []string{"help", "bogus"}, want: `unknown command "bogus"`},
		{name: "help of two commands", args: []string{"help", "validate", "schema"}, want: "at most one"},
		{name: "validate without a path", args: []string{"validate"}, want: "one argument"},
		{name: "validate with two paths", args: []string{"validate", "a.yaml", "b.yaml"}, want: "one argument"},
		{name: "validate with an unknown flag", args: []string{"validate", "--bogus", "a.yaml"}, want: "bogus"},
		{name: "schema with an argument", args: []string{"schema", "x"}, want: "no argument"},
		{name: "version with an argument", args: []string{"version", "x"}, want: "no argument"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := run(t, tt.args...)
			r.wantStatus(t, 2)
			if msg := r.errorMessage(t); !strings.Contains(msg, tt.want) {
				t.Errorf("message = %q, want it to contain %q", msg, tt.want)
			}
		})
	}
}

func TestVersion(t *testing.T) {
	r := run(t, "version")
	r.wantStatus(t, 0)
	if got := r.decode(t)["version"]; got != "1.2.3" {
		t.Errorf("version = %v, want 1.2.3", got)
	}
}

func TestSchema(t *testing.T) {
	r := run(t, "schema")
	r.wantStatus(t, 0)
	doc := r.decode(t)
	if doc["$schema"] == nil || doc["properties"] == nil {
		t.Errorf("the output is not the schema: %s", r.stdout)
	}
}

func TestValidateValid(t *testing.T) {
	root := t.TempDir()
	path := projectFile(t, "app", "name: app\nroot: "+root+"\ndefault_branch: main\n")
	r := run(t, "validate", path)
	r.wantStatus(t, 0)
	doc := r.decode(t)
	if doc["valid"] != true || doc["path"] != path {
		t.Errorf("unexpected output: %s", r.stdout)
	}
	if !strings.Contains(r.stdout, `"errors": []`) || !strings.Contains(r.stdout, `"warnings": []`) {
		t.Errorf("errors and warnings must be empty arrays: %s", r.stdout)
	}
}

func TestValidateInvalid(t *testing.T) {
	root := t.TempDir()
	path := projectFile(t, "bar", "name: foo\nroot: "+root+"\ndefault_branch: main\nextra: 1\n")
	r := run(t, "validate", path)
	r.wantStatus(t, 1)
	doc := r.decode(t)
	if doc["valid"] != false {
		t.Errorf("valid = %v, want false", doc["valid"])
	}
	errs, _ := doc["errors"].([]any)
	if len(errs) != 2 {
		t.Fatalf("errors = %v, want 2", errs)
	}
	first := errs[0].(map[string]any)
	if first["source"] != "schema" || first["pointer"] != "/extra" || first["line"] != float64(4) || first["message"] == "" {
		t.Errorf("errors[0] = %v", first)
	}
	second := errs[1].(map[string]any)
	if second["source"] != "check" || second["pointer"] != "/name" {
		t.Errorf("errors[1] = %v", second)
	}
}

func TestValidateWarningKeepsItValid(t *testing.T) {
	path := projectFile(t, "app", "name: app\nroot: "+filepath.Join(t.TempDir(), "missing")+"\ndefault_branch: main\n")
	r := run(t, "validate", path)
	r.wantStatus(t, 0)
	warns, _ := r.decode(t)["warnings"].([]any)
	if len(warns) != 1 {
		t.Errorf("warnings = %v, want 1", warns)
	}
}

func TestValidateOmitsUnknownLine(t *testing.T) {
	path := projectFile(t, "app", "")
	r := run(t, "validate", path)
	r.wantStatus(t, 1)
	if strings.Contains(r.stdout, `"line"`) {
		t.Errorf("a diagnostic without a line must omit the key: %s", r.stdout)
	}
}

func TestValidateUnreadableFile(t *testing.T) {
	r := run(t, "validate", filepath.Join(t.TempDir(), "missing.yaml"))
	r.wantStatus(t, 2)
	if msg := r.errorMessage(t); !strings.Contains(msg, "cannot read the file") {
		t.Errorf("message = %q", msg)
	}
}
