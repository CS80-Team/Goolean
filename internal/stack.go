package internal

type Stack[T any] struct {
	data []T
	idx  int
}

func NewStack[T any](params ...interface{}) *Stack[T] {
	if len(params) == 0 {
		return &Stack[T]{data: make([]T, 10), idx: -1}
	}
	switch v := params[0].(type) {
	case int:
		return &Stack[T]{data: make([]T, v), idx: -1}
	case []T:
		var newData = make([]T, len(v))
		copy(newData, v)
		return &Stack[T]{data: newData, idx: len(v) - 1}
	default:
		panic("Invalid parameter type")
	}
}

func (s *Stack[T]) Push(val T) {
	s.idx++
	if s.idx == len(s.data) {
		s.data = append(s.data, val)
	} else {
		s.data[s.idx] = val
	}
}

func (s *Stack[T]) Pop() T {
	if s.IsEmpty() {
		panic("Stack is empty")
	}
	s.idx--
	return s.data[s.idx+1]
}

func (s *Stack[T]) Peek() T {
	if s.IsEmpty() {
		panic("Stack is empty")
	}
	return s.data[s.idx]
}

func (s *Stack[T]) IsEmpty() bool {
	return s.idx == -1
}

func (s *Stack[T]) GetSize() int {
	return s.idx + 1
}

func (s *Stack[T]) GetCapacity() int {
	return len(s.data)
}

func (s *Stack[T]) Clear() {
	s.idx = -1
}
