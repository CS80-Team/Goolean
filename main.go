package main

import (
	"Boolean-IR-System/retrieval"
	"bufio"
	"fmt"
	"os"
)

func main() {
	rs := retrieval.NewRetrievalSystem()

	fmt.Print("Enter the text files directory path: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	dirPath := scanner.Text()

	err := rs.LoadTXTFilesFromDir(dirPath)
	if err != nil {
		fmt.Printf("Error loading text files: %v\n", err)
		return
	}

	fmt.Printf("Loaded %d text files\n", len(rs.Documents))
}
