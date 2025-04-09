package tokenizer

type Tokenizer struct {
	line             *string
	idx              int
	delimiterManager *DelimiterManager
}

func NewTokenizer(line *string, delimiterManager *DelimiterManager) *Tokenizer {
	tokenizer := Tokenizer{line: line, idx: 0, delimiterManager: delimiterManager}
	tokenizer.adjustIdxPointer()

	return &tokenizer
}

func (t *Tokenizer) adjustIdxPointer() {
	for t.idx < len(*t.line) && t.delimiterManager.IsDelimiter((*t.line)[t.idx]) {
		t.idx++
	}
}

func (t *Tokenizer) NextToken() string {
	if !t.HasNext() {
		panic("[Tokenizer]: No more tokens")
	}

	var token string
	t.adjustIdxPointer()

	for t.idx < len(*t.line) && !t.delimiterManager.IsDelimiter((*t.line)[t.idx]) {
		token += string(rune((*t.line)[t.idx]))
		t.idx++
	}

	t.adjustIdxPointer()

	return token
}

func (t *Tokenizer) HasNext() bool {
	return t.idx < len(*t.line)
}
