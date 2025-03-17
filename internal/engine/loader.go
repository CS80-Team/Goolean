package engine

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/CS80-Team/BooleanEngine/internal"
)

func isLegible(ext string) bool {
	return ext == ".txt"
}

func (e *Engine) LoadDirectory(path string) {
	var newDocs = LoadDocs(path)
	for _, doc := range newDocs {
		e.AddDocument(doc)
	}
}

func LoadDocs(path string) []*internal.Document {
	var docs []*internal.Document

	files, err := os.ReadDir(path)
	if err != nil {
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			panic(fmt.Sprintf("[Loader]: Path %s does not exist", path))
		}

		if isLegible(filepath.Ext(path)) {
			return []*internal.Document{{DirectoryPath: path, Name: filepath.Base(path)}}
		}
		return nil
	}

	for _, file := range files {
		if file.Type().IsDir() {
			docs = append(docs, LoadDocs(filepath.Join(path, file.Name()))...)
		} else if isLegible(filepath.Ext(file.Name())) {
			docs = append(docs, internal.NewDocument(filepath.Join(path, file.Name())))
		}
	}

	return docs
}
