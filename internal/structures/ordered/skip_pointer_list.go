package ordered

import (
	"math"

	"golang.org/x/exp/constraints"
)

type SkipNode[Entry constraints.Ordered] struct {
	entry Entry
	next  *SkipNode[Entry]
	skip  *SkipNode[Entry]
}


var _ OrderedStructure[int] = &SkipPointerList[int]{}

type SkipPointerList[Entry constraints.Ordered] struct {
	head *SkipNode[Entry]
	size int
}

func (s *SkipPointerList[Entry]) New() *SkipPointerList[Entry] {
	return NewSkipPointerList[Entry]()
}

func NewSkipPointerList[Entry constraints.Ordered]() *SkipPointerList[Entry] {
	return &SkipPointerList[Entry]{}
}

func NewSkipPointerListWithCapacity[Entry constraints.Ordered](capacity int) *SkipPointerList[Entry] {
	var zeroValue Entry
	s := NewSkipPointerList[Entry]()
	for i := 0; i < capacity; i++ {
		s.InsertSorted(zeroValue)
	}

	return s
}

func NewSkipPointerListWithSlice[Entry constraints.Ordered](slice []Entry) *SkipPointerList[Entry] {
	newS := NewSkipPointerList[Entry]()
	for _, val := range slice {
		newS.InsertSorted(val)
	}

	newS.UpdateSkipPointers()

	return newS
}

func (s *SkipPointerList[Entry]) UpdateSkipPointers() {
	if s.IsEmpty() {
		return
	}

	blockSize := int(math.Sqrt(float64(s.size)))
	prev := s.head
	curr := s.head.next
	for i := 1; i < s.size; i++ {
		if i%blockSize == 0 {
			prev.skip = curr
			prev = curr
		}
		curr = curr.next
	}
}

func (s *SkipPointerList[Entry]) InsertSorted(entry Entry) {
	newNode := &SkipNode[Entry]{entry: entry}
	if s.head == nil {
		s.head = newNode
		s.size++
		return
	}

	curr := s.head
	for curr.next != nil && curr.next.entry < entry {
		curr = curr.next
	}

	newNode.next = curr.next
	curr.next = newNode
	s.size++
}

func (s *SkipPointerList[Entry]) BinarySearch(Entry) int {
	return -1
}

func (s *SkipPointerList[Entry]) LowerBound(Entry) int {
	return -1
}

func (s *SkipPointerList[Entry]) UpperBound(Entry) int {
	return -1
}

func (s *SkipPointerList[Entry]) GetLength() int {
	return s.size
}

func (s *SkipPointerList[Entry]) IsEmpty() bool {
	return s.size == 0
}

func (s *SkipPointerList[Entry]) At(index int) Entry {
	if index < 0 || index >= s.size {
		panic("[SkipPointer]: Index out of range")
	}

	curr := s.head
	for i := 0; i < index; i++ {
		curr = curr.next
	}

	return curr.entry
}
