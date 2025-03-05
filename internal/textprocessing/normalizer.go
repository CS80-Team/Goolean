package textprocessing

import "strings"

type Normalizer struct{}

func NewNormalizer() *Normalizer {
	return &Normalizer{}
}

func (n *Normalizer) Process(token string) string {
	return strings.ToLower(token)
}
