package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/CS80-Team/Goolean/internal"
	"github.com/CS80-Team/Goolean/internal/engine"
	pb "github.com/CS80-Team/Goolean/internal/transport/file"
)

type FileServer struct {
	pb.UnimplementedFileServiceServer
	uploadDir string
	engine    *engine.Engine
}

func NewFileServer(uploadDir string, engine *engine.Engine) *FileServer {
	return &FileServer{
		uploadDir: uploadDir,
		engine:    engine,
	}
}

func (s *FileServer) UploadFile(stream pb.FileService_UploadFileServer) error {
	log.Print("New file upload request")

	var fileName string
	var fileExt string
	fileData := make(map[int32][]byte)
	var maxChunkID int32 = -1

	for {
		fileChunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to receive chunk: %v", err)
		}

		if fileName == "" {
			fileName = fileChunk.Name
			fileExt = fileChunk.Ext
		}

		for _, chunk := range fileChunk.Chunks {
			fileData[chunk.Id] = chunk.Data
			if chunk.Id > maxChunkID {
				maxChunkID = chunk.Id
			}
		}
	}

	if err := os.MkdirAll(s.uploadDir, 0755); err != nil {
		return fmt.Errorf("failed to create upload directory: %v", err)
	}

	fullPath := filepath.Join(s.uploadDir, fileName+"."+fileExt)
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	for i := int32(0); i <= maxChunkID; i++ {
		data, ok := fileData[i]
		if !ok {
			return fmt.Errorf("missing chunk %d", i)
		}
		if _, err := file.Write(data); err != nil {
			return fmt.Errorf("failed to write chunk %d: %v", i, err)
		}
	}

	s.engine.AddDocument(internal.NewDocument(fullPath))

	return stream.SendAndClose(&pb.FileStatus{
		Status:  "success",
		Message: fmt.Sprintf("File %s uploaded successfully", fileName),
	})
}

func (s *FileServer) validateFile(id *pb.DocumentID) (*internal.Document, error) {
	if id == nil {
		return nil, fmt.Errorf("document ID not provided")
	}

	doc_id, err := strconv.Atoi(id.Id)

	if err != nil {
		return nil, fmt.Errorf("document ID not provided")
	}

	if doc_id < 0 || doc_id >= s.engine.GetDocumentsSize() {
		return nil, fmt.Errorf("invalid document ID, must be in range [0, %d)", s.engine.GetDocumentsSize())
	}

	doc := s.engine.GetDocumentByID(doc_id)

	info, err := os.Stat(doc.GetFilePath())
	if err != nil {
		return nil, fmt.Errorf("file not found: %v", err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("%s is a directory, not a file", doc.GetFilePath())
	}

	return doc, nil
}

func (s *FileServer) DownloadFile(ctx context.Context, id *pb.DocumentID) (*pb.File, error) {
	log.Print("New file download request with id ", id.Id)
	doc, err := s.validateFile(id)

	if err != nil {
		return nil, err
	}

	file, err := os.Open(doc.GetFilePath())
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	chunk := &pb.File_Chunk{
		Id:   0,
		Data: data,
	}

	log.Print("File ", doc.Name, " downloaded successfully")

	return &pb.File{
		Name:   doc.Name,
		Chunks: []*pb.File_Chunk{chunk},
		Ext:    doc.Ext,
	}, nil
}
