package main

import (
	"Boolean-IR-System/internal/engine"
	"Boolean-IR-System/internal/textprocessing"
	"Boolean-IR-System/shell"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	engine := engine.NewEngine(textprocessing.NewDefaultProcessor(
		textprocessing.NewNormalizer(),
		textprocessing.NewStemmer(),
		textprocessing.NewStopWordRemover(),
	))

	engine.LoadDirectory(os.Getenv("DATASET_PATH"))

	s := shell.NewShell(os.Stdin, os.Stdout)

	RegisterCommands(s, engine)

	s.Run("Welcome to the Boolean IR System shell, type `help` for list of commands\n")
}
