package tokenizer

type DelimiterManager struct {
	delimiters map[rune]struct{}
}

func NewDelimiterManager(delimiters *map[rune]struct{}) *DelimiterManager {
	if delimiters == nil || len(*delimiters) == 0 {
		panic("[Tokenizer]: Delimiters cannot be nil or empty")
	}

	return &DelimiterManager{delimiters: *delimiters}
}

func (t *DelimiterManager) IsDelimiter(c byte) bool {
	_, ok := t.delimiters[rune(c)]
	return ok
}

func (t *DelimiterManager) AddDelimiter(c rune) {
	t.delimiters[c] = struct{}{}
}
