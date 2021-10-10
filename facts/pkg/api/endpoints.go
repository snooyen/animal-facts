package api

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateFactEndpoint        endpoint.Endpoint
	GetFactEndpoint           endpoint.Endpoint
	DeleteFactEndpoint        endpoint.Endpoint
	GetAnimalsEndpoint        endpoint.Endpoint
	GetRandAnimalFactEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		CreateFactEndpoint:        MakeCreateFactEndpoint(s),
		GetFactEndpoint:           MakeGetFactEndpoint(s),
		DeleteFactEndpoint:        MakeDeleteFactEndpoint(s),
		GetAnimalsEndpoint:        MakeGetAnimalsEndpoint(s),
		GetRandAnimalFactEndpoint: MakeGetRandAnimalFactEndpoint(s),
	}
}

func MakeCreateFactEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createFactRequest)
		e := s.CreateFact(ctx, req.Fact.Animal, req.Fact.Fact)
		return createFactResponse{Err: e}, nil
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
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
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

type createFactRequest struct {
	Fact Fact
}

type createFactResponse struct {
	Err error `json:"err,omitempty"`
}

type getFactRequest struct {
	ID int64
}

type getFactResponse struct {
	Fact Fact
	Err  error `json:"err,omitempty"`
}

type deleteFactRequest struct {
	ID int64
}

type deleteFactResponse struct {
	Err error `json:"err,omitempty"`
}

type getAnimalsResponse struct {
	Animals []string `json:"animals"`
	Err     error    `json:"err,omitempty"`
}

type getRandAnimalFactRequest struct {
	Animal string
}

type getRandAnimalFactResponse struct {
	Fact Fact
	Err  error `json:"err,omitempty"`
}
