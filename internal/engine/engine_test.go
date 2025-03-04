package engine

import (
	"Boolean-IR-System/internal"
	"Boolean-IR-System/internal/structures"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func MockEngine() *Engine {
	e := NewEngine()

	e.docs = []*internal.Document{
		{ID: 0, Name: "ahmedAndOmar.txt", Path: os.Getenv("DATASET_PATH")},
		{ID: 1, Name: "ahmed.txt", Path: os.Getenv("DATASET_PATH")},
		{ID: 2, Name: "omar.txt", Path: os.Getenv("DATASET_PATH")},
		{ID: 3, Name: "test.txt", Path: os.Getenv("DATASET_PATH")},
	}

	e.index = map[string]structures.OrderedStructure[int]{
		"ahmed": structures.NewSortedSlice(0, 1),
		"omar":  structures.NewSortedSlice(0, 2),
		"test":  structures.NewSortedSlice(3),
	}
	return e
}

func TestQuery(t *testing.T) {
	e := MockEngine()

	tests := []struct {
		name     string
		query    string
		expected []int // document ids
	}{
		{
			name:     "Single token",
			query:    "ahmed",
			expected: []int{0, 1},
		},
		{
			name:     "AND operation",
			query:    "ahmed AND omar",
			expected: []int{0},
		},
		{
			name:     "OR operation",
			query:    "ahmed OR omar",
			expected: []int{0, 1, 2},
		},
		{
			name:     "AND NOT operation",
			query:    "ahmed AND NOT omar",
			expected: []int{1},
		},
		{
			name:     "Multiple NOT operations",
			query:    "NOT ahmed AND NOT omar",
			expected: []int{3},
		},
		{
			name:     "Query with duplicate terms",
			query:    "ahmed AND ahmed",
			expected: []int{0, 1},
		},
		{
			name:     "Empty query",
			query:    "",
			expected: nil,
		},
		{
			name:     "Invalid query (missing operand)",
			query:    "AND",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Invalid query (missing operand)" {
				assert.Panics(t, func() {
					e.Query(tt.query)
				}, "Expected panic for invalid query")
				return
			}

			result := e.Query(tt.query)
			if tt.expected == nil {
				assert.Nil(t, result, "Expected nil result")
				return
			}

			assert.Equal(t, len(tt.expected), result.GetLength(), "Unexpected result length")
			for i, docID := range tt.expected {
				assert.Equal(t, docID, result.At(i), "Unexpected document ID at index %d", i)
			}
		})
	}
}
