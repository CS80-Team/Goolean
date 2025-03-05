package engine

import (
	"Boolean-IR-System/internal"
	"Boolean-IR-System/internal/structures"
	"Boolean-IR-System/internal/textprocessing"
	"testing"

	"github.com/stretchr/testify/assert"
)

var processor = textprocessing.NewDefaultProcessor(
	textprocessing.NewNormalizer(),
	textprocessing.NewStemmer(),
	textprocessing.NewStopWordRemover(),
)

func MockEngine() *Engine {
	e := NewEngine(processor)

	e.docs = []*internal.Document{
		{ID: 0, Name: "ahmedAndOmar.txt"},
		{ID: 1, Name: "ahmed.txt"},
		{ID: 2, Name: "omar.txt"},
		{ID: 3, Name: "test.txt"},
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
					e.QueryString(tt.query)
				}, "Expected panic for invalid query")
				return
			}

			result := e.QueryString(tt.query)
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

func TestInverse(t *testing.T) {
	e := MockEngine()

	tests := []struct {
		name     string
		input    structures.OrderedStructure[int]
		expected []int
	}{
		{
			name:     "Inverse of non-empty set",
			input:    structures.NewSortedSlice[int](0, 1),
			expected: []int{2, 3},
		},
		{
			name:     "Inverse of empty set",
			input:    structures.NewSortedSlice[int](),
			expected: []int{0, 1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := e.inverse(tt.input)
			assert.Equal(t, len(tt.expected), result.GetLength(), "Unexpected result length")
			for i, docID := range tt.expected {
				assert.Equal(t, docID, result.At(i), "Unexpected document ID at index %d", i)
			}
		})
	}
}

func TestIntersection(t *testing.T) {
	e := MockEngine()

	tests := []struct {
		name     string
		s1       structures.OrderedStructure[int]
		s2       structures.OrderedStructure[int]
		expected []int
	}{
		{
			name:     "Intersection of two sets",
			s1:       structures.NewSortedSlice[int](0, 1),
			s2:       structures.NewSortedSlice[int](0, 1, 2),
			expected: []int{0, 1},
		},
		{
			name:     "Intersection with empty set",
			s1:       structures.NewSortedSlice[int](0, 1),
			s2:       structures.NewSortedSlice[int](),
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := e.intersection(tt.s1, tt.s2)
			assert.Equal(t, len(tt.expected), result.GetLength(), "Unexpected result length")
			for i, docID := range tt.expected {
				assert.Equal(t, docID, result.At(i), "Unexpected document ID at index %d", i)
			}
		})
	}
}

func TestUnion(t *testing.T) {
	e := MockEngine()

	tests := []struct {
		name     string
		s1       structures.OrderedStructure[int]
		s2       structures.OrderedStructure[int]
		expected []int
	}{
		{
			name:     "Union of two sets",
			s1:       structures.NewSortedSlice[int](0, 1),
			s2:       structures.NewSortedSlice[int](0, 2),
			expected: []int{0, 1, 2},
		},
		{
			name:     "Union with empty set",
			s1:       structures.NewSortedSlice[int](0, 1),
			s2:       structures.NewSortedSlice[int](),
			expected: []int{0, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := e.union(tt.s1, tt.s2)
			assert.Equal(t, len(tt.expected), result.GetLength(), "Unexpected result length")
			for i, docID := range tt.expected {
				assert.Equal(t, docID, result.At(i), "Unexpected document ID at index %d", i)
			}
		})
	}
}
