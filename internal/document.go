package internal

import (
	"path/filepath"
	"strconv"
)

type Document struct {
	ID            int
	Name          string
	DirectoryPath string
	Ext           string
}

func NewDocument(path string) *Document {
	doc := &Document{}
	doc.Ext = filepath.Ext(path)
	doc.Name = filepath.Base(path)[:len(filepath.Base(path))-len(doc.Ext)]
	doc.DirectoryPath = filepath.Dir(path)
	return doc
}

func (d *Document) GetFileNameWithExt() string {
	return d.Name + d.Ext
}

func (d *Document) GetFilePath() string {
	return filepath.Join(d.DirectoryPath, d.Name+d.Ext)
}

func (d *Document) String() string {
	return "ID: " + strconv.Itoa(d.ID) + "\n" +
		"Name: " + d.Name + "\n" +
		"DirectoryPath: " + d.DirectoryPath + "\n" +
		"Ext: " + d.Ext + "\n"
}
