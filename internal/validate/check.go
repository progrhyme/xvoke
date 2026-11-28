package validate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.yaml.in/yaml/v4"
)

// check runs the checks that the schema cannot express.
func check(path string, root *yaml.Node, value any) (errs, warns []Diagnostic) {
	fields, _ := value.(map[string]any)

	if name, ok := fields["name"].(string); ok && name != "" {
		dir := directoryName(path)
		if dir != "" && name != dir {
			errs = append(errs, Diagnostic{
				Source: SourceCheck, Pointer: "/name", Line: lineOf(root, []string{"name"}),
				Message: fmt.Sprintf("name %q does not match the directory name %q", name, dir),
			})
		}
	}

	if rootDir, ok := fields["root"].(string); ok && rootDir != "" {
		if err := directoryExists(rootDir); err != nil {
			warns = append(warns, Diagnostic{
				Source: SourceCheck, Pointer: "/root", Line: lineOf(root, []string{"root"}),
				Message: err.Error(),
			})
		}
	}
	return errs, warns
}

// directoryName returns the name of the directory that holds the file, or "" if it cannot be told.
func directoryName(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return ""
	}
	return filepath.Base(filepath.Dir(abs))
}

// directoryExists reports why dir is not an existing directory. A leading "~" means the home directory.
func directoryExists(dir string) error {
	expanded, err := expandHome(dir)
	if err != nil {
		return err
	}
	info, err := os.Stat(expanded)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return fmt.Errorf("directory does not exist: %s", dir)
	case err != nil:
		return fmt.Errorf("cannot access the directory %s: %w", dir, err)
	case !info.IsDir():
		return fmt.Errorf("not a directory: %s", dir)
	}
	return nil
}

func expandHome(path string) (string, error) {
	if path != "~" && !strings.HasPrefix(path, "~/") {
		return path, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot expand ~ in %s: %w", path, err)
	}
	return filepath.Join(home, strings.TrimPrefix(path, "~")), nil
}
