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
	e.LoadDirectory(os.Getenv("TEST_DATASET_PATH"))
	fmt.Println(os.Getenv("TEST_DATASET_PATH"))
	// e.Query("omar", "ahmed", "AND")
}
