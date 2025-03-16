package tokenizer

type Tokenizer struct {
	line    *string
	idx     int
	tokener *Tokener
}

func NewTokenizer(line *string, tokener *Tokener) *Tokenizer {
	return &Tokenizer{line: line, idx: 0, tokener: tokener}
}

func (t *Tokenizer) NextToken() string {
	if !t.HasNext() {
		panic("[Tokenizer]: No more tokens")
	}

	var token string
	for t.idx < len(*t.line) && t.tokener.IsDelimiter((*t.line)[t.idx]) {
		t.idx++
	}

	for t.idx < len(*t.line) && !t.tokener.IsDelimiter((*t.line)[t.idx]) {
		token += string(rune((*t.line)[t.idx]))
		t.idx++
	}

	for t.idx < len(*t.line) && t.tokener.IsDelimiter((*t.line)[t.idx]) {
		t.idx++
	}

	return token
}

func (t *Tokenizer) HasNext() bool {
	return t.idx < len(*t.line)
}
