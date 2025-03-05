package textprocessing

var stopWords = map[string]bool{
	"the": true,
	"a":   true,
	"to":  true,
	"if":  true,
}

type StopWordRemover struct{}

func NewStopWordRemover() *StopWordRemover {
	return &StopWordRemover{}
}

func (s *StopWordRemover) Process(token string) string {
	if stopWords[token] {
		return ""
	}
	return token
}
