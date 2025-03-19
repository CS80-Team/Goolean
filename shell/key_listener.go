package shell

import (
	"github.com/chzyer/readline"
	"strings"
)

type KeyListener struct {
	shell *Shell
}

func (l *KeyListener) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	if key == readline.CharTab {
		input := string(line)[:pos-1]
		l.shell.logger.GetLogger().Info("Tab key pressed: " + input)

		parts := strings.Fields(input)
		if len(parts) == 0 {
			return nil, 0, false
		}

		cmd := parts[0]
		argPrefix := ""

		if len(parts) > 1 {
			argPrefix = parts[len(parts)-1]
		}

		var completion string
		var found bool
		if argPrefix == "" {
			completion, found = l.shell.autoCompleteCommand(cmd)
		} else {
			completion, found = l.shell.autoCompleteArg(cmd, argPrefix)
			completion = strings.Join(parts[:len(parts)-1], " ") + " " + completion
		}

		if found {
			return []rune(completion), len(completion), true
		}

		return nil, 0, false
	}

	return nil, 0, false
}
