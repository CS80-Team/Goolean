package shell

import (
	"io"
	"os"
	"strings"
)

type Status string

const (
	OK   Status = "OK"
	FAIL Status = "FAIL"
	EXIT Status = "EXIT"
)

type shell struct {
	commands  map[string]Command
	inStream  io.Reader
	outStream io.Writer
}

func NewShell() *shell {
	return &shell{
		commands: make(map[string]Command),
	}
}

func (s *shell) RegisterCommand(cmd Command) {
	s.commands[cmd.Name] = cmd
}

func (s *shell) Help() {
	for _, cmd := range s.GetCommands() {
		s.WriteOutput(cmd.Name + ": " + cmd.Description + "\n")
		s.WriteOutput("\tUsage: " + cmd.Usage + "\n")
	}
}

func (s *shell) ExecuteCommand(cmd string, args []string) Status {
	if strings.ToUpper(cmd) == string(EXIT) {
		return EXIT
	}

	if cmd == "" {
		return OK
	}

	if cmd == "help" {
		s.Help()
		return OK
	}

	if command, ok := s.commands[cmd]; ok {
		return command.Handler(args)
	}

	s.WriteOutput("Command not found\n")
	return FAIL
}

func (s *shell) GetCommands() []Command {
	var cmds []Command
	for _, cmd := range s.commands {
		cmds = append(cmds, cmd)
	}
	return cmds
}

func (s *shell) SetInputStream(in io.Reader) {
	s.inStream = in
}

func (s *shell) SetOutputStream(out io.Writer) {
	s.outStream = out
}

func (s *shell) ReadInput() string {
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

func (s *shell) WriteOutput(output string) {
	s.outStream.Write([]byte(output))
}

func (s *shell) Run() {
	var stat Status
	for {
		s.WriteOutput(">")
		input := s.ReadInput()
		cmd, args := parseInput(input)
		stat = s.ExecuteCommand(cmd, args)
		if stat == EXIT {
			break
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

func (s *shell) Exit() {
	os.Exit(0)
}
