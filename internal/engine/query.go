package engine

import (
	"Boolean-IR-System/internal/structures"
	"strings"
)

const (
	AND = "AND"
	OR  = "OR"
	NOT = "NOT"
)

func (e *Engine) Query(tokens []string) structures.OrderedStructure[int] {
	var ops = structures.NewStack[string]()
	var keys = structures.NewStack[string]()
	var res structures.OrderedStructure[int]

	if len(tokens) == 0 {
		return nil
	} else if len(tokens) == 1 {
		// not world
		if tokens[0] == NOT {
			return res
		}
		if tokens[0] == AND || tokens[0] == OR {
			panic("[Engine]: Invalid query, missing operand")
		}
		tokenized := e.ProcessToken(tokens[0])
		return e.index[tokenized]
	}

	for _, token := range tokens {
		if token == AND || token == OR || token == NOT {
			ops.Push(token)
		} else {
			token = e.ProcessToken(token)

			keys.Push(token)

			var notCount = 0
			for !ops.IsEmpty() && ops.Peek() == NOT {
				notCount++
				ops.Pop()
			}

			if !ops.IsEmpty() {
				if ops.Peek() == AND {
					if res == nil {
						if keys.GetSize() < 1 {
							panic("[Engine]: Invalid query, missing operator")
						}

						if (notCount % 2) == 1 {
							res = e.intersection(e.inverse(e.index[keys.Pop()]), e.index[keys.Pop()])
						} else {
							res = e.intersection(e.index[keys.Pop()], e.index[keys.Pop()])
						}
					} else {
						if keys.IsEmpty() {
							panic("[Engine]: Invalid query, missing operator")
						}

						if (notCount % 2) == 1 {
							res = e.intersection(res, e.inverse(e.index[keys.Pop()]))
						} else {
							res = e.intersection(res, e.index[keys.Pop()])
						}
					}
				} else {
					if res == nil {
						if keys.GetSize() < 1 {
							panic("[Engine]: Invalid query, missing operator")
						}

						if (notCount % 2) == 1 {
							res = e.union(e.inverse(e.index[keys.Pop()]), e.index[keys.Pop()])
						} else {
							res = e.union(e.index[keys.Pop()], e.index[keys.Pop()])
						}
					} else {
						if keys.IsEmpty() {
							panic("[Engine]: Invalid query, missing operator")
						}

						if (notCount % 2) == 1 {
							res = e.union(res, e.inverse(e.index[keys.Pop()]))
						} else {
							res = e.union(res, e.index[keys.Pop()])
						}
					}
				}
				ops.Pop()
				notCount = 0
			} else {
				if keys.GetSize() > 1 || (keys.GetSize() == 1 && res != nil) {
					panic("[Engine]: Invalid query, missing operator")
				}
			}

			for !ops.IsEmpty() && ops.Peek() == NOT {
				notCount++
				ops.Pop()
			}

			if notCount%2 == 1 {
				if res == nil {
					res = e.inverse(e.index[keys.Pop()])
				} else {
					res = e.inverse(res)
				}
			}
		}
	}

	if !ops.IsEmpty() && ops.Peek() == NOT {
		ops.Pop()
		if res == nil {
			res = e.inverse(nil)
		} else {
			res = e.inverse(res)
		}
	}

	return res
}

func (e *Engine) QueryString(query string) structures.OrderedStructure[int] {
	return e.Query(strings.Fields(query))
}

func (e *Engine) inverse(s structures.OrderedStructure[int]) structures.OrderedStructure[int] {
	var res = structures.NewSortedSlice[int]()
	if s == nil {
		for i := 0; i < e.GetDocumentsSize(); i++ {
			res.InsertSorted(i)
		}
		return res
	}

	j := 0
	for i := 0; i < e.GetDocumentsSize() && j < s.GetLength(); i++ {
		if i == s.At(j) {
			j++
		} else {
			res.InsertSorted(i)
		}
	}

	return res
}

func (e *Engine) intersection(s1, s2 structures.OrderedStructure[int]) structures.OrderedStructure[int] {
	if s1 == nil || s2 == nil {
		return nil
	}
	var res = structures.NewSortedSlice[int]()

	if s1.GetLength() > s2.GetLength() {
		s1, s2 = s2, s1
	}

	for i := 0; i < s1.GetLength(); i++ {
		if idx := s2.BinarySearch(s1.At(i)); idx != -1 {
			res.InsertSorted(s1.At(i))
		}
	}

	return res
}

func (e *Engine) union(s1, s2 structures.OrderedStructure[int]) structures.OrderedStructure[int] {
	if s1 == nil {
		return s2
	}
	if s2 == nil {
		return s1
	}

	var res = structures.NewSortedSlice[int]()

	for i := range s1.GetLength() {
		res.InsertSorted(s1.At(i))
	}

	for i := range s2.GetLength() {
		res.InsertSorted(s2.At(i))
	}

	return res
}
