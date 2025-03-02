package retrieval

import (
	"fmt"
	"os"
	"path/filepath"
)

func (rs *BooleanRetrievalSystem) LoadTXTFilesFromDir(folderPath string) error {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".txt" {
			continue
		}

		filePath := filepath.Join(folderPath, file.Name())
		content, err := os.ReadFile(filePath)

		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", filePath, err)
			continue
		}

		doc := Document{
			ID:      rs.NextDocumentID,
			Name:    file.Name(),
			Content: string(content),
		}
		rs.Documents = append(rs.Documents, doc)
		rs.NextDocumentID++
	}

	return nil
}
