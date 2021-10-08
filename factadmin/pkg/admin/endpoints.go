package admin

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func MakeApproveEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ApprovalRequest)
		v, err := s.Approve(ctx, req.Animal, req.Fact, req.Action)
		if err != nil {
			v.Err = err.Error()
		}
		return v, nil
	}
}
