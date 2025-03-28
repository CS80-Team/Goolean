package ordered

import "golang.org/x/exp/constraints"

type SetOperations[Entry constraints.Ordered] interface {
	Intersection(OrderedStructure[Entry]) OrderedStructure[Entry]
	Union(OrderedStructure[Entry]) OrderedStructure[Entry]
	Complement(int) OrderedStructure[Entry]
}

type OrderedStructure[Entry constraints.Ordered] interface {
	SetOperations[Entry]

	InsertSorted(Entry)

	GetLength() int
	IsEmpty() bool
	At(int) Entry
}

type SearchableOrderedSet[Entry constraints.Ordered] interface {
	BinarySearch(Entry) int
	LowerBound(Entry) int
	UpperBound(Entry) int
}
