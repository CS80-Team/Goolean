package engine_tests

import (
	"path/filepath"

	"github.com/CS80-Team/Goolean/internal/engine"
	"github.com/CS80-Team/Goolean/internal/engine/tokenizer"
	"github.com/CS80-Team/Goolean/internal/structures/factory"
	"github.com/CS80-Team/Goolean/internal/textprocessing"
)

var (
	processor = textprocessing.NewDefaultProcessor(
		textprocessing.NewNormalizer(),
		textprocessing.NewStemmer(),
		textprocessing.NewStopWordRemover(),
	)

	tokens = map[rune]struct{}{
		' ':  {},
		'\n': {},
		',':  {},
		'?':  {},
		'!':  {},
		'.':  {},
		';':  {},
	}

	DelimiterManager = tokenizer.NewDelimiterManager(
		&tokens,
	)

	OrderedSliceIndexManager    = engine.NewIndexManager(factory.NewOrderedSliceFactory[int]())
	SkipPointerListIndexManager = engine.NewIndexManager(factory.NewSkipPointerListFactory[int]())

	datasetPath = filepath.Join("..", "..", "..", "dataset")
)

func MockEngine(path string, indexManager engine.IndexManager) *engine.Engine {
	e := engine.NewEngine(processor, DelimiterManager, indexManager)

	e.LoadDirectory(path)

	return e
}
