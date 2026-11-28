package validate

import (
	"errors"
	"strconv"
	"strings"

	"go.yaml.in/yaml/v4"
)

var pointerEscaper = strings.NewReplacer("~", "~0", "/", "~1")

// pointerOf joins tokens into a JSON Pointer, escaping "~" and "/" (RFC 6901).
func pointerOf(tokens []string) string {
	var sb strings.Builder
	for _, t := range tokens {
		sb.WriteByte('/')
		sb.WriteString(pointerEscaper.Replace(t))
	}
	return sb.String()
}

// lineOf returns the line of the key (for a mapping entry) or of the item (for a sequence) that the location points to.
// If the location leads out of the document, it returns the line of the deepest node found.
// It returns 0 for an empty document.
func lineOf(root *yaml.Node, location []string) int {
	node := root
	line := node.Line
	for _, token := range location {
		switch node.Kind {
		case yaml.MappingNode:
			next := -1
			for i := 0; i+1 < len(node.Content); i += 2 {
				if node.Content[i].Value == token {
					next = i
					break
				}
			}
			if next < 0 {
				return line
			}
			line = node.Content[next].Line
			node = node.Content[next+1]
		case yaml.SequenceNode:
			index, err := strconv.Atoi(token)
			if err != nil || index < 0 || index >= len(node.Content) {
				return line
			}
			node = node.Content[index]
			line = node.Line
		default:
			return line
		}
	}
	return line
}

// yamlDiagnostic turns an error of the YAML parser into a diagnostic, with the line when the parser tells it.
func yamlDiagnostic(err error) *Diagnostic {
	var le *yaml.LoadError
	if errors.As(err, &le) {
		return &Diagnostic{Source: SourceCheck, Line: le.Mark.Line, Message: "invalid YAML: " + le.Message}
	}
	return &Diagnostic{Source: SourceCheck, Message: "invalid YAML: " + err.Error()}
}
