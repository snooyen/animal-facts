package publisher

import (
	"context"
	"time"

	"github.com/go-kit/log"
)

func LoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return logmw{logger, next}
	}
}

type logmw struct {
	logger log.Logger
	Service
}

func (mw logmw) Publish(ctx context.Context, animal string) (response PublishResponse, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "publish",
			"animal", animal,
			"fact", response.Fact,
			"id", response.ID,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	response, err = mw.Service.Publish(ctx, animal)
	return
}
