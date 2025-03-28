package ordered

import (
	"math"

	"golang.org/x/exp/constraints"
)

type SkipNode[Entry constraints.Integer] struct {
	entry Entry
	next  *SkipNode[Entry]
	skip  *SkipNode[Entry]
}

var _ OrderedStructure[int] = &SkipPointerList[int]{}

type SkipPointerList[Entry constraints.Integer] struct {
	head *SkipNode[Entry]
	size int
}

func (s *SkipPointerList[Entry]) New() *SkipPointerList[Entry] {
	return NewSkipPointerList[Entry]()
}

func NewSkipPointerList[Entry constraints.Integer]() *SkipPointerList[Entry] {
	return &SkipPointerList[Entry]{}
}

func NewSkipPointerListWithCapacity[Entry constraints.Integer](capacity int) *SkipPointerList[Entry] {
	var zeroValue Entry
	s := NewSkipPointerList[Entry]()
	for i := 0; i < capacity; i++ {
		s.InsertSorted(zeroValue)
	}

	return s
}

func NewSkipPointerListWithSlice[Entry constraints.Integer](slice []Entry) *SkipPointerList[Entry] {
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

func (s *SkipPointerList[Entry]) Complement(maxDocId int) OrderedStructure[Entry] {
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

func (s1 *SkipPointerList[Entry]) Intersection(s2 OrderedStructure[Entry]) OrderedStructure[Entry] {
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

func (s1 *SkipPointerList[Entry]) Union(s2 OrderedStructure[Entry]) OrderedStructure[Entry] {
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
