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
	SHELL_PREFIX   = "SHELL"
	COMMAND_PREFIX = "COMMAND"
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
	inStream          io.ReadCloser
	inStreamWriter    io.Writer
	outStream         io.Writer
	errStream         io.Writer
	inputHandler      *InputHandler
	prompt            string
	historyFile       string
	logger            *internal.Logger
}

func NewShell(
	istream io.ReadCloser,
	inStreamWriter io.Writer,
	ostream io.Writer,
	errStream io.Writer,
	prompt string,
	historyFile string,
	logger *internal.Logger,
) *Shell {
	sh := &Shell{
		commands:    make(map[string]*Command),
		inStream:    istream,
		outStream:   ostream,
		prompt:      prompt,
		historyFile: historyFile,
		rootCommand: make(map[string]string),
		logger:      logger,
	}

	listener := &KeyListener{shell: sh}
	inputHandler, err := NewInputHandler(
		prompt,
		historyFile,
		listener,
		istream,
		inStreamWriter,
		ostream,
		ostream,
	)

	if err != nil {
		panic(err)
	}
	sh.inputHandler = inputHandler

	sh.registerBuiltInCommands()
	return sh
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

func (sh *Shell) GetCommands() []*Command {
	var cmds []*Command
	for _, cmd := range sh.commands {
		cmds = append(cmds, cmd)
	}
	return cmds
}

func (sh *Shell) SetInputStream(in io.ReadCloser) {
	sh.inStream = in
}

func (sh *Shell) SetInputStreamWriter(in io.Writer) {
	sh.inStreamWriter = in
}

func (sh *Shell) SetOutputStream(out io.Writer) {
	sh.outStream = out
}

func (sh *Shell) SetErrorStream(err io.Writer) {
	sh.errStream = err
}

func (sh *Shell) Error(prefix, err string) {
	prefix = "[" + prefix + " Error]: " + err + "\n"
	sh.WriteColored(COLOR_RED, prefix)
}

func (sh *Shell) Warn(prefix, warning string) {
	prefix = "[" + prefix + " Warning]: " + warning + "\n"
	sh.WriteColored(COLOR_YELLOW, prefix)
}

func (sh *Shell) Info(prefix, info string) {
	prefix = "[" + prefix + " Info]: " + info + "\n"
	sh.WriteColored(COLOR_BLUE, prefix)
}

func (sh *Shell) Success(prefix, success string) {
	prefix = "[" + prefix + " Success]: " + success + "\n"
	sh.WriteColored(COLOR_GREEN, prefix)
}

func (sh *Shell) WriteColored(color string, output string) {
	if !isTerminal(sh.outStream) {
		sh.Write(output)
		return
	}
	sh.Write(string(color) + output + string(COLOR_RESET))
}

func (sh *Shell) Write(output string) {
	_, _ = sh.outStream.Write([]byte(output))
}

func (sh *Shell) Run(welcMessage string) {
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

			if errors.Is(err, io.EOF) { // Ctrl+D to exit
				break
			}

			sh.Error(SHELL_PREFIX, "Error reading input: "+err.Error())
			continue
		}

		if sh.execute(&input) == EXIT {
			break
		}
	}

	sh.Exit()
}

func (sh *Shell) Exit() {
	_ = sh.logger.Close()
	sh.inputHandler.Close()
	sh.inStream.Close()
}

/*

- Private methods

*/

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

	errMsg := "Command (" + cmd + ") not found, "

	if nearestCmd != "" {
		if matchedAlias != "" {
			errMsg += "did you mean `" + matchedAlias + "` (alias for `" + nearestCmd + "`)?, "
		} else {
			errMsg += "did you mean `" + nearestCmd + "`?, "
		}
	}
	errMsg += "type `help` for list of commands"
	sh.Error(COMMAND_PREFIX, errMsg)
	sh.logger.GetLogger().Error(errMsg)
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

func (sh *Shell) parseInput(input *string) (string, []string) {
	tokens := strings.Fields(*input)
	if len(tokens) == 0 {
		return "", nil
	}

	return tokens[0], tokens[1:]
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
			if arg.Tag != EMPTY_TAG {
				if strings.HasPrefix(arg.Tag, argPrefix) {
					return arg.Tag, true
				}
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
			sh.Error(COMMAND_PREFIX, "Invalid arguments, "+err)
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

func (sh *Shell) execute(input *string) Status {
	commandOrAlias, args := sh.parseInput(input)

	switch sh.executeCommand(commandOrAlias, args) {
	case EXIT:
		return EXIT
	case FAIL:
		command, found := sh.findCommandByNameOrAlias(commandOrAlias)
		if !found {
			sh.handleCommandOrAliasNotFound(commandOrAlias)
		} else {
			sh.Error(SHELL_PREFIX, "Command failed, Usage: "+command.Usage)
		}
	case NOT_FOUND:
		sh.handleCommandOrAliasNotFound(commandOrAlias)
	}

	return OK
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
	if exsistCmd, ok := sh.rootCommand[alias]; ok {
		warn := fmt.Sprintf("Alias %s already exists for command %s, alias overrided.", alias, sh.commands[exsistCmd].Name)
		sh.Warn(COMMAND_PREFIX, warn)
		sh.logger.GetLogger().Warn(warn)
	}
	sh.rootCommand[alias] = cmd
}
