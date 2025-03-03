package internal

type Stack struct {
	data []int
	idx  int
}

func NewStack(params ...interface{}) *Stack {
	if len(params) == 0 {
		return &Stack{data: make([]int, 10), idx: -1}
	}
	switch v := params[0].(type) {
	case int:
		return &Stack{data: make([]int, v), idx: -1}
	case []int:
		var newData = make([]int, len(v))
		copy(newData, v)
		return &Stack{data: newData, idx: len(v) - 1}
	default:
		panic("Invalid parameter type")
	}
}

func (s *Stack) Push(val int) {
	s.idx++
	if s.idx == s.GetCapacity() {
		s.data = append(s.data, val)
	} else {
		s.data[s.idx] = val
	}
}

func (s *Stack) Pop() int {
	if s.IsEmpty() {
		panic("Stack is empty")
	}

	s.idx--
	return s.data[s.idx+1]
}

func (s *Stack) Peek() int {
	if s.IsEmpty() {
		panic("Stack is empty")
	}

	return s.data[s.idx]
}

func (s *Stack) IsEmpty() bool {
	return s.idx == -1
}

func (s *Stack) GetSize() int {
	return s.idx + 1
}

func (s *Stack) GetCapacity() int {
	return len(s.data)
}

func (s *Stack) Clear() {
	s.idx = -1
}
