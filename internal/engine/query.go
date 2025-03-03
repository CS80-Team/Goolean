package engine

import (
	"Boolean-IR-System/internal"
	"fmt"
)

func (e *Engine) Query(key1, key2, op string) {
	var result internal.SortedStructure

	switch op {
	case "AND":
		result = e.intersection(e.index[key1], e.index[key2])
	case "OR":
		result = e.union(e.index[key1], e.index[key2])
	case "NOT":
		result = e.inverse(e.index[key1])
	default:
		fmt.Println("Invalid operation")
		return
	}

	for i := 0; i < result.GetLength(); i++ {
		fmt.Println(e.docs[result.At(i)].Name)
	}
}

func (e *Engine) inverse(s internal.SortedStructure) internal.SortedStructure {
	var res = internal.NewSortedSlice()

	if s == nil {
		return res
	}

	for i := 0; i < e.nextDocID; i++ {
		if _, found := s.BinarySearch(i); !found {
			res.InsertSorted(i)
		}
	}

	return res
}

func (e *Engine) intersection(s1, s2 internal.SortedStructure) internal.SortedStructure {
	var res = internal.NewSortedSlice()

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
