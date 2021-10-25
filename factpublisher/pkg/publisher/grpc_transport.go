package publisher

import (
	"context"

	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"

	"github.com/snooyen/animal-facts/factpublisher/pb"
)

type grpcServer struct {
	publishFact grpctransport.Handler
	pb.UnimplementedFactsServer
}

func NewGRPCServer(endpoints Endpoints, logger log.Logger) pb.FactsServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return &grpcServer{
		publishFact: grpctransport.NewServer(
			endpoints.PublishFactEndpoint,
			decodeGRPCPublishRequest,
			encodeGRPCPublishFactResponse,
			options...,
		),
	}
}

func (s *grpcServer) PublishFact(ctx context.Context, req *pb.PublishRequest) (*pb.PublishReply, error) {
	_, rep, err := s.publishFact.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.PublishReply), nil
}

func decodeGRPCPublishRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.PublishRequest)
	return publishFactRequest{Animal: req.Animal}, nil
}

func encodeGRPCPublishFactResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(publishFactResponse)
	return &pb.PublishReply{Err: errToStr(resp.Err)}, nil
}

func errToStr(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
