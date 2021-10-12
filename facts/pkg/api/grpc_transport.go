package api

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	"github.com/snooyen/animal-facts/facts/pb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type grpcServer struct {
	createFact        grpctransport.Handler
	deleteFact        grpctransport.Handler
	getFact           grpctransport.Handler
	getAnimals        grpctransport.Handler
	getRandAnimalFact grpctransport.Handler
	pb.UnimplementedFactsServer
}

func NewGRPCServer(endpoints Endpoints, logger log.Logger) pb.FactsServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return &grpcServer{
		createFact: grpctransport.NewServer(
			endpoints.CreateFactEndpoint,
			decodeGRPCCreateFactRequest,
			encodeGRPCCreateFactResponse,
			options...,
		),
		getFact: grpctransport.NewServer(
			endpoints.GetFactEndpoint,
			decodeGRPCGetFactRequest,
			encodeGRPCGetFactResponse,
			options...,
		),
		getAnimals: grpctransport.NewServer(
			endpoints.GetAnimalsEndpoint,
			decodeGRPCGetAnimalsRequest,
			encodeGRPCGetAnimalsResponse,
			options...,
		),
		getRandAnimalFact: grpctransport.NewServer(
			endpoints.GetRandAnimalFactEndpoint,
			decodeGRPCGetRandAnimalFactRequest,
			encodeGRPCGetRandAnimalFactResponse,
			options...,
		),
		deleteFact: grpctransport.NewServer(
			endpoints.DeleteFactEndpoint,
			decodeGRPCDeleteFactRequest,
			encodeGRPCDeleteFactResponse,
			options...,
		),
	}
}

func (s *grpcServer) CreateFact(ctx context.Context, req *pb.CreateFactRequest) (*pb.CreateFactReply, error) {
	_, rep, err := s.createFact.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.CreateFactReply), nil
}

func (s *grpcServer) DeleteFact(ctx context.Context, req *pb.DeleteFactRequest) (*pb.DeleteFactReply, error) {
	_, rep, err := s.deleteFact.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.DeleteFactReply), nil
}

func (s *grpcServer) GetFact(ctx context.Context, req *pb.GetFactRequest) (*pb.GetFactReply, error) {
	_, rep, err := s.getFact.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetFactReply), nil
}

func (s *grpcServer) GetAnimals(ctx context.Context, req *emptypb.Empty) (*pb.GetAnimalsReply, error) {
	_, rep, err := s.getAnimals.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetAnimalsReply), nil
}

func (s *grpcServer) GetRandAnimalFact(ctx context.Context, req *pb.GetRandAnimalFactRequest) (*pb.GetRandAnimalFactReply, error) {
	_, rep, err := s.getRandAnimalFact.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetRandAnimalFactReply), nil
}

func decodeGRPCCreateFactRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CreateFactRequest)
	return createFactRequest{Animal: req.Animal, Fact: req.Fact}, nil
}

func encodeGRPCCreateFactResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(createFactResponse)
	return &pb.CreateFactReply{Err: errToStr(resp.Err)}, nil
}

func decodeGRPCGetFactRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetFactRequest)
	return getFactRequest{ID: req.ID}, nil
}

func encodeGRPCGetFactResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(getFactResponse)
	return &pb.GetFactReply{Err: errToStr(resp.Err)}, nil
}

func decodeGRPCDeleteFactRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.DeleteFactRequest)
	return deleteFactRequest{ID: req.ID}, nil
}

func encodeGRPCDeleteFactResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(deleteFactResponse)
	return &pb.DeleteFactReply{Err: errToStr(resp.Err)}, nil
}

func decodeGRPCGetAnimalsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return nil, nil
}

func encodeGRPCGetAnimalsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(getAnimalsResponse)
	return &pb.GetAnimalsReply{Animals: resp.Animals, Err: errToStr(resp.Err)}, nil
}

func decodeGRPCGetRandAnimalFactRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetRandAnimalFactRequest)
	return getRandAnimalFactRequest{Animal: req.Animal}, nil
}

func encodeGRPCGetRandAnimalFactResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(getRandAnimalFactResponse)
	return &pb.GetRandAnimalFactReply{Fact: resp.Fact.Fact, Err: errToStr(resp.Err)}, nil
}

func errToStr(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
