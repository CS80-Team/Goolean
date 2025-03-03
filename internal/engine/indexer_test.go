package engine

import (
	// "Boolean-IR-System/internal"
	"Boolean-IR-System/internal"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"
	"github.com/joho/godotenv"
)

const (
	BasePath = "../../"
)

func getAllTokens(text string) []string {
	var idx int
	var tokens []string

	for idx < len(text) {
		token := getNextToken(&text, &idx)
		tokens = append(tokens, token)
	}

	return tokens
}

func TestTokenization(t *testing.T) {
	t.Run("Test Tokenizing strings", func(t *testing.T) {
		testcases := []struct {
			input    string
			expected []string
		}{
			{
				input:    "          Hello, World!",
				expected: []string{"hello", "world"},
			},
			{
				input:    "Hello, World! 123",
				expected: []string{"hello", "world", "123"},
			},
			{
				input:    "Hello, World! 123 456",
				expected: []string{"hello", "world", "123", "456"},
			},
			{
				input:    "...A,,     B! C,,,??,,..  ,",
				expected: []string{"a", "b", "c"},
			},
			{
				input:    "read,write,query,delete..page2,in,a,book",
				expected: []string{"read", "write", "query", "delete", "page2", "in", "a", "book"},
			},
		}

		for _, tc := range testcases {
			tokens := getAllTokens(tc.input)
			assert.Equal(t, tc.expected, tokens, "Expected tokens to match")
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
			{
				input:    "Hello, World! 123",
				expected: []string{"hello", "world", "123"},
			},
		}

		for _, tc := range testcases {
			tokens := getAllTokens(tc.input)
			assert.Equal(t, tc.expected, tokens, "Expected tokens to match")
		}
	})
}

func TestIndexing(t *testing.T) {
	t.Run("Test Indexing aplhabets", func(t *testing.T) {
		err := godotenv.Load(BasePath + ".env")
		if err != nil {
			panic("Error loading .env file")
		}

		e := NewEngine()
		e.parseDocument(&internal.Document{
			Path: BasePath + os.Getenv("TEST_DATASET_PATH"), Name: "alphabet.txt",
		})
		assert.Equal(t, 26, len(e.index), "Expected the engine to have 26 keyword")
	})

	t.Run("Test Indexing numbers", func(t *testing.T) {
		err := godotenv.Load(BasePath + ".env")
		if err != nil {
			panic("Error loading .env file")
		}

		e := NewEngine()
		e.parseDocument(&internal.Document{
			Path: BasePath + os.Getenv("TEST_DATASET_PATH"), Name: "1to100.txt",
		})
		assert.Equal(t, 100, len(e.index), "Expected the engine to have 100 keyword")
	})

	t.Run("Test Indexing duplicates removal", func(t *testing.T) {
		err := godotenv.Load(BasePath + ".env")
		if err != nil {
			panic("Error loading .env file")
		}

		e := NewEngine()
		e.parseDocument(&internal.Document{
			Path: BasePath + os.Getenv("TEST_DATASET_PATH"), Name: "duplicates.txt",
		})
		assert.Equal(t, 10, len(e.index), "Expected the engine to have 10 keyword")
	})

	t.Run("Test Indexing lowercase normalization", func(t *testing.T) {
		err := godotenv.Load(BasePath + ".env")
		if err != nil {
			panic("Error loading .env file")
		}

		e := NewEngine()
		e.parseDocument(&internal.Document{
			Path: BasePath + os.Getenv("TEST_DATASET_PATH"), Name: "lowercaseNormalization.txt",
		})
		assert.Equal(t, 10, len(e.index), "Expected the engine to have 10 keyword")
	})
}
