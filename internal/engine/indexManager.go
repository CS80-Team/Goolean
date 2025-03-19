package engine

import (
	"github.com/CS80-Team/Goolean/internal/engine/structuresFactory"
	"github.com/CS80-Team/Goolean/internal/structures/ordered"
)

type IndexManager struct {
	index   map[string]ordered.OrderedStructure[int]
	factory structuresFactory.StructuresFactory[int]
}

func NewIndexManager(fa structuresFactory.StructuresFactory[int]) *IndexManager {
	return &IndexManager{
		index:   make(map[string]ordered.OrderedStructure[int]),
		factory: fa,
	}
}

func (idx *IndexManager) Put(key string, value int) {
	if _, ok := idx.index[key]; !ok {
		idx.index[key] = idx.factory.New()
	}

	idx.index[key].InsertSorted(value)
}

func (idx *IndexManager) PutSlice(key string, values []int) {
	if _, ok := idx.index[key]; !ok {
		idx.index[key] = idx.factory.NewWithSlice(values)
	} else {
		for _, value := range values {
			idx.index[key].InsertSorted(value)
		}
	}
}

func (idx *IndexManager) Get(key string) ordered.OrderedStructure[int] {
	return idx.index[key]
}

func (idx *IndexManager) Size() int {
	return len(idx.index)
}
