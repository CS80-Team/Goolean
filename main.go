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
}
