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

func (mw serviceLoggingMiddleware) CreateFact(ctx context.Context, animal string, fact string) (ufid int64, err error) {
	defer func() {
		mw.logger.Log("method", "CreateFact", "animal", animal, "ufid", ufid, "err", err)
	}()

	return mw.next.CreateFact(ctx, animal, fact)
}

func (mw serviceLoggingMiddleware) GetFact(ctx context.Context, ufid int64) (fact Fact, err error) {
	defer func() {
		mw.logger.Log("method", "GetFact", "ufid", ufid, "err", err)
	}()

	return mw.next.GetFact(ctx, ufid)
}

func (mw serviceLoggingMiddleware) DeleteFact(ctx context.Context, ufid int64) (err error) {
	defer func() {
		mw.logger.Log("method", "DeleteFact", "ufid", ufid, "err", err)
	}()

	return mw.next.DeleteFact(ctx, ufid)
}

func (mw serviceLoggingMiddleware) GetAnimals(ctx context.Context) (animal []string, err error) {
	defer func() {
		mw.logger.Log("method", "GetAnimals", "err", err)
	}()

	return mw.next.GetAnimals(ctx)
}

func (mw serviceLoggingMiddleware) GetRandAnimalFact(ctx context.Context, animal string) (fact Fact, err error) {
	defer func() {
		mw.logger.Log("method", "GetRandAnimalFact", "ufid", fact.ID, "err", err)
	}()

	return mw.next.GetRandAnimalFact(ctx, animal)
}

func (mw serviceLoggingMiddleware) PublishFact(ctx context.Context, animal string) (fact Fact, err error) {
	defer func() {
		mw.logger.Log("method", "PublishFact", "animal", animal, "err", err)
	}()

	return mw.next.PublishFact(ctx, animal)
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
