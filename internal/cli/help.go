package cli

import "slices"

func isHelpFlag(arg string) bool {
	return slices.Contains(helpFlag.Names, arg)
}

const indexUsage = "xvoke <command> [flags]"

func helpCommand() command {
	return command{
		Name:    "help",
		Summary: "Show the usage of xvoke or of a command.",
		Usage:   "xvoke help [<command>]",
		Run: func(args []string, e env) int {
			switch len(args) {
			case 0:
				return e.printIndexHelp()
			case 1:
				c, ok := lookup(e.commands, args[0])
				if !ok {
					return e.fail(unknownMessage(args[0]))
				}
				return e.printCommandHelp(c)
			default:
				return e.fail("help takes at most one command; " + hint)
			}
		},
	}
}

// wantsHelp reports whether args has -h or --help before any "--".
func wantsHelp(args []string) bool {
	for _, a := range args {
		if a == "--" {
			return false
		}
		if isHelpFlag(a) {
			return true
		}
	}
	return false
}

type indexHelpOutput struct {
	Usage    string    `json:"usage"`
	Commands []command `json:"commands"`
}

type commandHelpOutput struct {
	Command command `json:"command"`
}

// withHelpFlag returns the command with the help flag added, which every command accepts.
func withHelpFlag(c command) command {
	c.Flags = slices.Concat(c.Flags, []flagSpec{helpFlag})
	return c
}

func (e env) printIndexHelp() int {
	out := indexHelpOutput{Usage: indexUsage, Commands: make([]command, len(e.commands))}
	for i, c := range e.commands {
		out.Commands[i] = withHelpFlag(c)
	}
	return e.print(out, exitOK)
}

func (e env) printCommandHelp(c command) int {
	return e.print(commandHelpOutput{Command: withHelpFlag(c)}, exitOK)
}
