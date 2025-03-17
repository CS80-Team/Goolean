package engine

import (
	"bufio"
	"fmt"
	"os"

	"github.com/CS80-Team/Boolean-IR-System/internal"
	"github.com/CS80-Team/Boolean-IR-System/internal/engine/tokenizer"
)

/*
* `parseDocument` reads the document and tokenizes it using the engine's tokener
* then it processes the tokens using the engine's processor and indexes them.
 */
func (e *Engine) parseDocument(doc *internal.Document) {
	if doc == nil {
		panic("[Indexer]: Document cannot be nil")
	}

	file, err := os.Open(doc.GetFilePath())
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
		var tokenizer = tokenizer.NewTokenizer(&line, e.delimiterManager)
		var token string

		for tokenizer.HasNext() {
			token = tokenizer.NextToken()

			token = e.ProcessToken(token)

			e.indexKey(&token, doc)
		}
	}
}

func (e *Engine) indexKey(key *string, doc *internal.Document) {
	if key == nil {
		panic("[Indexer]: Key cannot be nil")
	}

	if *key == "" {
		return
	}

	e.indexMgr.Put(*key, doc.ID)
}
