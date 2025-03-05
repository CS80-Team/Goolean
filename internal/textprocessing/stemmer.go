package textprocessing

import (
	"github.com/kljensen/snowball"
)

type Stemmer struct{}

func NewStemmer() *Stemmer {
	return &Stemmer{}
}

func (s *Stemmer) Process(token string) string {
	stemmed, err := snowball.Stem(token, "english", true)
	if err == nil {
		return stemmed
	}
	return ""
}
