package ordered

import (
	"golang.org/x/exp/constraints"
)

type OrderedSlice[Entry constraints.Ordered] struct {
	data []Entry
}

func NewOrderedSlice[Entry constraints.Ordered]() *OrderedSlice[Entry] {
	return &OrderedSlice[Entry]{}
}

func NewOrderedSliceWithCapacity[Entry constraints.Ordered](capacity int) OrderedStructure[Entry] {
	return &OrderedSlice[Entry]{data: make([]Entry, 0, capacity)}
}

func NewOrderedSliceWithSlice[Entry constraints.Ordered](slice []Entry) OrderedStructure[Entry] {
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
