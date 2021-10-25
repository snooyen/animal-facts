package publisher

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type Endpoints struct {
	PublishFactEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service, logger log.Logger) Endpoints {
	endpoints := Endpoints{
		PublishFactEndpoint: MakePublishFactEndpoint(s),
	}

	endpoints.PublishFactEndpoint = EndpointLoggingMiddleware(log.With(logger, "method", "PublishFact"))(endpoints.PublishFactEndpoint)

	return endpoints
}

func MakePublishFactEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(publishFactRequest)
		rsp, e := s.PublishFact(ctx, req.Animal)
		rsp.Err = e
		return rsp, nil
	}
}

type publishFactRequest struct {
	Animal string
}

type publishFactResponse struct {
	Fact string `json:"fact"`
	ID   int64  `json:"id"`
	Err  error  `json:"err,omitempty"`
}

func (r publishFactResponse) error() error { return r.Err }
