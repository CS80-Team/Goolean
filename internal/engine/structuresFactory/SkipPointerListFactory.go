package structuresFactory

import (
	"github.com/CS80-Team/Boolean-IR-System/internal/structures/ordered"
	"golang.org/x/exp/constraints"
)

type SkipPointerListFactory[Entry constraints.Ordered] struct {
}

func NewSkipPointerListFactory[Entry constraints.Ordered]() *SkipPointerListFactory[Entry] {
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
