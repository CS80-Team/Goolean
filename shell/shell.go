package shell

import (
	"io"
	"os"
	"strings"
)

type Status string

const (
	OK       Status = "OK"
	FAIL     Status = "FAIL"
	EXIT     Status = "EXIT"
	NotFound Status = "NOT_FOUND"
)

type Shell struct {
	commands          map[string]Command
	earlyExecCommands []EarlyCommand
	inStream          io.Reader
	outStream         io.Writer
}

func NewShell() *Shell {
	return &Shell{
		commands:  make(map[string]Command),
		inStream:  os.Stdin,
		outStream: os.Stdout,
	}
}

func (s *Shell) RegisterCommand(cmd Command) {
	s.commands[cmd.Name] = cmd
}

func (s *Shell) RegisterEarlyExecCommand(cmd EarlyCommand) {
	s.earlyExecCommands = append(s.earlyExecCommands, cmd)
}

func (s *Shell) executeEarlyCommands() {
	for _, cmd := range s.earlyExecCommands {
		cmd.Handler(s)
	}
}

func (s *Shell) executeCommand(cmd string, args []string) Status {
	if strings.ToUpper(cmd) == string(EXIT) {
		return EXIT
	}

	if cmd == "" {
		return OK
	}

	if command, ok := s.commands[cmd]; ok {
		return command.Handler(s, args)
	}

	return NotFound
}

func (s *Shell) GetCommands() []Command {
	var cmds []Command
	for _, cmd := range s.commands {
		cmds = append(cmds, cmd)
	}
	return cmds
}

func (s *Shell) SetInputStream(in io.Reader) {
	s.inStream = in
}

func (s *Shell) SetOutputStream(out io.Writer) {
	s.outStream = out
}

func (s *Shell) read() string {
	var input string
	buf := make([]byte, 1024)
	for {
		n, err := s.inStream.Read(buf)

		if n > 0 {
			input += string(buf[:n])
		}

		if err != nil || n == 0 || buf[n-1] == '\n' {
			break
		}
	}
	return input
}

func (s *Shell) Write(output string) {
	_, _ = s.outStream.Write([]byte(output))
}

func (s *Shell) Run(welcomeMessage string) {
	var stat Status
	s.executeCommand("clear", nil)
	s.Write(welcomeMessage)
	for {
		s.Write("\n")
		s.executeEarlyCommands()
		s.Write(">")

		input := s.read()
		cmd, args := parseInput(input)
		stat = s.executeCommand(cmd, args)

		if stat == EXIT {
			break
		} else if stat == FAIL {
			s.Write(s.commands[cmd].Usage + "\n")
		} else if stat == NotFound {
			s.Write("Command not found\n")
		}
	}

	s.Exit()
}

func parseInput(input string) (string, []string) {
	tokens := strings.Fields(input)
	if len(tokens) == 0 {
		return "", nil
	}

	return tokens[0], tokens[1:]
}

func (s *Shell) Exit() {
	os.Exit(0)
}
