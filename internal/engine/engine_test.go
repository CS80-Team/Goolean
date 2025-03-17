package engine

import (
	"path/filepath"
	"testing"

	"github.com/CS80-Team/BooleanEngine/internal/engine/structuresFactory"
	"github.com/CS80-Team/BooleanEngine/internal/engine/tokenizer"
	"github.com/CS80-Team/BooleanEngine/internal/structures/ordered"
	"github.com/CS80-Team/BooleanEngine/internal/textprocessing"
	"github.com/stretchr/testify/assert"
)

var (
	processor = textprocessing.NewDefaultProcessor(
		textprocessing.NewNormalizer(),
		textprocessing.NewStemmer(),
		textprocessing.NewStopWordRemover(),
	)

	tokens = map[rune]struct{}{
		' ':  {},
		'\n': {},
		',':  {},
		'?':  {},
		'!':  {},
		'.':  {},
		';':  {},
	}

	delimiterManager = tokenizer.NewDelimiterManager(
		&tokens,
	)

	indexManager = NewIndexManager(structuresFactory.NewOrderedSliceFactory[int]())

	ahmed = processor.Process("ahmed")
	omar  = processor.Process("omar")

	world = []int{0, 1, 2}
)

func TestComplement(t *testing.T) {
	e := MockEngine()

	tests := []struct {
		name     string
		input    ordered.OrderedStructure[int]
		expected []int
	}{
		{
			name:     "Inverse of non-empty set",
			input:    e.indexMgr.factory.NewWithSlice([]int{0, 1}),
			expected: []int{2},
		},
		{
			name:     "Inverse of empty set",
			input:    e.indexMgr.factory.New(),
			expected: world,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := e.complement(tt.input)
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
		s1       ordered.OrderedStructure[int]
		s2       ordered.OrderedStructure[int]
		expected []int
	}{
		{
			name:     "Intersection of two sets",
			s1:       e.indexMgr.factory.NewWithSlice([]int{0, 1}),
			s2:       e.indexMgr.factory.NewWithSlice(world),
			expected: []int{0, 1},
		},
		{
			name:     "Intersection with empty set",
			s1:       e.indexMgr.factory.NewWithSlice([]int{0, 1}),
			s2:       ordered.NewOrderedSlice[int](),
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
		s1       ordered.OrderedStructure[int]
		s2       ordered.OrderedStructure[int]
		expected []int
	}{
		{
			name:     "Union of two sets",
			s1:       e.indexMgr.factory.NewWithSlice([]int{0, 1}),
			s2:       e.indexMgr.factory.NewWithSlice([]int{0, 2}),
			expected: []int{0, 1, 2},
		},
		{
			name:     "Union with empty set",
			s1:       e.indexMgr.factory.NewWithSlice([]int{0, 1}),
			s2:       e.indexMgr.factory.New(),
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

func MockEngine() *Engine {
	e := NewEngine(processor, delimiterManager, *indexManager)

	e.LoadDirectory(filepath.Join("..", "..", "test_dataset", "engine_test_dataset"))

	return e
}

func TestQuery(t *testing.T) {
	e := MockEngine()

	tests := []struct {
		name     string
		query    []string
		expected []int // document ids
	}{
		{
			name:     "Single token",
			query:    []string{ahmed},
			expected: []int{0, 1},
		},
		{
			name:     "AND operation",
			query:    []string{ahmed, AND, omar},
			expected: []int{1},
		},
		{
			name:     "OR operation",
			query:    []string{ahmed, OR, omar},
			expected: world,
		},
		{
			name:     "AND NOT operation",
			query:    []string{ahmed, AND, NOT, omar},
			expected: []int{0},
		},
		{
			name:     "Multiple NOT operations",
			query:    []string{NOT, ahmed, AND, NOT, omar},
			expected: []int{},
		},
		{
			name:     "Query with duplicate terms",
			query:    []string{ahmed, AND, ahmed},
			expected: []int{0, 1},
		},
		{
			name:     "Empty query",
			query:    []string{""},
			expected: nil,
		},
		{
			name:     "Query engine empty-set complement",
			query:    []string{NOT},
			expected: world,
		},

		// CASES TO be Handled

		// {
		// 	name:     "Query engine world complement",
		// 	query:    []string{NOT, NOT},
		// 	expected: nil,
		// },
		// {
		// 	name:     "Query engine world intersection",
		// 	query:    []string{NOT, AND, NOT},
		// 	expected: world,
		// },
		{
			name:     "Query engine world union",
			query:    []string{NOT, OR, NOT},
			expected: world,
		},
		{
			name:     "Query a non existing key",
			query:    []string{"AMOGUS"},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := e.Query(tt.query)
			if tt.expected == nil {
				assert.Nil(t, result, "Expected nil result")
				return
			}

			assert.NotNil(t, result, "Expected non-nil result")
			assert.Equal(t, len(tt.expected), result.GetLength(), "Expected result length to match")
			for i, docID := range tt.expected {
				assert.Equal(t, docID, result.At(i), "Unexpected document ID at index %d", i)
			}
		})
	}

	invalidQueriesTests := []struct {
		name  string
		query []string
	}{
		{
			name:  "Invalid query (AND)",
			query: []string{AND},
		},
		{
			name:  "Invalid query (OR)",
			query: []string{OR},
		},
	}

	for _, tt := range invalidQueriesTests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Panics(t, func() {
				e.Query(tt.query)
			}, "Expected panic")
		})
	}
}
