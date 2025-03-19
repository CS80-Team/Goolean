package shell

import (
	"io"

	"github.com/chzyer/readline"
)

type InputHandler struct {
	reader      *readline.Instance
	prompt      string
	historyFile string
	listener    *KeyListener
}

func NewInputHandler(
	prompt,
	historyFile string,
	listener *KeyListener,
	stdin io.ReadCloser,
	stdinWriter io.Writer,
	stdout io.Writer,
	stderr io.Writer,
) (*InputHandler, error) {
	config := &readline.Config{
		Prompt:            prompt,
		HistoryFile:       historyFile,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
		Listener:          listener,
		Stdin:             stdin,
		StdinWriter:       stdinWriter,
		Stdout:            stdout,
		Stderr:            stderr,
	}

	rl, err := readline.NewEx(config)
	if err != nil {
		return nil, err
	}

	return &InputHandler{
		reader:      rl,
		prompt:      prompt,
		historyFile: historyFile,
		listener:    listener,
	}, nil
}

func (ih *InputHandler) ReadLine() (string, error) {
	return ih.reader.Readline()
}

func (ih *InputHandler) Close() {
	_ = ih.reader.Close()
}
