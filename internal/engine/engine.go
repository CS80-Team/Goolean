package engine

import (
	"github.com/CS80-Team/Boolean-IR-System/internal"
	"github.com/CS80-Team/Boolean-IR-System/internal/structures/ordered"

	"github.com/CS80-Team/Boolean-IR-System/internal/engine/tokenizer"
	"github.com/CS80-Team/Boolean-IR-System/internal/textprocessing"
)

type Engine struct {
	// `docs` stores Document structs that are stored in the engine,
	// they are sorted by the order they were added to the engine and assigned an ID which is their index in the slice,
	// it is used to retrieve documents by their ID.
	docs []*internal.Document

	// `index` maps a tokens (keys) to an ordered list of document IDs that contain that token,
	// the list stores the documents sotred by their ID's.
	index map[string]ordered.OrderedStructure[int]

	// `library` is a "set" that stores documents names to avoid adding the same document multiple times.
	library map[string]struct{}

	// `processor` is used to process the tokens before adding them to the index and before querying the index,
	// it removes stop words and apply stemming and normalization to the tokens.
	processor textprocessing.Processor

	// `delimiterManager` defines how the engine should tokenize the documents.
	delimiterManager *tokenizer.DelimiterManager
}

func NewEngine(processor textprocessing.Processor, tokener *tokenizer.DelimiterManager) *Engine {
	return &Engine{
		docs:             make([]*internal.Document, 0),
		index:            make(map[string]ordered.OrderedStructure[int]),
		library:          make(map[string]struct{}),
		processor:        processor,
		delimiterManager: tokener,
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

func (e *Engine) GetDocumentsCopy() []*internal.Document {
	docs := make([]*internal.Document, len(e.docs))
	copy(docs, e.docs)
	return docs
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
