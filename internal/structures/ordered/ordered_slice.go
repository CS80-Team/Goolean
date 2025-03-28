package ordered

import (
	"golang.org/x/exp/constraints"
)

type OrderedSlice[Entry constraints.Integer] struct {
	data []Entry
}

func NewOrderedSlice[Entry constraints.Integer]() *OrderedSlice[Entry] {
	return &OrderedSlice[Entry]{}
}

func NewOrderedSliceWithCapacity[Entry constraints.Integer](capacity int) OrderedStructure[Entry] {
	return &OrderedSlice[Entry]{data: make([]Entry, 0, capacity)}
}

func NewOrderedSliceWithSlice[Entry constraints.Integer](slice []Entry) OrderedStructure[Entry] {
	for i := 1; i < len(slice); i++ {
		if slice[i] < slice[i-1] {
			panic("Provided slice must be sorted")
		}
	}

	newSlice := make([]Entry, len(slice))
	copy(newSlice, slice)
	return &OrderedSlice[Entry]{data: newSlice}
}

func (o *OrderedSlice[Entry]) InsertSorted(entry Entry) {
	var idx = o.UpperBound(entry)

	if idx-1 >= 0 && o.data[idx-1] == entry { // neglect duplicates
		return
	}

	o.data = append(o.data, entry)
	var swp = o.data[idx]
	o.data[idx] = entry
	for i := idx + 1; i < len(o.data); i++ {
		o.data[i], swp = swp, o.data[i]
	}
}

func (o *OrderedSlice[Entry]) BinarySearch(entry Entry) int {
	low, high := 0, len(o.data)-1
	for low <= high {
		mid := low + ((high - low) >> 1)
		if o.data[mid] == entry {
			return mid
		}
		if o.data[mid] < entry {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return -1
}

func (o *OrderedSlice[Entry]) LowerBound(entry Entry) int {
	var low, high = 0, len(o.data) - 1
	var mid int
	for low < high {
		mid = low + ((high - low) >> 1)
		if o.data[mid] < entry {
			low = mid + 1
		} else {
			high = mid
		}
	}

	return low
}

func (o *OrderedSlice[Entry]) UpperBound(entry Entry) int {
	var low, high = 0, len(o.data) - 1
	var mid int
	for low <= high {
		mid = low + ((high - low) >> 1)
		if o.data[mid] <= entry {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return low
}

func (o *OrderedSlice[Entry]) GetLength() int {
	return len(o.data)
}

func (o *OrderedSlice[Entry]) At(i int) Entry {
	if i < 0 || i >= len(o.data) {
		panic("Index out of range")
	}

	return o.data[i]
}

func (o *OrderedSlice[Entry]) IsEmpty() bool {
	return len(o.data) == 0
}

func (s *OrderedSlice[Entry]) Complement(maxDocId int) OrderedStructure[Entry] {
	var res = NewOrderedSlice[Entry]()

	var i, j = 0, 0
	for i < s.GetLength() && j <= maxDocId {
		if s.At(i) == Entry(j) {
			i++
		} else {
			res.InsertSorted(Entry(j))
		}
		j++
	}

	for j <= maxDocId {
		res.InsertSorted(Entry(j))
		j++
	}

	return res
}

func (s1 *OrderedSlice[Entry]) Intersection(s2 OrderedStructure[Entry]) OrderedStructure[Entry] {
	if s1 == nil || s2 == nil {
		return nil
	}
	var res = NewOrderedSlice[Entry]()

	i := 0
	j := 0

	for i < s1.GetLength() && j < s2.GetLength() {
		if s1.At(i) == s2.At(j) {
			res.InsertSorted(s1.At(i))
			i++
			j++
		} else if s1.At(i) < s2.At(j) {
			i++
		} else {
			j++
		}
	}

	return res
}

func (s1 *OrderedSlice[Entry]) Union(s2 OrderedStructure[Entry]) OrderedStructure[Entry] {
	if s1 == nil || s1.IsEmpty() {
		return s2
	}
	if s2 == nil {
		return s1
	}

	var res = NewOrderedSlice[Entry]()

	for i := range s1.GetLength() {
		res.InsertSorted(s1.At(i))
	}

	for i := range s2.GetLength() {
		res.InsertSorted(s2.At(i))
	}

	return res
}
