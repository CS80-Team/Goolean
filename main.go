package main

import (
	"os"
	"path/filepath"

	"github.com/CS80-Team/Boolean-IR-System/internal/engine"
	"github.com/CS80-Team/Boolean-IR-System/internal/engine/tokenizer"
	"github.com/CS80-Team/Boolean-IR-System/internal/textprocessing"
	"github.com/CS80-Team/Boolean-IR-System/shell"
)

func main() {
	engine := engine.NewEngine(
		textprocessing.NewDefaultProcessor(
			textprocessing.NewNormalizer(),
			textprocessing.NewStemmer(),
			textprocessing.NewStopWordRemover(),
		),
		tokenizer.NewTokener(
			&map[rune]struct{}{
				' ': {},

				',':  {},
				'?':  {},
				'!':  {},
				'.':  {},
				';':  {},
				':':  {},
				'\\': {},

				'(': {},
				')': {},
				'[': {},
				']': {},
				'{': {},
				'}': {},

				'=': {},
				'+': {},
				'-': {},
				'*': {},
				'/': {},
				'%': {},
				'^': {},
			},
		),
	)

	engine.LoadDirectory(filepath.Join(filepath.Base("."), "dataset"))

	s := shell.NewShell(os.Stdin, os.Stdout)

	RegisterCommands(s, engine)

	s.Run("Welcome to the Boolean IR System shell, type `help` for list of commands\n")
}
