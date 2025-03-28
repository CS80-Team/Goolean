package factory

import (
	"github.com/CS80-Team/Goolean/internal/structures/ordered"
	"golang.org/x/exp/constraints"
)

type StructuresFactory[Entry constraints.Ordered] interface {
	New() ordered.OrderedStructure[Entry]
	NewWithCapacity(capacity int) ordered.OrderedStructure[Entry]
	NewWithSlice(slice []Entry) ordered.OrderedStructure[Entry]
}
