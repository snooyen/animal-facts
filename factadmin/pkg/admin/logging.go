package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

// ServiceMiddleware is a chainable behavior modifier for Service.
type ServiceMiddleware func(Service) Service

// ServiceLoggingMiddleware takes a logger as a dependency
// and returns a service Middleware.
func ServiceLoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return serviceLoggingMiddleware{logger, next}
	}
}

type serviceLoggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw serviceLoggingMiddleware) ProcessApprovalRequests(ctx context.Context) (err error) {
	defer func() {
		mw.logger.Log("method", "ProcessApprovalRequests", "err", err)
	}()

	return mw.next.ProcessApprovalRequests(ctx)
}

func (mw serviceLoggingMiddleware) HandleSMS(ctx context.Context, req handleSMSRequest) (b string, err error) {
	mw.logger.Log("method", "HandleSMS", "msg", "start", "req", fmt.Sprintf("%+v", req))
	b, err = mw.next.HandleSMS(ctx, req)
	mw.logger.Log("method", "HandleSMS", "rsp", b, "err", err)
	return b, err
}

/* Endpoint Logging Middleware
Returns an endpoint middleware that logs the
duration of each invocation, and the resulting error, if any.
*/
func EndpointLoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			defer func(begin time.Time) {
				logger.Log("transport_error", err, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)

		}
	}
}
