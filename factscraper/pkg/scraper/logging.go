package scraper

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
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

func (mw logmw) Scrape(ctx context.Context, animal string) (visited []string, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "scrape",
			"animal", animal,
			"visited", fmt.Sprintf("%v", visited),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	visited, err = mw.Service.Scrape(ctx, animal)
	return
}
