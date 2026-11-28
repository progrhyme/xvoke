package cli

import "slices"

// flagSpec describes a flag of a command. It is also the form of a flag in the help.
type flagSpec struct {
	Names       []string `json:"names"`
	Description string   `json:"description"`
}

// command is the definition of a command. It is also the form of a command in the help.
type command struct {
	Name    string                         `json:"name"`
	Summary string                         `json:"summary"`
	Usage   string                         `json:"usage"`
	Flags   []flagSpec                     `json:"flags"`
	Run     func(args []string, e env) int `json:"-"`
}

var helpFlag = flagSpec{Names: []string{"-h", "--help"}, Description: "Show the usage of the command."}

// commands returns the commands, in the order that the help lists them.
// Adding a command here is enough: the help is built from this list.
func commands() []command {
	return []command{validateCommand(), schemaCommand(), versionCommand(), helpCommand()}
}

// lookup returns the command named name, and whether it exists.
func lookup(cmds []command, name string) (command, bool) {
	i := slices.IndexFunc(cmds, func(c command) bool { return c.Name == name })
	if i < 0 {
		return command{}, false
	}
	return cmds[i], true
}
