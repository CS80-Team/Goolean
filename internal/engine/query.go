package engine

import (
	"fmt"
	"strings"

	"github.com/CS80-Team/Goolean/internal/structures"
	"github.com/CS80-Team/Goolean/internal/structures/ordered"
)

const (
	AND = "AND"
	OR  = "OR"
	NOT = "NOT"
)

func (e *Engine) Query(tokens []string) (ordered.OrderedStructure[int], error) {
	var ops = structures.NewStack[string]()
	var keys = structures.NewStack[string]()
	var res ordered.OrderedStructure[int]

	if len(tokens) == 0 {
		return nil, fmt.Errorf("no query provided")
	} else if len(tokens) == 1 {
		// not world
		if tokens[0] == NOT {
			return e.complement(nil), nil
		}
		if tokens[0] == AND || tokens[0] == OR {
			logger.Error(QueryPrefix, "invalid query, missing operator")
			return nil, fmt.Errorf("invalid query, missing operator")
		}

		tokenized := e.ProcessToken(tokens[0])
		res = e.indexMgr.Get(tokenized)
		if res == nil {
			return nil, fmt.Errorf("no results found")
		}
		return res, nil
	}

	// HANDLE NOT AND NOT

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
							logger.Error(QueryPrefix, "invalid query, missing operator")
							return nil, fmt.Errorf("invalid query, missing operator")
						}

						if (notCount % 2) == 1 {
							res = e.intersection(e.complement(e.indexMgr.Get(keys.Pop())), e.indexMgr.Get(keys.Pop()))
						} else {
							res = e.intersection(e.indexMgr.Get(keys.Pop()), e.indexMgr.Get(keys.Pop()))
						}
					} else {
						if keys.IsEmpty() {
							logger.Error(QueryPrefix, "invalid query, missing operator")
							return nil, fmt.Errorf("invalid query, missing operator")
						}

						if (notCount % 2) == 1 {
							res = e.intersection(res, e.complement(e.indexMgr.Get(keys.Pop())))
						} else {
							res = e.intersection(res, e.indexMgr.Get(keys.Pop()))
						}
					}
				} else {
					if res == nil {
						if keys.GetSize() < 1 {
							logger.Error(QueryPrefix, "invalid query, missing operator")
							return nil, fmt.Errorf("invalid query, missing operator")
						}

						if (notCount % 2) == 1 {
							res = e.union(e.complement(e.indexMgr.Get(keys.Pop())), e.indexMgr.Get(keys.Pop()))
						} else {
							res = e.union(e.indexMgr.Get(keys.Pop()), e.indexMgr.Get(keys.Pop()))
						}
					} else {
						if keys.IsEmpty() {
							logger.Error(QueryPrefix, "invalid query, missing operator")
							return nil, fmt.Errorf("invalid query, missing operator")
						}

						if (notCount % 2) == 1 {
							res = e.union(res, e.complement(e.indexMgr.Get(keys.Pop())))
						} else {
							res = e.union(res, e.indexMgr.Get(keys.Pop()))
						}
					}
				}
				ops.Pop()
				notCount = 0
			} else {
				if keys.GetSize() > 1 || (keys.GetSize() == 1 && res != nil) {
					logger.Error(QueryPrefix, "invalid query, missing operator")
					return nil, fmt.Errorf("invalid query, missing operator")
				}
			}

			for !ops.IsEmpty() && ops.Peek() == NOT {
				notCount++
				ops.Pop()
			}

			if notCount%2 == 1 {
				if res == nil {
					res = e.complement(e.indexMgr.Get(keys.Pop()))
				} else {
					res = e.complement(res)
				}
			}
		}
	}

	if !ops.IsEmpty() && ops.Peek() == NOT {
		fmt.Println("Complement")
		ops.Pop()
		if res == nil {
			res = e.complement(nil)
		} else {
			res = e.complement(res)
		}
	}

	return res, nil
}

func (e *Engine) QueryString(query string) (ordered.OrderedStructure[int], error) {
	return e.Query(strings.Fields(query))
}

func (e *Engine) complement(s ordered.OrderedStructure[int]) ordered.OrderedStructure[int] {
	var res = e.indexMgr.factory.New()
	if s == nil {
		for i := 0; i < e.GetDocumentsSize(); i++ {
			res.InsertSorted(i)
		}
		return res
	}

	j := 0
	i := 0
	for j < s.GetLength() {
		if i == s.At(j) {
			j++
		} else {
			res.InsertSorted(i)
		}
		i++
	}

	for i < e.GetDocumentsSize() {
		res.InsertSorted(i)
		i++
	}

	return res
}

func (e *Engine) intersection(s1, s2 ordered.OrderedStructure[int]) ordered.OrderedStructure[int] {
	if s1 == nil || s2 == nil {
		return nil
	}
	var res = e.indexMgr.factory.New()

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

func (e *Engine) union(s1, s2 ordered.OrderedStructure[int]) ordered.OrderedStructure[int] {
	if s1 == nil {
		return s2
	}
	if s2 == nil {
		return s1
	}

	var res = e.indexMgr.factory.New()

	for i := range s1.GetLength() {
		res.InsertSorted(s1.At(i))
	}

	for i := range s2.GetLength() {
		res.InsertSorted(s2.At(i))
	}

	return res
}
