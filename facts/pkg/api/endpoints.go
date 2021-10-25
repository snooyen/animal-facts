package api

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type Endpoints struct {
	CreateFactEndpoint        endpoint.Endpoint
	GetFactEndpoint           endpoint.Endpoint
	DeleteFactEndpoint        endpoint.Endpoint
	GetAnimalsEndpoint        endpoint.Endpoint
	GetRandAnimalFactEndpoint endpoint.Endpoint
	PublishFactEndpoint       endpoint.Endpoint
}

func MakeServerEndpoints(s Service, logger log.Logger) Endpoints {
	endpoints := Endpoints{
		CreateFactEndpoint:        MakeCreateFactEndpoint(s),
		GetFactEndpoint:           MakeGetFactEndpoint(s),
		DeleteFactEndpoint:        MakeDeleteFactEndpoint(s),
		GetAnimalsEndpoint:        MakeGetAnimalsEndpoint(s),
		GetRandAnimalFactEndpoint: MakeGetRandAnimalFactEndpoint(s),
		PublishFactEndpoint:       MakePublishFactEndpoint(s),
	}

	endpoints.CreateFactEndpoint = EndpointLoggingMiddleware(log.With(logger, "method", "CreateFact"))(endpoints.CreateFactEndpoint)
	endpoints.GetFactEndpoint = EndpointLoggingMiddleware(log.With(logger, "method", "GetFact"))(endpoints.GetFactEndpoint)
	endpoints.DeleteFactEndpoint = EndpointLoggingMiddleware(log.With(logger, "method", "DeleteFact"))(endpoints.DeleteFactEndpoint)
	endpoints.GetAnimalsEndpoint = EndpointLoggingMiddleware(log.With(logger, "method", "GetAnimals"))(endpoints.GetAnimalsEndpoint)
	endpoints.GetRandAnimalFactEndpoint = EndpointLoggingMiddleware(log.With(logger, "method", "GetRandAnimalFact"))(endpoints.GetRandAnimalFactEndpoint)
	endpoints.PublishFactEndpoint = EndpointLoggingMiddleware(log.With(logger, "method", "PublishFact"))(endpoints.PublishFactEndpoint)

	return endpoints
}

func MakeCreateFactEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createFactRequest)
		ufid, e := s.CreateFact(ctx, req.Animal, req.Fact)
		return createFactResponse{ID: ufid, Err: e}, nil
	}
}

func MakeGetFactEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getFactRequest)
		fact, e := s.GetFact(ctx, req.ID)
		return getFactResponse{Fact: fact, Err: e}, nil
	}
}

func MakeDeleteFactEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteFactRequest)
		e := s.DeleteFact(ctx, req.ID)
		return deleteFactResponse{Err: e}, nil
	}
}

func MakeGetAnimalsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (response interface{}, err error) {
		animals, e := s.GetAnimals(ctx)
		return getAnimalsResponse{Animals: animals, Err: e}, nil
	}
}

func MakeGetRandAnimalFactEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getRandAnimalFactRequest)
		fact, e := s.GetRandAnimalFact(ctx, req.Animal)
		return getRandAnimalFactResponse{Fact: fact, Err: e}, nil
	}
}

func MakePublishFactEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(publishFactRequest)
		fact, e := s.PublishFact(ctx, req.Animal)
		return publishFactResponse{Fact: fact.Fact, ID: fact.ID}, e
	}
}

type createFactRequest struct {
	Animal string `json:"Animal"`
	Fact   string `json:"Fact"`
}

type createFactResponse struct {
	ID  int64 `json:"id"`
	Err error `json:"err,omitempty"`
}

func (r createFactResponse) error() error { return r.Err }

type getFactRequest struct {
	ID int64
}

type getFactResponse struct {
	Fact Fact
	Err  error `json:"err,omitempty"`
}

func (r getFactResponse) error() error { return r.Err }

type deleteFactRequest struct {
	ID int64
}

type deleteFactResponse struct {
	Err error `json:"err,omitempty"`
}

func (r deleteFactResponse) error() error { return r.Err }

type getAnimalsResponse struct {
	Animals []string `json:"animals"`
	Err     error    `json:"err,omitempty"`
}

func (r getAnimalsResponse) error() error { return r.Err }

type getRandAnimalFactRequest struct {
	Animal string
}

type getRandAnimalFactResponse struct {
	Fact Fact
	Err  error `json:"err,omitempty"`
}

func (r getRandAnimalFactResponse) error() error { return r.Err }

type publishFactRequest struct {
	Animal string
}

type publishFactResponse struct {
	Fact string `json:"fact"`
	ID   int64  `json:"id"`
	Err  error  `json:"err,omitempty"`
}

func (r publishFactResponse) error() error { return r.Err }
