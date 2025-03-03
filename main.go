package main

import (
	"Boolean-IR-System/internal/engine"
)

func main() {
	// rs := retrieval.LoadDocs("dataset")

	// for _, doc := range rs {
	// 	fmt.Println(doc.Path, doc.Name)
	// }

	myEngine := engine.NewEngine()

	// wip for cross-platform paths
	myEngine.LoadDirectory("C:\\Users\\jett\\Boolean-IR-System\\dataset")

	myEngine.Query("omar", "ahmed", "AND")
	myEngine.Query("omar", "ahmed", "OR")
	myEngine.Query("omar", "ahmed", "NOT")
	myEngine.Query("ahmed", "ahmed", "NOT")
	myEngine.Query("SYSCALL", "ahmed", "NOT")
}
