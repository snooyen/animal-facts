package admin

import (
	"context"
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

func (mw logmw) Approve(ctx context.Context, animal string, fact string, action string) (response ApprovalResponse, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "approve",
			"action", action,
			"animal", animal,
			"fact", fact,
			"msg", response.Msg,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	response, err = mw.Service.Approve(ctx, animal, fact, action)
	return
}
