package main

import (
	"Boolean-IR-System/internal/engine"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	e := engine.NewEngine()
	e.LoadDirectory(os.Getenv("DATASET_PATH"))

	ret := e.Query("ahmed")
	for i := 0; i < ret.GetLength(); i++ {
		fmt.Println(e.GetDocumentByID(ret.At(i)).Name)
	}
}
