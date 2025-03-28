// This is an implementation of a skip pointer list, which is a single linked list with a skip pointer that skips a fixed number of nodes.
// It is interduced to benchmark the performance of the search engine.

package ordered

import (
	"math"

	"golang.org/x/exp/constraints"
)

var _ OrderedStructure[int] = &SkipPointerList[int]{}
var _ SetOperations[int] = &SkipPointerList[int]{}

type SkipNode[Entry constraints.Integer] struct {
	entry Entry
	next  *SkipNode[Entry]
	skip  *SkipNode[Entry]
}

type SkipPointerList[Entry constraints.Integer] struct {
	head          *SkipNode[Entry]
	tail          *SkipNode[Entry]
	size          int
	currBlockSize int
}

func NewSkipPointerList[Entry constraints.Integer]() *SkipPointerList[Entry] {
	return &SkipPointerList[Entry]{
		head:          nil,
		tail:          nil,
		size:          0,
		currBlockSize: 0,
	}
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

func calculateBlockSize(size int) int {
	return int(math.Sqrt(float64(size)))
}

func (s *SkipPointerList[Entry]) UpdateSkipPointers() {
	if s.IsEmpty() {
		return
	}

	// remove all skip pointers to avoid memory leaks and unexpected behavior
	curr := s.head
	for curr != nil {
		curr.skip = nil
		curr = curr.next
	}

	s.currBlockSize = calculateBlockSize(s.size)

	prev := s.head
	curr = s.head.next
	for i := 1; i < s.size; i++ {
		if i%s.currBlockSize == 0 {
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
		s.tail = newNode
		s.size++
		return
	}

	if s.head.entry > entry {
		s.pushFront(entry)
		return
	}

	if s.tail.entry < entry {
		s.pushBack(entry)
		return
	}

	curr := s.head
	for curr.next != nil && curr.next.entry < entry {
		if curr.skip != nil && curr.skip.entry < entry {
			curr = curr.skip
		} else {
			curr = curr.next
		}
	}

	// neglect duplicates
	if curr.next != nil && curr.next.entry == entry {
		return
	}

	newNode.next = curr.next
	curr.next = newNode
	s.size++

	if newNode.next == nil {
		s.tail = newNode
	}

	// Update skip pointers if the size of the list is a perfect square
	if s.size >= (s.currBlockSize+1)*(s.currBlockSize+1) {
		s.UpdateSkipPointers()
	}
}

func (s *SkipPointerList[Entry]) GetLength() int {
	return s.size
}

func (s *SkipPointerList[Entry]) IsEmpty() bool {
	return s.size == 0
}

// Used internally for `SetOperations` functions,
// where it is guaranteed that the insertion is always sorted and added to the end
func (s *SkipPointerList[Entry]) pushBack(entry Entry) {
	if s.IsEmpty() {
		s.InsertSorted(entry)
		return
	}

	if s.tail.entry > entry {
		panic("[SkipPointer]: Cannot push back an entry that is less than the tail")
	}

	if s.tail.entry == entry {
		return
	}

	s.tail.next = &SkipNode[Entry]{entry: entry}
	s.tail = s.tail.next
	s.size++
}

func (s *SkipPointerList[Entry]) pushFront(entry Entry) {
	if s.IsEmpty() {
		s.InsertSorted(entry)
		return
	}

	if s.head.entry < entry {
		panic("[SkipPointer]: Cannot push front an entry that is greater than the head")
	}

	if s.head.entry == entry {
		return
	}

	newNode := &SkipNode[Entry]{entry: entry, next: s.head}
	s.head = newNode
	s.size++
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
	var res = NewSkipPointerList[Entry]()

	var i = s.head
	var j = 0
	for i != nil && j <= maxDocId {
		if i.entry == Entry(j) {
			i = i.next
		} else {
			res.pushBack(Entry(j))
		}
		j++
	}

	for j <= maxDocId {
		res.pushBack(Entry(j))
		j++
	}

	return res
}

func (s1 *SkipPointerList[Entry]) Intersection(s2 OrderedStructure[Entry]) OrderedStructure[Entry] {
	if s1 == nil || s2 == nil {
		return nil
	}
	var res = NewSkipPointerList[Entry]()

	var i = s1.head
	var j = s2.(*SkipPointerList[Entry]).head

	for i != nil && j != nil {
		if i.entry < j.entry {
			if i.skip != nil && i.skip.entry <= j.entry {
				i = i.skip
			} else {
				i = i.next
			}
		} else if i.entry > j.entry {
			if j.skip != nil && j.skip.entry <= i.entry {
				j = j.skip
			} else {
				j = j.next
			}
		} else {
			res.pushBack(i.entry)
			i = i.next
			j = j.next
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

	var res = NewSkipPointerList[Entry]()

	var i = s1.head
	var j = s2.(*SkipPointerList[Entry]).head

	for i != nil && j != nil {
		if i.entry < j.entry {
			res.pushBack(i.entry)
			i = i.next
		} else if i.entry > j.entry {
			res.pushBack(j.entry)
			j = j.next
		} else {
			res.pushBack(i.entry)
			i = i.next
			j = j.next
		}
	}

	for i != nil {
		res.pushBack(i.entry)
		i = i.next
	}

	for j != nil {
		res.pushBack(j.entry)
		j = j.next
	}

	return res
}
