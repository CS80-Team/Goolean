package shell

import (
	"os"
	"os/exec"
	"path/filepath"
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
			"alias <alias> <command>",
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
				cmd, _ := sh.findCommandByNameOrAlias(args[1])
				sh.addAlias(args[0], cmd.Name)
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

	sh.RegisterCommand(
		NewCommand(
			"exec",
			"Execute a an external command",
			"exec <command>",
			[]Argument{
				{
					Name:        "Execute command",
					Description: "The command to execute",
					Required:    true,
					Type:        "string",
					Default:     "",
				},
			},
			[]string{},
			func(s *Shell, args []string) Status {
				cmd := exec.Command(args[0], args[1:]...)

				cmd.Stdout = sh.outStream
				cmd.Stderr = sh.errStream
				cmd.Stdin = sh.inStream

                err := cmd.Run()

				if err != nil {
					sh.Write("Error executing command: " + err.Error() + "\n")
					return FAIL
				}

				return OK
			},
			func(args []string) (bool, string) {
                if len(args) < 1 {
                    return false, "No command provided"
                }

                return true, ""
            },
		),
	)

	sh.RegisterCommand(
		NewCommand(
			"run",
			"Run a script",
			"run <script_path.shell>",
			[]Argument{
				{
					Name:        "Script Path",
					Description: "The path to the script to run",
					Required:    true,
					Type:        "string",
					Default:     "",
				},
			},
			[]string{},
			func(s *Shell, args []string) Status {
				file, err := os.ReadFile(args[0])
				if err != nil {
					sh.Write("Error reading script file: " + err.Error() + "\n")
					return FAIL
				}

				lines := strings.Split(string(file), "\n")
				for _, line := range lines {
					sh.execute(&line)
				}

				return OK
			},
			func(args []string) (bool, string) {
				if len(args) != 1 {
					return false, "Invalid number of arguments"
				}

				if info, err := os.Stat(args[0]); os.IsNotExist(err) {
					return false, "File not found"
				} else if info.IsDir() || filepath.Ext(args[0]) != ".shell" {
					return false, "Invalid file type"
				}

				return true, ""
			},
		),
	)
}
