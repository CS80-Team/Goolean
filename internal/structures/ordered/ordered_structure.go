package ordered

import "golang.org/x/exp/constraints"

type OrderedStructure[Entry constraints.Ordered] interface {
	InsertSorted(Entry)

    BinarySearch(Entry) int
	LowerBound(Entry) int
	UpperBound(Entry) int

    GetLength() int
	IsEmpty() bool

    At(int) Entry

	Intersection(OrderedStructure[Entry]) OrderedStructure[Entry]
    Union(OrderedStructure[Entry]) OrderedStructure[Entry]
    Complement(int) OrderedStructure[Entry]
}
