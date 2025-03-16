package ordered

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertion(t *testing.T) {
	t.Run("Test Insertion from large to small", func(t *testing.T) {
		s := NewSortedSlice[int]()
		for i := 5; i > 0; i-- {
			s.InsertSorted(i)
		}

		assert.Equal(t, 5, s.GetLength(), "Expected length 5")
		for i := 0; i < 5; i++ {
			assert.Equal(t, i+1, s.At(i), "Expected element %d", i+1)
		}
	})

	t.Run("Test Insertion from small to large", func(t *testing.T) {
		s := NewSortedSlice[int]()
		for i := 0; i < 5; i++ {
			s.InsertSorted(i)
		}

		assert.Equal(t, 5, s.GetLength(), "Expected length 5")
		for i := 0; i < 5; i++ {
			assert.Equal(t, i, s.At(i), "Expected element %d", i)
		}
	})

	t.Run("Test Insertion with random order", func(t *testing.T) {
		s := NewSortedSlice[int]()
		s.InsertSorted(5)
		s.InsertSorted(3)
		s.InsertSorted(1)
		s.InsertSorted(4)
		s.InsertSorted(2)
		s.InsertSorted(6)
		s.InsertSorted(8)
		s.InsertSorted(7)
		s.InsertSorted(9)

		assert.Equal(t, 9, s.GetLength(), "Expected length 9")
		for i := 0; i < s.GetLength(); i++ {
			assert.Equal(t, i+1, s.At(i), "Expected element %d", i+1)
		}
	})

	t.Run("Test inserting the minimum value", func(t *testing.T) {
		s := NewSortedSlice[int]()

		for i := 0; i < 1000; i++ {
			s.InsertSorted(i)
		}

		s.InsertSorted(-1)
		assert.Equal(t, 1001, s.GetLength(), "Expected length 1001")
		assert.Equal(t, -1, s.At(0), "Expected element -1")
		for i := 0; i < 1000; i++ {
			assert.Equal(t, i, s.At(i+1), "Expected element %d", i)
		}
	})

	t.Run("Test inserting the maximum value", func(t *testing.T) {
		s := NewSortedSlice[int]()

		for i := 0; i < 1000; i++ {
			s.InsertSorted(i)
		}

		s.InsertSorted(1000000)
		assert.Equal(t, 1001, s.GetLength(), "Expected length 1001")
		assert.Equal(t, 1000000, s.At(1000), "Expected element 1000000")
		for i := 0; i < 1000; i++ {
			assert.Equal(t, i, s.At(i), "Expected element %d", i)
		}
	})

	t.Run("Test Insertion with duplicates", func(t *testing.T) {
		s := NewSortedSlice[int]()
		s.InsertSorted(1)
		s.InsertSorted(1)
		s.InsertSorted(1000)
		s.InsertSorted(2323)
		s.InsertSorted(1132)
		s.InsertSorted(2)
		s.InsertSorted(2)
		s.InsertSorted(1000)

		assert.Equal(t, 5, s.GetLength(), "Expected length 5")
		assert.Equal(t, 1, s.At(0), "Expected element 1")
		assert.Equal(t, 2, s.At(1), "Expected element 2")
		assert.Equal(t, 1000, s.At(2), "Expected element 1000")
		assert.Equal(t, 1132, s.At(3), "Expected element 1132")
		assert.Equal(t, 2323, s.At(4), "Expected element 2323")
	})
}

func TestSearching(t *testing.T) {
	t.Run("Test BinarySearch", func(t *testing.T) {
		s := NewSortedSlice[int]()
		s.InsertSorted(1)
		s.InsertSorted(2)
		s.InsertSorted(3)
		s.InsertSorted(4)
		s.InsertSorted(5)

		assert.Equal(t, 2, s.BinarySearch(3), "Expected index 2")
		assert.Equal(t, -1, s.BinarySearch(10), "Expected index -1")
	})
}
