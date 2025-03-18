package structuresFactory

import (
	"github.com/CS80-Team/Goolean/internal/structures/ordered"
	"golang.org/x/exp/constraints"
)

type OrderedSliceFactory[Entry constraints.Ordered] struct {
}

func NewOrderedSliceFactory[Entry constraints.Ordered]() *OrderedSliceFactory[Entry] {
	return &OrderedSliceFactory[Entry]{}
}

func (o OrderedSliceFactory[Entry]) New() ordered.OrderedStructure[Entry] {
	return ordered.NewOrderedSlice[Entry]()
}

func (o OrderedSliceFactory[Entry]) NewWithCapacity(capacity int) ordered.OrderedStructure[Entry] {
	return ordered.NewOrderedSliceWithCapacity[Entry](capacity)
}

func (o OrderedSliceFactory[Entry]) NewWithSlice(slice []Entry) ordered.OrderedStructure[Entry] {
	return ordered.NewOrderedSliceWithSlice[Entry](slice)
}
