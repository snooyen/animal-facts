package publisher

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/go-kit/kit/log"
	"github.com/go-redis/redis/v8"

	pb "github.com/snooyen/animal-facts/facts/pb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

var (
	ErrAnimalUnsupported = errors.New("Unsupported Animal")
)

// Service describes a service that publishs the web for animal-facts
type Service interface {
	Publish(ctx context.Context, animal string) (PublishResponse, error)
}

type service struct {
	rdb       *redis.Client
	logger    log.Logger
	facts     pb.FactsClient
	scheduler *gocron.Scheduler
}

// ServiceMiddleware is a chainable behavior modifier for Service.
type ServiceMiddleware func(Service) Service

func New(ctx context.Context, redisClient *redis.Client, logger log.Logger, factsApiAddr string, cronSchedule string) (Service, error) {
	s := service{
		rdb:       redisClient,
		logger:    logger,
		facts:     NewFactsClient(factsApiAddr),
		scheduler: gocron.NewScheduler(time.UTC),
	}
	if err := s.schedulePublishJobs(ctx, cronSchedule); err != nil {
		return s, err
	}

	s.scheduler.StartAsync()

	return s, nil
}

func (s service) Publish(ctx context.Context, animal string) (response PublishResponse, err error) {
	s.logger.Log("msg", "start", "method", "publish", "animal", animal)
	response = PublishResponse{}

	req := pb.GetRandAnimalFactRequest{
		Animal: animal,
	}

	res, err := s.facts.GetRandAnimalFact(ctx, &req)
	if err != nil {
		return
	}
	response.Fact = res.Fact
	response.ID = res.Fact.ID

	// Send fact for approval
	approvalChan := fmt.Sprintf("approvals:%s", animal)
	approvalMsg := fmt.Sprintf("%s:%s", res.Fact.ID, response.Fact)
	err = s.rdb.Publish(ctx, approvalChan, approvalMsg).Err()

	s.logger.Log("msg", "end", "method", "publish", "fact", response.Fact)
	return
}

// SchedulePublishJobs runs service.Publish on a regular interval
func (s service) schedulePublishJobs(ctx context.Context, cronSchedule string) error {
	s.scheduler.TagsUnique()
	// Get list of Animals to publish fact(s) for
	res, err := s.facts.GetAnimals(ctx, new(emptypb.Empty))
	if err != nil {
		return err
	}

	s.logger.Log()

	for _, animal := range res.Animals {
		s.logger.Log("method", "scheduler", "animal", animal, "cronSchedule", cronSchedule)
		s.scheduler.Cron(cronSchedule).Tag(animal).Do(s.Publish, ctx, animal)
	}

	return nil
}
