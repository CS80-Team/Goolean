package internal

type SortedStructure interface {
	InsertSorted(int)
	BinarySearch(int) (int, bool)
	LowerBound(int) int
	UpperBound(int) int
	GetLength() int
	At(int) int
}

type SortedSlice struct {
	data []int
}

func NewSortedSlice() *SortedSlice {
	return &SortedSlice{data: make([]int, 0)}
}

func (s *SortedSlice) BinarySearch(val int) (int, bool) {
	var low, high = 0, len(s.data) - 1
	var mid int
	for low <= high {
		mid = low + ((high - low) >> 1)
		if s.data[mid] == val {
			return mid, true
		}

		if s.data[mid] < val {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return -1, false
}

func (s *SortedSlice) GetLength() int {
	return len(s.data)
}

func (s *SortedSlice) At(idx int) int {
	if idx < 0 || idx >= len(s.data) {
		panic("Index out of range")
	}

	return s.data[idx]
}

func (s *SortedSlice) LowerBound(val int) int {
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

func (s *SortedSlice) UpperBound(val int) int {
	var low, high = 0, len(s.data) - 1
	var mid int
	for low < high {
		mid = low + ((high - low) >> 1)
		if s.data[mid] <= val {
			low = mid + 1
		} else {
			high = mid
		}
	}

	return low
}

func (s *SortedSlice) InsertSorted(val int) {
	var idx = s.LowerBound(val)
	s.data = append(s.data, 0)
	copy(s.data[idx+1:], s.data[idx:])
	s.data[idx] = val
}
