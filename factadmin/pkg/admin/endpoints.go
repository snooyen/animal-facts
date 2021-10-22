package admin

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type Endpoints struct {
	ApproveFactEndpoint endpoint.Endpoint
	DeferFactEndpoint   endpoint.Endpoint
	DeleteFactEndpoint  endpoint.Endpoint
	HandleSMSEndpoint   endpoint.Endpoint
}

func MakeServerEndpoints(s Service, logger log.Logger) Endpoints {
	endpoints := Endpoints{
		ApproveFactEndpoint: MakeApproveFactEndpoint(s),
		DeferFactEndpoint:   MakeDeferFactEndpoint(s),
		DeleteFactEndpoint:  MakeDeleteFactEndpoint(s),
		HandleSMSEndpoint:   MakeHandleSMSEndpoint(s),
	}

	endpoints.ApproveFactEndpoint = EndpointLoggingMiddleware(log.With(logger, "method", "ApproveFact"))(endpoints.ApproveFactEndpoint)
	endpoints.DeferFactEndpoint = EndpointLoggingMiddleware(log.With(logger, "method", "DeferFact"))(endpoints.DeferFactEndpoint)
	endpoints.DeleteFactEndpoint = EndpointLoggingMiddleware(log.With(logger, "method", "DeleteFact"))(endpoints.DeleteFactEndpoint)
	endpoints.HandleSMSEndpoint = EndpointLoggingMiddleware(log.With(logger, "method", "HandleSMS"))(endpoints.HandleSMSEndpoint)

	return endpoints
}

func MakeApproveFactEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return struct{}{}, nil
	}
}

func MakeDeferFactEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return struct{}{}, nil
	}
}

func MakeDeleteFactEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return struct{}{}, nil
	}
}

func MakeHandleSMSEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		e := s.HandleSMS(ctx)
		return handleSMSResponse{Err: e}, nil
	}
}

type approveFactRequest struct {
	ID int64
}

type approveFactResponse struct {
	Err error `json:"err,omitempty"`
}

func (r approveFactResponse) error() error { return r.Err }

type deferFactRequest struct {
	ID int64
}

type deferFactResponse struct {
	Err error `json:"err,omitempty"`
}

func (r deferFactResponse) error() error { return r.Err }

type deleteFactRequest struct {
	ID int64
}

type deleteFactResponse struct {
	Err error `json:"err,omitempty"`
}

func (r deleteFactResponse) error() error { return r.Err }

type handleSMSRequest struct {
	req string
}

type handleSMSResponse struct {
	Err error `json:"err,omitempty"`
}

func (r handleSMSResponse) error() error { return r.Err }
