package factory

import (
	"github.com/CS80-Team/Goolean/internal/structures/ordered"
	"golang.org/x/exp/constraints"
)

var _ StructuresFactory[int] = &SkipPointerListFactory[int]{}

type SkipPointerListFactory[Entry constraints.Integer] struct {
}

func NewSkipPointerListFactory[Entry constraints.Integer]() *SkipPointerListFactory[Entry] {
	return &SkipPointerListFactory[Entry]{}
}

func (o SkipPointerListFactory[Entry]) New() ordered.OrderedStructure[Entry] {
	return ordered.NewSkipPointerList[Entry]()
}

func (o SkipPointerListFactory[Entry]) NewWithCapacity(capacity int) ordered.OrderedStructure[Entry] {
	return ordered.NewSkipPointerListWithCapacity[Entry](capacity)
}

func (o SkipPointerListFactory[Entry]) NewWithSlice(slice []Entry) ordered.OrderedStructure[Entry] {
	return ordered.NewSkipPointerListWithSlice[Entry](slice)
}
