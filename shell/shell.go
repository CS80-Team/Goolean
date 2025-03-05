package shell

import (
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

type shell struct {
	commands          map[string]Command
	earlyExecCommands []EarlyCommand
	inStream          io.Reader
	outStream         io.Writer
}

func NewShell() *shell {
	sh := &shell{
		commands:  make(map[string]Command),
		inStream:  os.Stdin,
		outStream: os.Stdout,
	}

	sh.RegisterCommand(Command{
		Name:        "exit",
		Description: "Exit the shell",
		Handler: func(args []string) Status {
			return EXIT
		},
		Usage: "exit",
	})

	sh.RegisterCommand(Command{
		Name:        "help",
		Description: "List all available commands",
		Handler: func(args []string) Status {
			for _, cmd := range sh.GetCommands() {
				sh.Write(cmd.Name + ": " + cmd.Description + "\n")
				sh.Write("    Usage: " + cmd.Usage + "\n")
			}
			return OK
		},
		Usage: "help",
	})

	sh.RegisterCommand(Command{
		Name:        "clear",
		Description: "Clear the screen",
		Handler: func(args []string) Status {
			switch runtime.GOOS {
			case "windows":
				cmd := exec.Command("cmd", "/c", "cls")
				cmd.Stdout = os.Stdout
				cmd.Run()
			default:
				cmd := exec.Command("clear")
				cmd.Stdout = os.Stdout
				cmd.Run()
			}

			return OK
		},
		Usage: "clear",
	})

	return sh
}

func (s *shell) RegisterCommand(cmd Command) {
	s.commands[cmd.Name] = cmd
}

func (s *shell) RegisterEarlyExecCommand(cmd EarlyCommand) {
	s.earlyExecCommands = append(s.earlyExecCommands, cmd)
}

func (s *shell) executeEarlyCommands() {
	for _, cmd := range s.earlyExecCommands {
		cmd.Handler()
	}
}

func (s *shell) executeCommand(cmd string, args []string) Status {
	if strings.ToUpper(cmd) == string(EXIT) {
		return EXIT
	}

	if cmd == "" {
		return OK
	}

	if command, ok := s.commands[cmd]; ok {
		return command.Handler(args)
	}

	return NOT_FOUND
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

func (s *shell) read() string {
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

func (s *shell) Write(output string) {
	s.outStream.Write([]byte(output))
}

func (s *shell) Run(welcMessege string) {
	var stat Status
	s.executeCommand("clear", nil)
	s.Write(welcMessege)
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
		} else if stat == NOT_FOUND {
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

func (s *shell) Exit() {
	os.Exit(0)
}
