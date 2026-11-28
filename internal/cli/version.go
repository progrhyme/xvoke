package cli

type versionOutput struct {
	Version string `json:"version"`
}

func versionCommand() command {
	return command{
		Name:    "version",
		Summary: "Show the version of xvoke.",
		Usage:   "xvoke version",
		Run: func(args []string, e env) int {
			if msg := noArgs("version", args); msg != "" {
				return e.fail(msg)
			}
			return e.print(versionOutput{Version: e.deps.Version}, exitOK)
		},
	}
}
