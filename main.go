package main

import (
	"Boolean-IR-System/internal/engine"
)

func main() {
	// rs := retrieval.Load("dataset")

	// for _, doc := range rs {
	// 	fmt.Println(doc.Path, doc.Name)
	// }

	engine.ReadDir("dataset")

	engine.Query("omar", "ahmed", "AND")
	engine.Query("omar", "ahmed", "OR")
	engine.Query("omar", "ahmed", "NOT")
	engine.Query("ahmed", "ahmed", "NOT")
	engine.Query("SYSCALL", "ahmed", "NOT")
}
