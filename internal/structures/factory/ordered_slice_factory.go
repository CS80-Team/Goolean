package factory

import (
	"github.com/CS80-Team/Goolean/internal/structures/ordered"
	"golang.org/x/exp/constraints"
)

var _ StructuresFactory[int] = &OrderedSliceFactory[int]{}

type OrderedSliceFactory[Entry constraints.Integer] struct {
}

func NewOrderedSliceFactory[Entry constraints.Integer]() *OrderedSliceFactory[Entry] {
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
