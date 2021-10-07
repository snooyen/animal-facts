package publisher

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func MakePublishEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PublishRequest)
		v, err := s.Publish(ctx, req.Animal)
		if err != nil {
			v.Err = err.Error()
		}
		return v, nil
	}
}
