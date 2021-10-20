package admin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/go-redis/redis/v8"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

var (
	ErrInvalidFactFormat         = errors.New("Invalid Fact Format (should be <SCORE>:<FACT>)")
	ErrAnimalUnsupported         = errors.New("Unsupported Animal")
	ErrApprovalActionUnsupported = errors.New("Unsupported Action")
)

// Service describes a service that publishs the web for animal-facts
type Service interface {
	Approve(ctx context.Context, animal string, fact string, action string) (ApprovalResponse, error)
	ProcessApprovalRequests(ctx context.Context) (err error)
	subscribeApprovalChannel(ctx context.Context, animal string)
}

type service struct {
	rdb    *redis.Client
	twilio *twilio.RestClient
	logger log.Logger
}

// ServiceMiddleware is a chainable behavior modifier for Service.
type ServiceMiddleware func(Service) Service

func New(redisClient *redis.Client, twilioClient *twilio.RestClient, logger log.Logger) (s Service) {
	s = service{
		rdb:    redisClient,
		twilio: twilioClient,
		logger: logger,
	}

	return
}

func (s service) Approve(ctx context.Context, animal string, fact string, action string) (response ApprovalResponse, err error) {
	response = ApprovalResponse{Animal: animal, Action: action}

	// Get the Fact Score and Fact Text
	fArray := strings.SplitN(fact, ":", 2)
	if len(fArray) != 2 {
		err = ErrInvalidFactFormat
		return
	}
	factScore := fArray[0]
	factText := fArray[1]
	response.Fact = factText

	factSetKey := fmt.Sprintf("facts:%s", animal)
	disposalSetKey := fmt.Sprintf("disposal:s", animal)
	switch action {
	case "DISPOSE": // Throw this fact away
		s.rdb.SAdd(ctx, disposalSetKey, factText)
		response.Msg = fmt.Sprintf("%s fact disposed", animal)
	case "DEFER": // Add this fact back into the fact set
		score, _ := strconv.Atoi(factScore)
		s.rdb.ZAdd(ctx, factSetKey, &redis.Z{float64(score + 1), factText})
		response.Msg = fmt.Sprintf("%s fact defered", animal)
	case "PUBLISH": // Publish this fact
		subChan := fmt.Sprintf("sub:%s", animal)
		err = s.rdb.Publish(ctx, subChan, factText).Err()
		response.Msg = fmt.Sprintf("%s fact published", animal)
	default:
		err = ErrApprovalActionUnsupported
	}

	return
}

func (s service) ProcessApprovalRequests(ctx context.Context) (err error) {
	animals, err := s.rdb.SMembers(ctx, "animals").Result()
	if err != nil {
		return
	}

	for _, animal := range animals {
		go s.subscribeApprovalChannel(ctx, animal)
	}

	return
}

func (s service) subscribeApprovalChannel(ctx context.Context, animal string) {
	chanName := fmt.Sprintf("approvals:%s", animal)
	sub := s.rdb.Subscribe(ctx, chanName)
	approvalChannel := sub.Channel()

	rateLimit := time.Tick(1 * time.Second)

	s.logger.Log("msg", "start background", "subscription", chanName)
	for msg := range approvalChannel {
		<-rateLimit
		s.logger.Log("msg", "tick", "subscription", chanName, "time", time.Now())
		select {
		case <-ctx.Done():
			return
		default:
			s.logger.Log("msg", "received message", "subscription", chanName, "text", fmt.Sprintf("%+v", msg))
			resp, err := s.sendTextForApproval(fmt.Sprintf("%+v", msg))
			s.logger.Log("msg", "sent text for approval", "subscription", chanName, "response", resp, "err", err)

		}
	}
}

func (s service) sendTextForApproval(msg string) (string, error) {
	from := "(805) 500-0105"
	to := "(626) 607-7393"

	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(from)
	params.SetBody(msg)

	resp, err := s.twilio.ApiV2010.CreateMessage(params)
	if err != nil {
		return "", err
	}
	response, _ := json.Marshal(*resp)

	return string(response), err
}
