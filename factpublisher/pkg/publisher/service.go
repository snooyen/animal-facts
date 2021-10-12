package publisher

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"

	pb "github.com/snooyen/animal-facts/facts/pb"
)

var (
	ErrAnimalUnsupported = errors.New("Unsupported Animal")
)

// Service describes a service that publishs the web for animal-facts
type Service interface {
	Publish(ctx context.Context, animal string) (PublishResponse, error)
}

type service struct {
	rdb   *redis.Client
	facts pb.FactsClient
}

// ServiceMiddleware is a chainable behavior modifier for Service.
type ServiceMiddleware func(Service) Service

func New(redisClient *redis.Client, factsApiAddr string) Service {
	return service{
		rdb:   redisClient,
		facts: NewFactsClient(factsApiAddr),
	}
}

func (s service) Publish(ctx context.Context, animal string) (response PublishResponse, err error) {
	response = PublishResponse{}

	req := pb.GetRandAnimalFactRequest{
		Animal: animal,
	}

	res, err := s.facts.GetRandAnimalFact(ctx, &req)
	if err != nil {
		return
	}
	response.Fact = res.Fact

	// Send fact for approval
	approvalChan := fmt.Sprintf("approvals:%s", animal)
	approvalMsg := fmt.Sprintf("%s:%s", strconv.FormatFloat(response.Score, 'f', -1, 64), response.Fact)
	err = s.rdb.Publish(ctx, approvalChan, approvalMsg).Err()

	return
}
