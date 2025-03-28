package ordered

import "golang.org/x/exp/constraints"

type BasicStructure[Entry constraints.Ordered] interface {
	GetLength() int
	IsEmpty() bool
	At(int) Entry
}

type SetOperations[Entry constraints.Ordered] interface {
	Intersection(OrderedStructure[Entry]) OrderedStructure[Entry]
	Union(OrderedStructure[Entry]) OrderedStructure[Entry]
	Complement(int) OrderedStructure[Entry]
}

type OrderedStructure[Entry constraints.Ordered] interface {
	BasicStructure[Entry]
	SetOperations[Entry]

	InsertSorted(Entry)
	BinarySearch(Entry) int
	LowerBound(Entry) int
	UpperBound(Entry) int
}
