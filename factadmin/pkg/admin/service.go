package admin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/go-redis/redis/v8"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

const (
	hostname = "https://dev.seannguyen.dev"
)

var (
	ErrAnimalUnsupported         = errors.New("unsupported animal")
	ErrApprovalActionUnsupported = errors.New("unsupported action")
)

// Service describes a service that publishs the web for animal-facts
type Service interface {
	ApproveFact(ctx context.Context, ufid int64) error
	DeferFact(ctx context.Context, ufid int64) error
	DeleteFact(ctx context.Context, ufid int64) error
	HandleSMS(ctx context.Context) error
	ProcessApprovalRequests(ctx context.Context) (err error)
}

type service struct {
	rdb          *redis.Client
	twilio       *twilio.RestClient
	logger       log.Logger
	twilioNumber string
	adminNumber  string
}

func New(redisClient *redis.Client, twilioClient *twilio.RestClient, logger log.Logger, twilioNumber, adminNumber string) (s Service) {
	s = service{
		rdb:          redisClient,
		twilio:       twilioClient,
		logger:       logger,
		twilioNumber: twilioNumber,
		adminNumber:  adminNumber,
	}

	return ServiceLoggingMiddleware(logger)(s)
}

func (s service) HandleSMS(ctx context.Context) error {
	return nil
}

func (s service) ApproveFact(ctx context.Context, ufid int64) error {
	return nil
}

func (s service) DeferFact(ctx context.Context, ufid int64) error {
	return nil
}

func (s service) DeleteFact(ctx context.Context, ufid int64) error {
	return nil
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

	rateLimit := time.NewTicker(1 * time.Second)
	defer rateLimit.Stop()

	s.logger.Log("msg", "start background", "subscription", chanName)
	for msg := range approvalChannel {
		<-rateLimit.C
		s.logger.Log("msg", "tick", "subscription", chanName, "time", time.Now())
		select {
		case <-ctx.Done():
			return
		default:
			msgStr := msg.String()
			s.logger.Log("msg", "received message", "subscription", chanName, "text", msgStr)
			id := strings.TrimSpace(strings.SplitN(msgStr, ":", 4)[2])
			smsMsg := fmt.Sprintf(
				"%s fact awaiting approval:\n"+
					"%s/facts/%s\n"+
					"approve: %s/admin/approve/%s\n"+
					"defer: %s/admin/defer/%s\n"+
					"delete: %s/admin/delete/%s\n",
				animal,
				hostname, id,
				hostname, id,
				hostname, id,
				hostname, id,
			)

			resp, err := s.sendTextForApproval(smsMsg)
			s.logger.Log("msg", "sent text for approval", "subscription", chanName, "response", resp, "err", err)

		}
	}
}

func (s service) sendTextForApproval(msg string) (string, error) {
	from := s.twilioNumber
	to := s.adminNumber

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
