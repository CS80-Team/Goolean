package dto

import (
	"github.com/CS80-Team/Goolean/internal"
	"github.com/CS80-Team/Goolean/internal/transport"
)

func MapEngineDocumentToGRPCDocument(engineDoc *internal.Document) *transport.Document {
	return &transport.Document{
		Id:   int32(engineDoc.ID),
		Name: engineDoc.Name,
		Path: engineDoc.DirectoryPath,
		Ext:  engineDoc.Ext,
	}
}

func MapGRPCDocumentToEngineDocument(grpcDoc *transport.Document) *internal.Document {
	return &internal.Document{
		ID:            int(grpcDoc.Id),
		Name:          grpcDoc.Name,
		DirectoryPath: grpcDoc.Path,
		Ext:           grpcDoc.Ext,
	}
}
