package engine

import (
	"Boolean-IR-System/internal"
	"Boolean-IR-System/internal/structures"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func (e *Engine) parseDocument(doc *internal.Document) {
	filePath := filepath.Join(doc.Path, doc.Name)

	file, err := os.Open(filePath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[Indexer]: Error opening file: %s\n", err)
		return
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	scan := bufio.NewScanner(file)
	for scan.Scan() {
		var line = scan.Text()
		var idx = 0
		var token string
		for idx < len(line) {
			token = getNextToken(&line, &idx)

			token = e.ProcessToken(token)
			if token == "" {
				continue
			}

			if _, ok := e.index[token]; !ok {
				e.index[token] = &structures.OrderedSlice[int]{}
			}

			e.index[token].InsertSorted(doc.ID)
		}
	}
}

func getNextToken(line *string, idx *int) string {
	var token string
	for *idx < len(*line) && isDelimiter(rune((*line)[*idx])) {
		*idx++
	}

	for *idx < len(*line) && !isDelimiter(rune((*line)[*idx])) {
		token += string(rune((*line)[*idx]))
		*idx++
	}

	for *idx < len(*line) && isDelimiter(rune((*line)[*idx])) {
		*idx++
	}

	return token
}

func isDelimiter(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == ',' || c == '.' || c == '!' || c == '?' || c == ':' || c == ';'
}
