package ordered

import (
	"golang.org/x/exp/constraints"
)

type OrderedSlice[T constraints.Ordered] struct {
	data []T
}

func NewSortedSlice[T constraints.Ordered](initialData ...T) *OrderedSlice[T] {
	if len(initialData) > 0 {
		return &OrderedSlice[T]{data: initialData}
	}
	return &OrderedSlice[T]{data: make([]T, 0)}
}

func (s *OrderedSlice[T]) BinarySearch(val T) int {
	low, high := 0, len(s.data)-1
	for low <= high {
		mid := low + ((high - low) >> 1)
		if s.data[mid] == val {
			return mid
		}
		if s.data[mid] < val {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return -1
}

func (s *OrderedSlice[Entry]) GetLength() int {
	return len(s.data)
}

func (s *OrderedSlice[Entry]) At(idx int) Entry {
	if idx < 0 || idx >= len(s.data) {
		panic("Index out of range")
	}

	return s.data[idx]
}

func (s *OrderedSlice[Entry]) LowerBound(val Entry) int {
	var low, high = 0, len(s.data) - 1
	var mid int
	for low < high {
		mid = low + ((high - low) >> 1)
		if s.data[mid] < val {
			low = mid + 1
		} else {
			high = mid
		}
	}

	return low
}

func (s *OrderedSlice[Entry]) UpperBound(val Entry) int {
	var low, high = 0, len(s.data) - 1
	var mid int
	for low <= high {
		mid = low + ((high - low) >> 1)
		if s.data[mid] <= val {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return low
}

func (s *OrderedSlice[Entry]) InsertSorted(val Entry) {
	var idx = s.UpperBound(val)

	if idx-1 >= 0 && s.data[idx-1] == val { // neglect duplicates
		return
	}

	s.data = append(s.data, val)
	var swp = s.data[idx]
	s.data[idx] = val
	for i := idx + 1; i < len(s.data); i++ {
		s.data[i], swp = swp, s.data[i]
	}
}
