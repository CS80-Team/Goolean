package shell

import (
	"os"
	"strings"
)

func (sh *Shell) registerBuiltInCommands() {
	sh.RegisterCommand(
		NewCommand(
			"exit",
			"Exit the shell",
			"exit",
			[]Argument{},
			[]string{},
			func(s *Shell, args []string) Status {
				return EXIT
			},
			func(args []string) (bool, string) {
				return true, ""
			},
		),
	)

	sh.RegisterCommand(
		NewCommand(
			"help",
			"List all available commands",
			"help",
			[]Argument{},
			[]string{},
			func(s *Shell, args []string) Status {
				for _, cmd := range sh.GetCommands() {
					sh.Write(cmd.Name + ": " + cmd.Description + "\n")
					sh.Write("    Aliases: ")
					if len(cmd.Aliases) > 0 {
						sh.Write(strings.Join(cmd.Aliases, ", ") + "\n")
					} else {
						sh.Write("No aliases found.\n")
					}
					sh.Write("    Usage: " + cmd.Usage + "\n\n")
				}
				return OK
			},
			func(args []string) (bool, string) {
				return true, ""
			},
		),
	)

	sh.RegisterCommand(
		NewCommand(
			"clear",
			"Clear the screen",
			"clear",
			[]Argument{},
			[]string{"cls"},
			func(s *Shell, args []string) Status {
				s.clearScreen()
				return OK
			},
			func(args []string) (bool, string) {
				return true, ""
			},
		),
	)

	sh.RegisterCommand(
		NewCommand(
			"history",
			"Display the shell history",
			"history",
			[]Argument{},
			[]string{"hist"},
			func(s *Shell, args []string) Status {
				file, err := os.ReadFile(sh.historyFile)
				if err != nil {
					sh.Write("Error reading history file: " + err.Error() + "\n")
					return FAIL
				}

				sh.Write(string(file))
				return OK
			},
			func(args []string) (bool, string) {
				return true, ""
			},
		),
	)

	sh.RegisterCommand(
		NewCommand(
			"alias",
			"Create an alias for a command",
			"alias <Alias> <Command>",
			[]Argument{
				{
					Name:        "Alias",
					Description: "The alias to create",
					Required:    true,
					Type:        "string",
					Default:     "",
				},
				{
					Name:        "Command",
					Description: "The command to create an alias for",
					Required:    true,
					Type:        "string",
					Default:     "",
				},
			},
			[]string{},
			func(s *Shell, args []string) Status {
				sh.addAlias(args[0], args[1])
				return OK
			},
			func(args []string) (bool, string) {
				if len(args) != 2 {
					return false, "Invalid number of arguments"
				}

				if _, ok := sh.findCommandByNameOrAlias(args[1]); !ok {
					return false, "Command or alias not found"
				}
				return true, ""
			},
		),
	)
}
