package tokenizer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenization(t *testing.T) {
	var (
		tokens = map[rune]struct{}{
			' ': {},
			',': {},
			'?': {},
			'!': {},
			'.': {},
			';': {},
		}
	)

	t.Run("Test Tokenizing strings", func(t *testing.T) {
		testcases := []struct {
			input    string
			expected []string
		}{
			{
				input:    "          Hello, World!",
				expected: []string{"Hello", "World"},
			},
			{
				input:    "Hello, World! 123",
				expected: []string{"Hello", "World", "123"},
			},
			{
				input:    "Hello, World! 123 456",
				expected: []string{"Hello", "World", "123", "456"},
			},
			{
				input:    "read,write,query,delete..page2,in,book",
				expected: []string{"read", "write", "query", "delete", "page2", "in", "book"},
			},
		}

		for _, tc := range testcases {
			tokenizer := NewTokenizer(&tc.input, NewDelimiterManager(&tokens))
			i := 0
			for tokenizer.HasNext() {
				assert.NotEqual(t, i, len(tc.expected), fmt.Sprintf("Expected number of tokens to match, got %d", i))
				token := tokenizer.NextToken()
				assert.Equal(t, tc.expected[i], token, "Expected token to match")
				i++
			}

			assert.Equal(t, len(tc.expected), i, "Expected number of tokens to match")
		}
	})

	t.Run("Test Tokenizing chars", func(t *testing.T) {
		testcases := []struct {
			input    string
			expected []string
		}{
			{
				input:    "a b c d e f g...h,.e,w,.t,q",
				expected: []string{"a", "b", "c", "d", "e", "f", "g", "h", "e", "w", "t", "q"},
			},
		}

		for _, tc := range testcases {
			tokenizer := NewTokenizer(&tc.input, NewDelimiterManager(&tokens))
			i := 0
			for tokenizer.HasNext() {
				assert.NotEqual(t, i, len(tc.expected), fmt.Sprintf("Expected number of tokens to match, got %d", i))
				token := tokenizer.NextToken()
				assert.Equal(t, tc.expected[i], token, "Expected token to match")
				i++
			}

			assert.Equal(t, len(tc.expected), i, "Expected number of tokens to match")
		}
	})
}
