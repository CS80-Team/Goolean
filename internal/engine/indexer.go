package engine

import (
	"Boolean-IR-System/internal"
	"bufio"
	"fmt"
	"os"
)

var (
	nextDocID int = 0
	docs      []*internal.Document
	index     map[string]internal.SortedStructure = make(map[string]internal.SortedStructure)
	library   map[string]struct{}                 = make(map[string]struct{})
)

func ReadDir(path string) {
	var newDocs = Load(path)
	for _, doc := range newDocs {
		if _, ok := library[doc.Name]; !ok {
			docs = append(docs, doc)
			library[doc.Path] = struct{}{}
			doc.ID = nextDocID
			nextDocID++
			parseDocument(doc)
		}
	}
}

func parseDocument(doc *internal.Document) {
	file, err := os.Open(doc.Path + "/" + doc.Name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Indexer]: Error opening file: %s\n", err)
		return
	}

	defer file.Close()

	scan := bufio.NewScanner(file)
	for scan.Scan() {
		var line = scan.Text()
		var idx = 0
		var token string
		for idx < len(line) {
			token = getNextToken(&line, &idx)

			// token = internal.Normalize(token)
			// fmt.Println(token)

			if _, ok := index[token]; !ok {
				index[token] = &internal.SortedSlice{}
			}

			index[token].InsertSorted(doc.ID)
		}
	}
}

func getNextToken(line *string, idx *int) string {
	var token string
	for *idx < len(*line) && isDelimeter(rune((*line)[*idx])) {
		*idx++
	}

	for *idx < len(*line) && !isDelimeter(rune((*line)[*idx])) {
		token += string((*line)[*idx])
		*idx++
	}

	for *idx < len(*line) && isDelimeter(rune((*line)[*idx])) {
		*idx++
	}

	return token
}

func isDelimeter(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == ',' || c == '.' || c == '!' || c == '?' || c == ':' || c == ';'
}

func Query(key1, key2, op string) {
	if op == "AND" {
		var res = intersection(index[key1], index[key2])
		for i := 0; i < res.GetLength(); i++ {
			fmt.Println(docs[res.At(i)].Name)
		}
	} else if op == "OR" {
		var res = union(index[key1], index[key2])
		for i := 0; i < res.GetLength(); i++ {
			fmt.Println(docs[res.At(i)].Name)
		}
	} else {
		var res = inverse(index[key1])
		for i := 0; i < res.GetLength(); i++ {
			fmt.Println(docs[res.At(i)].Name)
		}
	}
	fmt.Println()
}

func inverse(s internal.SortedStructure) internal.SortedStructure {
	var res = internal.NewSortedSlice()

	for i := 0; i < nextDocID; i++ {
		if _, found := s.BinarySearch(i); !found {
			res.InsertSorted(i)
		}
	}

	return res
}

func intersection(s1, s2 internal.SortedStructure) internal.SortedStructure {
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

func union(s1, s2 internal.SortedStructure) internal.SortedStructure {
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
