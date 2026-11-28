package cli

import "github.com/progrhyme/xvoke/internal/schema"

func schemaCommand() command {
	return command{
		Name:    "schema",
		Summary: "Show the latest JSON Schema of project.yaml.",
		Usage:   "xvoke schema",
		Run: func(args []string, e env) int {
			if msg := noArgs("schema", args); msg != "" {
				return e.fail(msg)
			}
			doc, err := schema.Document(schema.Latest())
			if err != nil {
				return e.fail(err.Error())
			}
			return e.print(doc, exitOK)
		},
	}
}
