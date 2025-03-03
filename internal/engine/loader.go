package engine

import (
	"Boolean-IR-System/internal"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func isLegiable(ext string) bool {
	return ext == ".txt"
}

func Load(path string) []*internal.Document {
	var docs []*internal.Document

	files, err := os.ReadDir(path)
	if err != nil {
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			fmt.Fprint(os.Stderr, "[Loader]: File does not exist\n")
			os.Exit(1)
		}

		if isLegiable(filepath.Ext(path)) {
			return []*internal.Document{{Path: path, Name: filepath.Base(path)}}
		}
		return nil
	}

	for _, file := range files {
		if file.Type().IsDir() {
			docs = append(docs, Load(filepath.Join(path, file.Name()))...)
		} else if isLegiable(filepath.Ext(file.Name())) {
			docs = append(docs, &internal.Document{Path: path, Name: file.Name()})
		}
	}

	return docs
}
