package tokenizer

type Tokener struct {
	delimiters map[rune]struct{}
}

func NewTokener(delimiters *map[rune]struct{}) *Tokener {
	if delimiters == nil || len(*delimiters) == 0 {
		panic("[Tokenizer]: Delimiters cannot be nil or empty")
	}

	return &Tokener{delimiters: *delimiters}
}

func (t *Tokener) IsDelimiter(c byte) bool {
	_, ok := t.delimiters[rune(c)]
	return ok
}

func (t *Tokener) AddDelimiter(c rune) {
	t.delimiters[c] = struct{}{}
}
