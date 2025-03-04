package structures

type Stack[Entry any] struct {
	data []Entry
	idx  int
}

func NewStack[Entry any](params ...interface{}) *Stack[Entry] {
	if len(params) == 0 {
		return &Stack[Entry]{data: make([]Entry, 10), idx: -1}
	}
	switch v := params[0].(type) {
	case int:
		return &Stack[Entry]{data: make([]Entry, v), idx: -1}
	case []Entry:
		var newData = make([]Entry, len(v))
		copy(newData, v)
		return &Stack[Entry]{data: newData, idx: len(v) - 1}
	default:
		panic("Invalid parameter type")
	}
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
