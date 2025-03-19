package engine

import (
	"errors"
	"fmt"
	"github.com/CS80-Team/Goolean/internal"
	"os"
	"path/filepath"
	"sync"
)

func isLegible(ext string) bool {
	return ext == ".txt"
}

func (e *Engine) LoadDirectory(path string) {
	docs := LoadDocuments(path)
	for _, doc := range docs {
		e.AddDocument(doc)
	}
}

func LoadDocuments(path string) []*internal.Document {
	docCh := make(chan *internal.Document, 100)
	errCh := make(chan error, 10)

	// To wait for all goroutines to finish
	var wg sync.WaitGroup

	// add a goroutine in the WaitGroup
	wg.Add(1)
	go LoadDocumentsRecursive(path, &wg, docCh, errCh)

	go func() {
		wg.Wait()
		close(docCh)
		close(errCh)
	}()

	errList := collectErrors(errCh)
	documents := collectDocuments(docCh)

	if len(errList) > 0 {
		fmt.Printf("[Loader]: %d errors occurred during loading\n", len(errList))
	}
	return documents
}

func collectErrors(errCh chan error) []error {
	var errList []error
	go func() {
		for err := range errCh {
			errList = append(errList, err)
			fmt.Printf("[Loader]: Warning: %v\n", err) // Log the error
		}
	}()
	return errList
}

func collectDocuments(docCh chan *internal.Document) []*internal.Document {
	var documents []*internal.Document
	for doc := range docCh {
		documents = append(documents, doc)
	}
	return documents
}

func LoadDocumentsRecursive(path string, wg *sync.WaitGroup, docCh chan<- *internal.Document, errCh chan<- error) {
	defer wg.Done()

	files, err := os.ReadDir(path)
	if err != nil {
		if os.IsPermission(err) {
			errCh <- fmt.Errorf("permission denied: %s", path)
			return
		}

		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			errCh <- fmt.Errorf("path %s does not exist", path)
			return
		}

		if isLegible(filepath.Ext(path)) {
			docCh <- &internal.Document{DirectoryPath: path, Name: filepath.Base(path)}
			return
		}
		errCh <- err
		return
	}

	for _, file := range files {
		filePath := filepath.Join(path, file.Name())
		if file.IsDir() {
			// If it's a directory, start a new goroutine to process it
			wg.Add(1)
			go LoadDocumentsRecursive(filePath, wg, docCh, errCh)
		} else if isLegible(filepath.Ext(file.Name())) {
			docCh <- internal.NewDocument(filePath)
		}
	}
}
