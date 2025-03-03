package engine

import (
	"Boolean-IR-System/internal"
)

type Engine struct {
	docs    []*internal.Document
	index   map[string]internal.SortedStructure
	library map[string]struct{}
}

func NewEngine() *Engine {
	return &Engine{
		docs:    make([]*internal.Document, 0),
		index:   make(map[string]internal.SortedStructure),
		library: make(map[string]struct{}),
	}
}

func (e *Engine) GetDocuments() []*internal.Document {
	return e.docs
}

func (e *Engine) GetDocumentByID(id int) *internal.Document {
	if id < 0 && id >= len(e.docs) {
		panic("[Engine]: Document ID out of range")
	}

	return e.docs[id]
}

func (e *Engine) GetDocumentsSize() int {
	return len(e.docs)
}

func (e *Engine) GetNextDocID() int {
	return e.GetDocumentsSize()
}
