package engine

import "strings"

func NormalizeToken(token string) string {
	return strings.ToLower(token)
}
