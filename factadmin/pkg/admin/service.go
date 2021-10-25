package admin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/go-redis/redis/v8"
	"github.com/snooyen/animal-facts/facts/pb"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

const (
	hostname = "https://dev.seannguyen.dev"
)

var (
	ErrAnimalUnsupported    = errors.New("unsupported animal")
	ErrSMSTypeUnsupported   = errors.New("unsupported sms type")
	ErrSMSActionUnsupported = errors.New("unsupported sms action")
	ErrSMSBadFactData       = errors.New("could not convert FACT data to int")
	smsRegExp               = regexp.MustCompile(`(?m)^(?P<type>FACT|USER):(?P<action>[[:alpha:]]*):(?P<data>.*)$`)
)

// Service describes a service that publishs the web for animal-facts
type Service interface {
	HandleSMS(ctx context.Context, req handleSMSRequest) (string, error)
	ProcessApprovalRequests(ctx context.Context) (err error)
}

type service struct {
	facts        pb.FactsClient
	rdb          *redis.Client
	twilio       *twilio.RestClient
	logger       log.Logger
	twilioNumber string
	adminNumber  string
}

func New(factsClient pb.FactsClient, redisClient *redis.Client, twilioClient *twilio.RestClient, logger log.Logger, twilioNumber, adminNumber string) (s Service) {
	s = service{
		facts:        factsClient,
		rdb:          redisClient,
		twilio:       twilioClient,
		logger:       logger,
		twilioNumber: twilioNumber,
		adminNumber:  adminNumber,
	}

	return ServiceLoggingMiddleware(logger)(s)
}

func (s service) HandleSMS(ctx context.Context, req handleSMSRequest) (string, error) {

	// Parse SMS Body for Action
	match := smsRegExp.FindStringSubmatch(req.Body)
	body := make(map[string]string)
	for i, name := range smsRegExp.SubexpNames() {
		if i != 0 && name != "" {
			body[name] = match[i]
		}
	}

	switch body["type"] {
	case "FACT":
		return s.handleSMSFact(ctx, body["action"], body["data"])
	case "USER":
		return s.handleSMSUser(ctx, body["action"], body["data"])
	default:
		return "", ErrSMSTypeUnsupported
	}
}

func (s service) handleSMSFact(ctx context.Context, action string, data string) (string, error) {
	switch action {
	case "APPROVE":
		return "", fmt.Errorf("not implemented")
	case "DEFER":
		return "", fmt.Errorf("not implemented")
	case "DELETE":
		ufid, err := strconv.Atoi(data)
		if err == nil {
			return data, err
		}
		r, err := s.facts.DeleteFact(ctx, &pb.DeleteFactRequest{ID: int64(ufid)})
		if err != nil {
			return "", err
		}
		return r.Err, err
	default:
		return "", ErrSMSActionUnsupported
	}
}

func (s service) handleSMSUser(ctx context.Context, action string, data string) (string, error) {
	return fmt.Sprintf("GOT %s action on USER with data: %s", action, data), nil
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
