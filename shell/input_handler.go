package shell

import (
	"github.com/chzyer/readline"
)

type InputHandler struct {
	reader      *readline.Instance
	prompt      string
	historyFile string
	listener    *KeyListener
}

func NewInputHandler(prompt, historyFile string, listener *KeyListener) (*InputHandler, error) {
	config := &readline.Config{
		Prompt:            prompt,
		HistoryFile:       historyFile,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
		Listener:          listener,
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
