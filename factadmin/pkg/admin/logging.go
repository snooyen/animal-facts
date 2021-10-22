package admin

import (
	"context"
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

func (mw serviceLoggingMiddleware) ApproveFact(ctx context.Context, ufid int64) (err error) {
	defer func() {
		mw.logger.Log("method", "ApproveFact", "ufid", ufid, "err", err)
	}()

	return mw.next.ApproveFact(ctx, ufid)
}

func (mw serviceLoggingMiddleware) DeferFact(ctx context.Context, ufid int64) (err error) {
	defer func() {
		mw.logger.Log("method", "DeferFact", "ufid", ufid, "err", err)
	}()

	return mw.next.DeferFact(ctx, ufid)
}

func (mw serviceLoggingMiddleware) DeleteFact(ctx context.Context, ufid int64) (err error) {
	defer func() {
		mw.logger.Log("method", "DeleteFact", "ufid", ufid, "err", err)
	}()

	return mw.next.DeleteFact(ctx, ufid)
}

func (mw serviceLoggingMiddleware) ProcessApprovalRequests(ctx context.Context) (err error) {
	defer func() {
		mw.logger.Log("method", "ProcessApprovalRequests", "err", err)
	}()

	return mw.next.ProcessApprovalRequests(ctx)
}

func (mw serviceLoggingMiddleware) HandleSMS(ctx context.Context, body string) (b string, err error) {
	defer func() {
		mw.logger.Log("method", "HandleSMS", "err", err, "body", b)
	}()
	b, err = mw.next.HandleSMS(ctx, body)
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
