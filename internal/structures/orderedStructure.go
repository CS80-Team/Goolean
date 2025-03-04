package structures

import "golang.org/x/exp/constraints"

type OrderedStructure[Entry constraints.Ordered] interface {
	InsertSorted(Entry)
	BinarySearch(Entry) int
	LowerBound(Entry) int
	UpperBound(Entry) int
	GetLength() int
	At(int) Entry
}
