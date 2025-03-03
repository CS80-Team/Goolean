package main

import (
	"Boolean-IR-System/internal/engine"
	"fmt"
	"os"
)

func main() {
	// args := os.Args

	// if len(args) < 2 {
	// 	panic("Usage: ./Boolean-IR-System <path>")
	// }

	e := engine.NewEngine()
	e.LoadDirectory(os.Getenv("DATASET_PATH"))
	fmt.Println(os.Getenv("DATASET_PATH"))
	e.Query("omar", "ahmed", "AND")
}
