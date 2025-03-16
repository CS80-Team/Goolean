package main

import (
	"os"

	"github.com/CS80-Team/Boolean-IR-System/internal/engine"
	"github.com/CS80-Team/Boolean-IR-System/internal/textprocessing"
	"github.com/CS80-Team/Boolean-IR-System/shell"
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
