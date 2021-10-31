package subscriptions

import (
	"context"
	"time"

	"github.com/go-kit/log"
	"github.com/go-redis/redis/v8"
)

const (
	subscriptionKey = "subscription"
)

type subscription struct {
	Animal string
	rdb    *redis.Client
	logger log.Logger
}

func New(animal string, logger log.Logger, rdb *redis.Client) *subscription {
	return &subscription{
		Animal: animal,
		logger: logger,
		rdb:    rdb,
	}
}

func (s *subscription) Start(ctx context.Context) error {
	chName := subscriptionKey + ":" + s.Animal
	s.logger.Log("subscribing to %s", chName)
	ch := s.rdb.Subscribe(ctx, chName).Channel()

	rateLimit := time.NewTicker(1 * time.Second)
	defer rateLimit.Stop()

	for msg := range ch {
		<-rateLimit.C
		select {
		case <-ctx.Done():
			return nil
		default:
			s.logger.Log("received message: %s", msg.Payload)
		}
	}
	return nil
}
