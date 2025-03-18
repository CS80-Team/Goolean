package shell

import (
    "github.com/chzyer/readline"
)

type KeyListener struct {
	shell *Shell
}

func (l *KeyListener) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	if key == readline.CharTab {
		l.shell.logger.GetLogger().Info("Tab key pressed: " + string(line)[:pos-1])
		cmd, found := l.shell.autoCompleteCommand(string(line)[:pos-1])
		if found {
			return []rune(cmd), len(cmd), true
		}

		return nil, 0, false
	}

	return nil, 0, false
}
