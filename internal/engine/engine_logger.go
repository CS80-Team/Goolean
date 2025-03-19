// Global logger for the engine package

package engine

import (
	"path/filepath"

	"github.com/CS80-Team/Goolean/internal"
)

const (
	// DefaultLoggerPath is the default path for the logger
	LoggerPath         = "engine.log"
	EnginePrefix       = "[Engine]: "
	LoaderPrefix       = "[Loader]: "
	IndexManagerPrefix = "[IndexManager]: "
	QueryPrefix        = "[Query]: "
	TokenizerPrefix    = "[Tokenizer]: "
)

var logger = internal.NewLogger(filepath.Join(getCurrAbsPath(), LoggerPath))

func getCurrAbsPath() string {
	f, err := filepath.Abs(".")

	if err != nil {
		panic(err)
	}

	return f
}
