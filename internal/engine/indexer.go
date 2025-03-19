package engine

import (
	"bufio"
	"fmt"
	"os"

	"github.com/CS80-Team/Goolean/internal"
	"github.com/CS80-Team/Goolean/internal/engine/tokenizer"
)

/*
* `parseDocument` reads the document and tokenizes it using the engine's delimiterManager
* then it processes the tokens using the engine's processor and indexes them.
 */
func (e *Engine) parseDocument(doc *internal.Document) {
	if doc == nil {
		logger.Error(IndexManagerPrefix, "Document cannot be nil")
		panic(IndexManagerPrefix + "Document cannot be nil")
	}

	file, err := os.Open(doc.GetFilePath())
	if err != nil {
		logger.Error(IndexManagerPrefix, fmt.Sprintf("Error opening file %s", doc.GetFilePath()))
		return
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	scan := bufio.NewScanner(file)
	for scan.Scan() {
		var line = scan.Text()
		var newTokenizer = tokenizer.NewTokenizer(&line, e.delimiterManager)
		var token string

		for newTokenizer.HasNext() {
			token = newTokenizer.NextToken()

			token = e.ProcessToken(token)

			e.indexKey(&token, doc)
		}
	}
}

func (e *Engine) indexKey(key *string, doc *internal.Document) {
	if key == nil {
		logger.Error(IndexManagerPrefix, "Key cannot be nil")
		panic(IndexManagerPrefix + "Key cannot be nil")
	}

	if *key == "" {
		return
	}

	e.indexMgr.Put(*key, doc.ID)
}
