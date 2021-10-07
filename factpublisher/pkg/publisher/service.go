package publisher

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var (
	ErrAnimalUnsupported = errors.New("Unsupported Animal")
)

// Service describes a service that publishs the web for animal-facts
type Service interface {
	Publish(ctx context.Context, animal string) (PublishResponse, error)
}

type service struct {
	rdb *redis.Client
}

// ServiceMiddleware is a chainable behavior modifier for Service.
type ServiceMiddleware func(Service) Service

func New(redisClient *redis.Client) Service {
	return service{
		rdb: redisClient,
	}
}

func (s service) Publish(ctx context.Context, animal string) (response PublishResponse, err error) {
	response = PublishResponse{}

	disposalSetKey := fmt.Sprintf("disposal:s", animal)
	z := s.popFact(ctx, animal)

	// If fact is in disposal set, pop facts until we have a fact not in the disposal set.
	for {
		if s.rdb.SIsMember(ctx, disposalSetKey, z.Member).Val() {
			z = s.popFact(ctx, animal)
		} else {
			break
		}
	}

	response.Fact = z.Member.(string)
	response.Score = z.Score

	// Send fact for approval
	approvalChan := fmt.Sprintf("approvals:%s", animal)
	approvalMsg := fmt.Sprintf("%s:%s", strconv.FormatFloat(response.Score, 'f', -1, 64), response.Fact)
	err = s.rdb.Publish(ctx, approvalChan, approvalMsg).Err()

	return
}

func (s service) popFact(ctx context.Context, animal string) (z redis.Z) {
	// Pop a fact from the animal's fact set
	// TODO: Handle unsupported animal(s)
	// TODO: Handle empty fact set
	factSetKey := fmt.Sprintf("facts:%s", animal)
	r, _ := s.rdb.ZPopMin(ctx, factSetKey).Result()

	z = r[0] // TODO: Will panic if fact set empty

	return
}
