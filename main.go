package main

import (
	"Boolean-IR-System/commands"
	"Boolean-IR-System/internal/engine"
	"Boolean-IR-System/internal/textprocessing"
	"Boolean-IR-System/shell"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	txtProcessor := textprocessing.NewDefaultProcessor(
		textprocessing.NewNormalizer(),
		textprocessing.NewStemmer(),
		textprocessing.NewStopWordRemover(),
	)

	engine := engine.NewEngine(txtProcessor)
	engine.LoadDirectory(os.Getenv("DATASET_PATH"))

	s := shell.NewShell()

	s.SetInputStream(os.Stdin)
	s.SetOutputStream(os.Stdout)

	commands.RegisterCommands(s, engine)

	s.Run("Welcome to the Boolean IR System shell, type `help` for list of commands\n")
}
