package shell

import (
	"errors"
	"github.com/chzyer/readline"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Status string

const (
	OK        Status = "OK"
	FAIL      Status = "FAIL"
	EXIT      Status = "EXIT"
	NOT_FOUND Status = "NOT_FOUND"
)

const (
	SHELL_PROMPT = ">>> "
	SHELL_PREFIX = "[SHELL]: "
)

type Shell struct {
	commands          map[string]Command
	earlyExecCommands []EarlyCommand
	inStream          io.Reader
	outStream         io.Writer
	inputHandler      *InputHandler
	prompt            string
	historyFile       string
}

func NewShell(istream io.Reader, ostream io.Writer) *Shell {
	sh := &Shell{
		commands:    make(map[string]Command),
		inStream:    istream,
		outStream:   ostream,
		prompt:      SHELL_PROMPT,
		historyFile: ".shell_history",
	}

	// Initialize the InputHandler
	listener := &KeyListener{shell: sh}
	inputHandler, err := NewInputHandler(sh.prompt, sh.historyFile, listener)
	if err != nil {
		panic(err)
	}
	sh.inputHandler = inputHandler

	sh.registerBuiltInCommands()

	return sh
}

func (sh *Shell) registerBuiltInCommands() {
	sh.RegisterCommand(Command{
		Name:        "exit",
		Description: "Exit the shell",
		Handler: func(s *Shell, args []string) Status {
			return EXIT
		},
		Usage: "exit",
	})

	sh.RegisterCommand(Command{
		Name:        "help",
		Description: "List all available commands",
		Handler: func(s *Shell, args []string) Status {
			for _, cmd := range sh.GetCommands() {
				sh.Write(cmd.Name + ": " + cmd.Description + "\n")
				sh.Write("    Usage: " + cmd.Usage + "\n\n")
			}
			return OK
		},
		Usage: "help",
	})

	sh.RegisterCommand(Command{
		Name:        "clear",
		Aliases:     []string{"cls"},
		Description: "Clear the screen",
		Handler: func(s *Shell, args []string) Status {
			s.clearScreen()
			return OK
		},
		Usage: "clear",
	})
}

func (sh *Shell) clearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	}
	if sh.inputHandler.reader != nil {
		sh.inputHandler.reader.Refresh()
	}
}

func (sh *Shell) RegisterCommand(cmd Command) {
	sh.commands[cmd.Name] = cmd
}

func (sh *Shell) RegisterEarlyExecCommand(cmd EarlyCommand) {
	sh.earlyExecCommands = append(sh.earlyExecCommands, cmd)
}

func (sh *Shell) executeEarlyCommands() {
	for _, cmd := range sh.earlyExecCommands {
		cmd.Handler(sh)
	}
}

func (sh *Shell) executeCommand(cmd string, args []string) Status {
	if strings.ToUpper(cmd) == string(EXIT) {
		return EXIT
	}

	if cmd == "" {
		return OK
	}

	if command, ok := sh.findCommandByNameOrAlias(cmd); ok {
		return command.Handler(sh, args)
	}

	return NOT_FOUND
}

func (sh *Shell) findCommandByNameOrAlias(cmd string) (Command, bool) {
	if command, ok := sh.commands[cmd]; ok {
		return command, true
	}

	for _, command := range sh.commands {
		for _, alias := range command.Aliases {
			if alias == cmd {
				return command, true
			}
		}
	}

	return Command{}, false
}

func (sh *Shell) GetCommands() []Command {
	var cmds []Command
	for _, cmd := range sh.commands {
		cmds = append(cmds, cmd)
	}
	return cmds
}

func (sh *Shell) SetInputStream(in io.Reader) {
	sh.inStream = in
}

func (sh *Shell) SetOutputStream(out io.Writer) {
	sh.outStream = out
}

func (sh *Shell) read() string {
	var input string
	buf := make([]byte, 1024)
	for {
		n, err := sh.inStream.Read(buf)

		if n > 0 {
			input += string(buf[:n])
		}

		if err != nil || n == 0 || buf[n-1] == '\n' {
			break
		}
	}
	return input
}

func (sh *Shell) Write(output string) {
	_, _ = sh.outStream.Write([]byte(output))
}

func (sh *Shell) Run(welcMessage string) {
	var stat Status

	sh.clearScreen()
	sh.Write(welcMessage)
	for {
		sh.Write("\n")
		sh.executeEarlyCommands()

		input, err := sh.inputHandler.reader.Readline()
		if err != nil {
			if errors.Is(err, readline.ErrInterrupt) {
				continue
			} else if err == io.EOF {
				sh.Write("\n")
				break
			}
			sh.Write("Error reading input: " + err.Error() + "\n")
			continue
		}

		cmd, args := sh.parseInput(input)
		stat = sh.executeCommand(cmd, args)

		if stat == EXIT {
			break
		} else if stat == FAIL {
			sh.Write(SHELL_PREFIX + sh.commands[cmd].Usage + "\n")
		} else if stat == NOT_FOUND {
			sh.handleCommandOrAliasNotFound(cmd)
		}
	}
}

func (sh *Shell) handleCommandOrAliasNotFound(cmd string) {
	nearestCmd, matchedAlias := sh.getNearestCommandOrAlias(cmd)
	if len(cmd) > 20 {
		cmd = cmd[:20] + "..."
	}

	sh.Write("Command (" + cmd + ") not found, ")

	if nearestCmd != "" {
		if matchedAlias != "" {
			sh.Write("did you mean `" + matchedAlias + "` (alias for `" + nearestCmd + "`)?, ")
		} else {
			sh.Write("did you mean `" + nearestCmd + "`?, ")
		}
	}
	sh.Write("type `help` for list of commands\n")
}

func (sh *Shell) getNearestCommandOrAlias(cmd string) (string, string) {
	best := 2
	nearestCmd := ""
	matchedAlias := ""

	for _, c := range sh.commands {
		dist := editDistance(c.Name, cmd)
		if dist <= best {
			best = dist
			nearestCmd = c.Name
			matchedAlias = ""
		}
		// Also check aliases
		for _, alias := range c.Aliases {
			dist := editDistance(alias, cmd)
			if dist <= best {
				best = dist
				nearestCmd = c.Name
				matchedAlias = alias
			}
		}
	}
	return nearestCmd, matchedAlias
}

func (sh *Shell) parseInput(input string) (string, []string) {
	tokens := strings.Fields(input)
	if len(tokens) == 0 {
		return "", nil
	}

	return tokens[0], tokens[1:]
}
