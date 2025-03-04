package engine

import (
	"Boolean-IR-System/internal/structures"
	"fmt"
	"strings"
)

const (
	AND = "AND"
	OR  = "OR"
	NOT = "NOT"
)

func (e *Engine) Query(query string) structures.OrderedStructure[int] {
	var ops = structures.NewStack[string]()
	var keys = structures.NewStack[string]()
	var res structures.OrderedStructure[int]
	fmt.Println("Constructing query:", query)
	tokens := strings.Fields(query)

	if len(tokens) == 0 {
		return nil
	} else if len(tokens) == 1 {
		if tokens[0] == NOT || tokens[0] == AND || tokens[0] == OR {
			panic("[Engine]: Invalid query, missing operand")
		}
		return e.index[tokens[0]]
	}

	for _, token := range tokens {
		if token == AND || token == OR || token == NOT {
			ops.Push(token)
		} else {
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

	return res
}

func (e *Engine) inverse(s structures.OrderedStructure[int]) structures.OrderedStructure[int] {
	var res = structures.NewSortedSlice[int]()
	if s == nil {
		return res
	}

	for i := 0; i < e.GetDocumentsSize(); i++ {
		if idx := s.BinarySearch(i); idx != -1 {
			res.InsertSorted(i)
		}
	}

	return res
}

func (e *Engine) intersection(s1, s2 structures.OrderedStructure[int]) structures.OrderedStructure[int] {
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
	var res = structures.NewSortedSlice[int]()

	var i, j = 0, 0
	for i < s1.GetLength() && j < s2.GetLength() {
		if s1.At(i) == s2.At(j) {
			res.InsertSorted(s1.At(i))
			i++
			j++
		} else {
			for i < s1.GetLength() && s1.At(i) < s2.At(j) {
				res.InsertSorted(s1.At(i))
				i++
			}

			if i < s1.GetLength() && s1.At(i) > s2.At(j) {
				for j < s2.GetLength() && s2.At(j) < s1.At(i) {
					res.InsertSorted(s2.At(j))
					j++
				}
			}
		}
	}

	for i < s1.GetLength() {
		res.InsertSorted(s1.At(i))
		i++
	}
	for j < s2.GetLength() {
		res.InsertSorted(s2.At(j))
		j++
	}
	return res
}
