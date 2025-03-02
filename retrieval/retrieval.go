package retrieval

type Document struct {
	ID      int
	Name    string
	Content string
}

type BooleanRetrievalSystem struct {
	Documents      []Document
	NextDocumentID int
}

func NewRetrievalSystem() *BooleanRetrievalSystem {
	return &BooleanRetrievalSystem{
		Documents:      []Document{},
		NextDocumentID: 0,
	}
}
