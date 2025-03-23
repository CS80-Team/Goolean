package service

import (
	"context"

	"github.com/CS80-Team/Goolean/internal/dto"
	"github.com/CS80-Team/Goolean/internal/engine"
	pb "github.com/CS80-Team/Goolean/internal/transport/load"
)

type LoadServer struct {
	pb.UnimplementedLoadServer
	engine *engine.Engine
}

func NewLoadServer(engine *engine.Engine) *LoadServer {
	return &LoadServer{
		engine: engine,
	}
}

// TODO: Send a whole file and load it into the engine, instead of local
func (qs *LoadServer) Load(ctx context.Context, request *pb.LoadRequest) (*pb.LoadResponse, error) {
	for _, doc := range request.Documents {
		qs.engine.AddDocument(dto.MapGRPCDocumentToEngineDocument(doc))
	}

	return &pb.LoadResponse{}, nil

	//switch {
	//case ctx.Done() != nil:
	//	return nil, ctx.Err()
	//default:
	//	for _, doc := range request.Documents {
	//		qs.engine.AddDocument(dto.MapGRPCDocumentToEngineDocument(doc))
	//	}
	//
	//	return &pb.LoadResponse{}, nil
	//}
}
