package shell

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"

	"github.com/CS80-Team/Goolean/internal"

	"github.com/chzyer/readline"
)

type Status string

const (
	OK        Status = "OK"
	FAIL      Status = "FAIL"
	EXIT      Status = "EXIT"
	NOT_FOUND Status = "NOT_FOUND"
)

const (
	SHELL_PROMPT   = ">>> "
	SHELL_PREFIX   = "[SHELL]: "
	COMMAND_PROMPT = "[COMMAND]: "
)

const (
	COLOR_RED     = "\033[31m"
	COLOR_GREEN   = "\033[32m"
	COLOR_YELLOW  = "\033[33m"
	COLOR_BLUE    = "\033[34m"
	COLOR_MAGENTA = "\033[35m"
	COLOR_CYAN    = "\033[36m"
	COLOR_RESET   = "\033[0m"
)

type Shell struct {
	commands          map[string]*Command
	rootCommand       map[string]string
	earlyExecCommands []EarlyCommand
	inStream          io.Reader
	outStream         io.Writer
	inputHandler      *InputHandler
	prompt            string
	historyFile       string
	logger            *internal.Logger
}

func NewShell(istream io.Reader, ostream io.Writer, prompt string, historyFile string, logger *internal.Logger) *Shell {
	sh := &Shell{
		commands:    make(map[string]*Command),
		inStream:    istream,
		outStream:   ostream,
		prompt:      prompt,
		historyFile: historyFile,
		rootCommand: make(map[string]string),
		logger:      logger,
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

func (sh *Shell) addAlias(alias string, cmd string) {
	if cm, ok := sh.commands[cmd]; ok {
		sh.logger.GetLogger().Warn(fmt.Sprintf("Alias %s already exists for command %s\n", alias, cm.Name))
	}
	sh.rootCommand[alias] = cmd
	// sh.commands[cmd].AddAlias(alias)
}

func (sh *Shell) RegisterCommand(cmd *Command) {
	for _, alias := range cmd.Aliases {
		sh.addAlias(alias, cmd.Name)
	}

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

func (sh *Shell) autoCompleteCommand(cmd string) (string, bool) {
	for c := range sh.commands {
		if strings.HasPrefix(c, cmd) {
			return c, true
		}
	}

	return "", false
}

func (sh *Shell) autoCompleteArg(cmd, argPrefix string) (string, bool) {
	if command, ok := sh.findCommandByNameOrAlias(cmd); ok {
		for _, arg := range command.Args {
			if strings.HasPrefix(arg.Name, argPrefix) {
				return arg.Name, true
			}
		}
	}

	return "", false
}

func (sh *Shell) executeCommand(cmdOrAlias string, args []string) Status {
	if strings.ToUpper(cmdOrAlias) == string(EXIT) {
		return EXIT
	}

	if cmdOrAlias == "" {
		return OK
	}

	if command, ok := sh.findCommandByNameOrAlias(cmdOrAlias); ok {
		ok, err := command.ValidateArgs(args)
		if !ok {
			sh.Write(COMMAND_PROMPT + "Invalid arguments, " + err + "\n")
			sh.logger.GetLogger().Error(fmt.Sprintf("Invalid arguments for command %s: %s", cmdOrAlias, err))
			return FAIL
		}
		return command.Handler(sh, args)
	}

	sh.logger.GetLogger().Error(fmt.Sprintf("Command %s not found\n", cmdOrAlias))
	return NOT_FOUND
}

func (sh *Shell) findCommandByNameOrAlias(cmdOrAlias string) (*Command, bool) {
	if command, ok := sh.commands[cmdOrAlias]; ok {
		return command, true
	}

	if command, ok := sh.rootCommand[cmdOrAlias]; ok {
		return sh.commands[command], true
	}
	return &Command{}, false
}

func (sh *Shell) GetCommands() []*Command {
	var cmds []*Command
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

func (sh *Shell) WriteColored(color string, output string) {
	sh.Write(color + output + COLOR_RESET)
}

func (sh *Shell) Write(output string) {
	_, _ = sh.outStream.Write([]byte(output))
}

func (sh *Shell) Run(welcMessage string) {
	var stat Status

	sh.clearScreen()
	sh.Write(welcMessage)
	sh.sortEarlyCommands()

	for {
		sh.Write("\n")
		sh.executeEarlyCommands()

		input, err := sh.inputHandler.ReadLine()
		if err != nil {
			if errors.Is(err, readline.ErrInterrupt) {
				// TODO: End running command (run the command in a goroutine)
				continue
			}

			sh.Write("Error reading input: " + err.Error() + "\n")
			continue
		}

		commandOrAlias, args := sh.parseInput(input)
		stat = sh.executeCommand(commandOrAlias, args)

		if stat == EXIT {
			break
		} else if stat == FAIL {
			command, found := sh.findCommandByNameOrAlias(commandOrAlias)
			if !found {
				sh.handleCommandOrAliasNotFound(commandOrAlias)
			} else {
				sh.Write(sh.prompt + "Command failed, Usage: " + command.Usage + "\n")
			}
		} else if stat == NOT_FOUND {
			sh.handleCommandOrAliasNotFound(commandOrAlias)
		}
	}

	sh.Exit()
}

func (sh *Shell) Exit() {
	_ = sh.logger.Close()
}

func (sh *Shell) sortEarlyCommands() {
	sort.SliceStable(sh.earlyExecCommands, func(i, j int) bool {
		return sh.earlyExecCommands[i].Priority > sh.earlyExecCommands[j].Priority
	})
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
