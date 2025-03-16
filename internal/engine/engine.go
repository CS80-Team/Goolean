package engine

import (
	"github.com/CS80-Team/Boolean-IR-System/internal"
	"github.com/CS80-Team/Boolean-IR-System/internal/structures/ordered"

	"github.com/CS80-Team/Boolean-IR-System/internal/engine/tokenizer"
	"github.com/CS80-Team/Boolean-IR-System/internal/textprocessing"
)

type Engine struct {
	docs      []*internal.Document
	index     map[string]ordered.OrderedStructure[int]
	library   map[string]struct{}
	processor textprocessing.Processor
	tokener   *tokenizer.Tokener
}

func NewEngine(processor textprocessing.Processor, tokener *tokenizer.Tokener) *Engine {
	return &Engine{
		docs:      make([]*internal.Document, 0),
		index:     make(map[string]ordered.OrderedStructure[int]),
		library:   make(map[string]struct{}),
		processor: processor,
		tokener:   tokener,
	}
}

func (e *Engine) AddDocument(doc *internal.Document) {
	if _, ok := e.library[doc.Name]; !ok {
		doc.ID = e.GetNextDocID()
		e.docs = append(e.docs, doc)
		e.library[doc.Path] = struct{}{}
		e.parseDocument(doc)
	}
}

func (e *Engine) ProcessToken(token string) string {
	return e.processor.Process(token)
}

func (e *Engine) GetDocuments() []*internal.Document {
	return e.docs
}

func (e *Engine) GetKeyIndex(key string) ordered.OrderedStructure[int] {
	if _, ok := e.index[key]; !ok {
		panic("[Engine]: Key index not found")
	}
	return e.index[key]
}

func (e *Engine) GetIndexSize() int {
	return len(e.index)
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
