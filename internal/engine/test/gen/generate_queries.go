// Generates random queries based on the dataset

package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/CS80-Team/Goolean/internal"
	"github.com/CS80-Team/Goolean/internal/engine"
	engine_tests "github.com/CS80-Team/Goolean/internal/engine/test"
	"github.com/CS80-Team/Goolean/internal/engine/tokenizer"
)

func parseDocument(doc *internal.Document) []string {
	file, err := os.Open(doc.GetFilePath())
	if err != nil {
		panic(fmt.Sprintf("Error opening file %s", doc.GetFilePath()))
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	scan := bufio.NewScanner(file)

	tokensList := make([]string, 0)
	for scan.Scan() {
		var line = scan.Text()
		var lineTokenizer = tokenizer.NewTokenizer(&line, engine_tests.DelimiterManager)
		var token string
		for lineTokenizer.HasNext() {
			token = lineTokenizer.NextToken()

			if token != "" {
				tokensList = append(tokensList, token)
			}
		}
	}

	return tokensList
}

func generateQueries(num, maxQueryLen int, keyWords []string, operators []string) []string {
	var queries []string

	if maxQueryLen <= 0 {
		panic("maxQueryLen must be greater than 0")
	}

	if num <= 0 {
		panic("num must be greater than 0")
	}

	if len(keyWords) == 0 {
		panic("keyWords must not be empty")
	}

	if len(operators) == 0 {
		panic("operators must not be empty")
	}

	// rand.Seed(time.Now().UnixNano())

	for i := 0; i < num; i++ {
		queryLen := rand.Intn(maxQueryLen) + 1
		query := ""

		// if rand.Intn(2) == 0 {
		// 	query += "NOT "
		// }

		for j := 0; j < queryLen-1; j++ {
			query += keyWords[rand.Intn(len(keyWords))] + " "
			query += operators[rand.Intn(len(operators))] + " "
		}

		query += keyWords[rand.Intn(len(keyWords))]

		queries = append(queries, query)
	}

	return queries
}

func main() {
	var keyWords []string
	operators := []string{"AND", "OR", "AND NOT", "OR NOT"}

	maxQueryLen := flag.Int("max", 5, "Maximum length of a query")
	numQueries := flag.Int("num", 1000, "Number of queries to generate")
	dataPath := flag.String("dataset", filepath.Join("..", "..", "..", "dataset"), "Path to the dataset")
	flag.Parse()

	docs := engine.LoadDocuments(*dataPath)

	for _, doc := range docs {
		keyWords = append(keyWords, parseDocument(doc)...)
	}

	fmt.Println("Details:")
	fmt.Println("\tDataset path:", *dataPath)
	fmt.Printf("\tNumber of keywords: %d\n", len(keyWords))
	fmt.Printf("\tNumber of operators: %d\n", len(operators))
	fmt.Printf("\tMax query length: %d\n", *maxQueryLen)
	fmt.Printf("\tNumber of queries: %d\n", *numQueries)
	fmt.Printf("\tNumber of documents loaded: %d\n", len(docs))
	fmt.Println("Generating queries...")

	queries := generateQueries(*numQueries, *maxQueryLen, keyWords, operators)

	file, err := os.Create("queries.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, query := range queries {
		_, err := file.WriteString(query + "\n")
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Queries generated successfully")
}
