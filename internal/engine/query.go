package engine

import (
	"Boolean-IR-System/internal"
	"fmt"
	"strings"
)

const (
	AND = "AND"
	OR  = "OR"
	NOT = "NOT"
)

func (e *Engine) Query(query string) internal.SortedStructure {
	var ops = internal.NewStack[string]()
	var keys = internal.NewStack[string]()
	var res internal.SortedStructure

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

			if notCount%2 == 1 {
				if res == nil {
					res = e.inverse(e.index[keys.Pop()])
				} else {
					res = e.inverse(res)
				}
			}

			if !ops.IsEmpty() {
				if ops.Peek() == AND {
					if res == nil {
						if keys.GetSize() == 1 {
							panic("[Engine]: Invalid query, missing operator")
						}
						res = e.intersection(e.index[keys.Pop()], e.index[keys.Pop()])
					} else {
						if keys.IsEmpty() {
							panic("[Engine]: Invalid query, missing operator")
						}
						res = e.intersection(res, e.index[keys.Pop()])
					}
				} else {
					if res == nil {
						if keys.GetSize() == 1 {
							panic("[Engine]: Invalid query, missing operator")
						}
						res = e.union(e.index[keys.Pop()], e.index[keys.Pop()])
					} else {
						if keys.IsEmpty() {
							panic("[Engine]: Invalid query, missing operator")
						}
						res = e.union(res, e.index[keys.Pop()])
					}
				}
				ops.Pop()
			} else {
				if keys.GetSize() > 1 || (keys.GetSize() == 1 && res != nil) {
					panic("[Engine]: Invalid query, missing operator")
				}
			}

			notCount = 0
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

func (e *Engine) inverse(s internal.SortedStructure) internal.SortedStructure {
	var res = internal.NewSortedSlice()
	fmt.Println("NOT")
	if s == nil {
		return res
	}

	for i := 0; i < e.GetDocumentsSize(); i++ {
		if _, found := s.BinarySearch(i); !found {
			res.InsertSorted(i)
		}
	}

	return res
}

func (e *Engine) intersection(s1, s2 internal.SortedStructure) internal.SortedStructure {
	var res = internal.NewSortedSlice()
	fmt.Println("AND")

	if s1.GetLength() > s2.GetLength() {
		s1, s2 = s2, s1
	}

	for i := 0; i < s1.GetLength(); i++ {
		if _, found := s2.BinarySearch(s1.At(i)); found {
			res.InsertSorted(s1.At(i))
		}
	}

	return res
}

func (e *Engine) union(s1, s2 internal.SortedStructure) internal.SortedStructure {
	var res = internal.NewSortedSlice()

	fmt.Println("OR")

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
