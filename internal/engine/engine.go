package engine

import (
	"Boolean-IR-System/internal"
)

type Engine struct {
	nextDocID int
	docs      []*internal.Document
	index     map[string]internal.SortedStructure
	library   map[string]struct{}
}

func NewEngine() *Engine {
	return &Engine{
		nextDocID: 0,
		docs:      make([]*internal.Document, 0),
		index:     make(map[string]internal.SortedStructure),
		library:   make(map[string]struct{}),
	}
}

func (e *Engine) GetDocuments() []*internal.Document {
	return e.docs
}

func (e *Engine) GetDocumentByID(id int) *internal.Document {
	if id >= 0 && id < len(e.docs) {
		return e.docs[id]
	}
	return nil
}

func (e *Engine) GetDocumentsSize() int {
	return e.nextDocID
}
