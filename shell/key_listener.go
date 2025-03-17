package shell

import (
	"github.com/chzyer/readline"
)

type KeyListener struct {
	shell *Shell
}

func (l *KeyListener) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	// ctrl + l to clear the screen
	if key == readline.CharCtrlL {
		l.shell.clearScreen()
		return nil, 0, false
	}
	return nil, 0, false
}
