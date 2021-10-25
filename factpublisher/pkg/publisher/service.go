package publisher

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/go-kit/log"
	"github.com/go-redis/redis/v8"

	pb "github.com/snooyen/animal-facts/facts/pb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

var (
	ErrAnimalUnsupported = errors.New("unsupported animal")
)

// Service describes a service that publishs the web for animal-facts
type Service interface {
	PublishFact(ctx context.Context, animal string) (publishFactResponse, error)
}

type service struct {
	rdb       *redis.Client
	logger    log.Logger
	facts     pb.FactsClient
	scheduler *gocron.Scheduler
}

func New(ctx context.Context, redisClient *redis.Client, logger log.Logger, factsApiAddr string, cronSchedule string) (Service, error) {
	s := service{
		rdb:       redisClient,
		logger:    logger,
		facts:     NewFactsClient(ctx, factsApiAddr),
		scheduler: gocron.NewScheduler(time.UTC),
	}
	if err := s.schedulePublishJobs(ctx, cronSchedule); err != nil {
		return s, err
	}

	s.scheduler.StartAsync()

	return s, nil
}

func (s service) PublishFact(ctx context.Context, animal string) (response publishFactResponse, err error) {
	s.logger.Log("msg", "start", "method", "publish", "animal", animal)
	response = publishFactResponse{}

	req := pb.GetRandAnimalFactRequest{
		Animal: animal,
	}

	res, err := s.facts.GetRandAnimalFact(ctx, &req)
	if err != nil {
		return
	}
	response.Fact = res.Fact
	response.ID = res.ID

	// Send fact for approval
	approvalChan := fmt.Sprintf("approvals:%s", animal)
	approvalMsg := fmt.Sprintf("%d:%s", response.ID, response.Fact)
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
		s.scheduler.Cron(cronSchedule).Tag(animal).Do(s.PublishFact, ctx, animal)
	}

	return nil
}
