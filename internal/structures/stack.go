package structures

type Stack[Entry any] struct {
	data []Entry
	idx  int
}

func NewStack[Entry any]() *Stack[Entry] {
	return &Stack[Entry]{data: make([]Entry, 10), idx: -1}
}

func NewStackWithCapacity[Entry any](capacity int) *Stack[Entry] {
	return &Stack[Entry]{data: make([]Entry, capacity), idx: -1}
}

func NewStackWithSlice[Entry any](slice []Entry) *Stack[Entry] {
	newS := make([]Entry, len(slice))
	copy(newS, slice)
	return &Stack[Entry]{data: newS, idx: len(slice) - 1}
}

func (s *Stack[Entry]) Push(val Entry) {
	s.idx++
	if s.idx == len(s.data) {
		s.data = append(s.data, val)
	} else {
		s.data[s.idx] = val
	}
}

func (s *Stack[Entry]) Pop() Entry {
	if s.IsEmpty() {
		panic("Stack is empty")
	}
	s.idx--
	return s.data[s.idx+1]
}

func (s *Stack[Entry]) Peek() Entry {
	if s.IsEmpty() {
		panic("Stack is empty")
	}
	return s.data[s.idx]
}

func (s *Stack[Entry]) IsEmpty() bool {
	return s.idx == -1
}

func (s *Stack[Entry]) GetSize() int {
	return s.idx + 1
}

func (s *Stack[Entry]) GetCapacity() int {
	return len(s.data)
}

func (s *Stack[Entry]) Clear() {
	s.idx = -1
}
