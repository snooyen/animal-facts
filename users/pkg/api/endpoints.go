package api

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type Endpoints struct {
	CreateUserEndpoint endpoint.Endpoint
	GetUserEndpoint    endpoint.Endpoint
	DeleteUserEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service, logger log.Logger) Endpoints {
	endpoints := Endpoints{
		CreateUserEndpoint: MakeCreateUserEndpoint(s),
		GetUserEndpoint:    MakeGetUserEndpoint(s),
		DeleteUserEndpoint: MakeDeleteUserEndpoint(s),
	}

	endpoints.CreateUserEndpoint = EndpointLoggingMiddleware(log.With(logger, "method", "CreateUser"))(endpoints.CreateUserEndpoint)
	endpoints.GetUserEndpoint = EndpointLoggingMiddleware(log.With(logger, "method", "GetUser"))(endpoints.GetUserEndpoint)
	endpoints.DeleteUserEndpoint = EndpointLoggingMiddleware(log.With(logger, "method", "DeleteUser"))(endpoints.DeleteUserEndpoint)

	return endpoints
}

func MakeCreateUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createUserRequest)
		ufid, e := s.CreateUser(ctx, req.User)
		return createUserResponse{ID: ufid, Err: e}, nil
	}
}

func MakeGetUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getUserRequest)
		user, e := s.GetUser(ctx, req.ID)
		return getUserResponse{User: user, Err: e}, nil
	}
}

func MakeDeleteUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteUserRequest)
		e := s.DeleteUser(ctx, req.ID)
		return deleteUserResponse{Err: e}, nil
	}
}

type createUserRequest struct {
	User User `json:"User"`
}

type createUserResponse struct {
	ID  int64 `json:"id"`
	Err error `json:"err,omitempty"`
}

func (r createUserResponse) error() error { return r.Err }

type getUserRequest struct {
	ID int64 `json:"id"`
}

type getUserResponse struct {
	User User  `json:"user"`
	Err  error `json:"err,omitempty"`
}

func (r getUserResponse) error() error { return r.Err }

type deleteUserRequest struct {
	ID int64 `json:"id"`
}

type deleteUserResponse struct {
	Err error `json:"err,omitempty"`
}

func (r deleteUserResponse) error() error { return r.Err }
