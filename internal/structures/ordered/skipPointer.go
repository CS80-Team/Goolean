package ordered

import "golang.org/x/exp/constraints"

type SkipNode[Entry constraints.Ordered] struct {
	entry Entry
	next  *SkipNode[Entry]
	skip  *SkipNode[Entry]
}

type SkipPointer[Entry constraints.Ordered] struct {
	head *SkipNode[Entry]
	size int
}

func NewSkipPointer[Entry constraints.Ordered]() *SkipPointer[Entry] {
	return &SkipPointer[Entry]{}
}

func NewSkipPointerWithSlice[Entry constraints.Ordered](slice []Entry) *SkipPointer[Entry] {
	newS := NewSkipPointer[Entry]()
	for _, val := range slice {
		newS.InsertSorted(val)
	}

	// newS.updateSkipPointers()

	return newS
}

func (s *SkipPointer[Entry]) InsertSorted(entry Entry) {
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

func (s *SkipPointer[Entry]) GetLength() int {
	return s.size
}

func (s *SkipPointer[Entry]) IsEmpty() bool {
	return s.size == 0
}

func (s *SkipPointer[Entry]) At(index int) Entry {
	if index < 0 || index >= s.size {
		panic("[SkipPointer]: Index out of range")
	}

	curr := s.head
	for i := 0; i < index; i++ {
		curr = curr.next
	}

	return curr.entry
}
