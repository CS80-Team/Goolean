package engine

import (
	"Boolean-IR-System/internal"
	"bufio"
	"fmt"
	"os"
)

var (
	nextDocID int = 1
	// docs []internal.Document
	index   map[string][]int    = make(map[string][]int)
	library map[string]struct{} = make(map[string]struct{})
)

func ReadDir(path string) {
	var newDocs = Load(path)
	for _, doc := range newDocs {
		if _, ok := library[doc.Path]; !ok {
			// docs = append(docs, fPath)
			library[doc.Path] = struct{}{}
			doc.ID = nextDocID
			nextDocID++
			parseDocument(doc)
		}
	}
}

func parseDocument(doc *internal.Document) {
	file, err := os.Open(doc.Path + "/" + doc.Name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Indexer]: Error opening file: %s\n", err)
		return
	}

	defer file.Close()

	scan := bufio.NewScanner(file)
	for scan.Scan() {
		var line = scan.Text()
		var idx = 0
		var token string
		for idx < len(line) {
			token = getNextToken(&line, &idx)

			// token = internal.Normalize(token)

			fmt.Println(token)

			if _, ok := index[token]; !ok {
				index[token] = []int{doc.ID}
			} else {
				index[token] = append(index[token], doc.ID)
			}
		}
	}
}

func getNextToken(line *string, idx *int) string {
	var token string
	for *idx < len(*line) && isDelimeter(rune((*line)[*idx])) {
		*idx++
	}

	for *idx < len(*line) && !isDelimeter(rune((*line)[*idx])) {
		token += string((*line)[*idx])
		*idx++
	}

	for *idx < len(*line) && isDelimeter(rune((*line)[*idx])) {
		*idx++
	}

	return token
}

func isDelimeter(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == ',' || c == '.' || c == '!' || c == '?' || c == ':' || c == ';'
}
