package service

import (
	"context"
	"log"

	"github.com/CS80-Team/Goolean/internal/engine"
	"github.com/CS80-Team/Goolean/internal/service/dto"
	"github.com/CS80-Team/Goolean/internal/transport"
	pb "github.com/CS80-Team/Goolean/internal/transport/query"
)

type QueryServer struct {
	pb.UnimplementedQueryServer
	engine *engine.Engine
}

func NewQueryServer(engine *engine.Engine) *QueryServer {
	return &QueryServer{
		engine: engine,
	}
}

func (qs *QueryServer) Query(ctx context.Context, request *pb.QueryRequest) (*pb.QueryResponse, error) {
	queryLine := request.QueryLine

	log.Printf("New query: %s", queryLine)

	res, err := qs.engine.QueryString(queryLine)

	if res == nil {
		return nil, err
	}

	documents := make([]*transport.Document, res.GetLength())
	for i := 0; i < res.GetLength(); i++ {
		doc := qs.engine.GetDocumentByID(res.At(i))

		documents[i] = dto.MapEngineDocumentToGRPCDocument(doc)
	}

	return &pb.QueryResponse{
		Documents: documents,
	}, err
}
