package api

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

/* Service Logging Middleware */
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

func (mw serviceLoggingMiddleware) CreateUser(ctx context.Context, user User) (uuid int64, err error) {
	defer func() {
		mw.logger.Log("method", "CreateUser", "uuid", uuid, "err", err)
	}()

	return mw.next.CreateUser(ctx, user)
}

func (mw serviceLoggingMiddleware) GetUser(ctx context.Context, uuid int64) (user User, err error) {
	defer func() {
		mw.logger.Log("method", "GetUser", "uuid", uuid, "err", err)
	}()

	return mw.next.GetUser(ctx, uuid)
}

func (mw serviceLoggingMiddleware) DeleteUser(ctx context.Context, uuid int64) (err error) {
	defer func() {
		mw.logger.Log("method", "DeleteUser", "uuid", uuid, "err", err)
	}()

	return mw.next.DeleteUser(ctx, uuid)
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
