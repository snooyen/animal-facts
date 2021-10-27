package api

import (
	"context"

	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"

	"github.com/snooyen/animal-facts/users/pb"
)

type grpcServer struct {
	createUser grpctransport.Handler
	getUser    grpctransport.Handler
	deleteUser grpctransport.Handler
	pb.UnimplementedUsersServer
}

func NewGRPCServer(endpoints Endpoints, logger log.Logger) pb.UsersServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return &grpcServer{
		createUser: grpctransport.NewServer(
			endpoints.CreateUserEndpoint,
			decodeGRPCCreateUserRequest,
			encodeGRPCCreateUserResponse,
			options...,
		),
		getUser: grpctransport.NewServer(
			endpoints.GetUserEndpoint,
			decodeGRPCGetUserRequest,
			encodeGRPCGetUserResponse,
			options...,
		),
		deleteUser: grpctransport.NewServer(
			endpoints.DeleteUserEndpoint,
			decodeGRPCDeleteUserRequest,
			encodeGRPCDeleteUserResponse,
			options...,
		),
	}
}

func (s *grpcServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	_, rep, err := s.createUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.CreateUserReply), nil
}

func (s *grpcServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	_, rep, err := s.getUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetUserReply), nil
}

func (s *grpcServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserReply, error) {
	_, rep, err := s.deleteUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.DeleteUserReply), nil
}

func decodeGRPCCreateUserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CreateUserRequest)
	return createUserRequest{User: User{
		Name:           req.Name,
		Phone:          req.Phone,
		WelcomeMessage: req.WelcomeMessage,
		Subscriptions:  req.Subscriptions,
		Deleted:        false,
	}}, nil
}

func encodeGRPCCreateUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(createUserResponse)
	return &pb.CreateUserReply{Err: errToStr(resp.Err)}, nil
}

func decodeGRPCGetUserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetUserRequest)
	return getUserRequest{ID: req.ID}, nil
}

func encodeGRPCGetUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(getUserResponse)
	return &pb.GetUserReply{
		Name:           resp.User.Name,
		Phone:          resp.User.Phone,
		WelcomeMessage: resp.User.WelcomeMessage,
		Subscriptions:  resp.User.Subscriptions,
		Deleted:        resp.User.Deleted,
		Err:            errToStr(resp.Err),
	}, nil
}

func decodeGRPCDeleteUserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.DeleteUserRequest)
	return deleteUserRequest{ID: req.ID}, nil
}

func encodeGRPCDeleteUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(deleteUserResponse)
	return &pb.DeleteUserReply{Err: errToStr(resp.Err)}, nil
}

func errToStr(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
