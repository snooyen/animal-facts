package subscriptions

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-kit/log"
	"github.com/go-redis/redis/v8"
	"github.com/snooyen/animal-facts/users/pb"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"google.golang.org/grpc"
)

const (
	subscriptionKey = "subscription"
)

type subscription struct {
	Animal       string
	rdb          *redis.Client
	logger       log.Logger
	twilio       *twilio.RestClient
	twilioNumber string
	users        pb.UsersClient
}

func New(animal string, logger log.Logger, rdb *redis.Client, usersConn *grpc.ClientConn, twilioClient *twilio.RestClient, twilioNumber string) *subscription {
	return &subscription{
		Animal:       animal,
		logger:       logger,
		rdb:          rdb,
		users:        pb.NewUsersClient(usersConn),
		twilio:       twilioClient,
		twilioNumber: twilioNumber,
	}
}

func (s *subscription) Start(ctx context.Context) error {
	chName := subscriptionKey + ":" + s.Animal
	s.logger.Log("method", "Start", "message", "start subscription", "channel", chName)
	ch := s.rdb.Subscribe(ctx, chName).Channel()

	rateLimit := time.NewTicker(1 * time.Second)
	defer rateLimit.Stop()

	for msg := range ch {
		<-rateLimit.C
		select {
		case <-ctx.Done():
			return nil
		default:
			s.logger.Log("method", "Start", "channel", chName, "message", msg.Payload)
			subscribers, err := s.rdb.SMembers(ctx, "subscribers:"+s.Animal).Result()
			if err != nil {
				s.logger.Log("method", "Start", "channel", chName, "error", err)
			}
			for _, subscriber := range subscribers {
				go s.sendFact(ctx, subscriber, msg.Payload)
			}
		}
	}
	return nil
}

func (s *subscription) sendFact(ctx context.Context, subscriber string, fact string) {
	// Parse subscriber string to uuid and get user
	uuid, err := strconv.Atoi(subscriber)
	if err != nil {
		s.logger.Log("method", "sendFact", "error", err)
		return
	}
	r, err := s.users.GetUser(ctx, &pb.GetUserRequest{ID: int64(uuid)})
	if err != nil {
		s.logger.Log("method", "sendFact", "error", err)
	}

	// Send fact to user
	s.logger.Log("method", "sendFact", "user", r.GetPhone(), "fact", fact)
	date := time.Now().Format("2006-01-02")
	msg := fmt.Sprintf("%s's %s fact for %s\n\n%s", r.GetName(), s.Animal, date, fact)
	params := &openapi.CreateMessageParams{}
	params.SetFrom(s.twilioNumber)
	params.SetTo(r.GetPhone())

	if r.GetFactsReceived() < 1 {
		s.logger.Log("method", "sendFact", "user", r.GetPhone(), "message", "first fact")
		params.SetBody(r.GetWelcomeMessage() + "\n\n" + msg)
	} else {
		params.SetBody(msg)
	}

	rsp, err := s.twilio.ApiV2010.CreateMessage(params)
	if err != nil {
		s.logger.Log("method", "sendFact", "error", err, "params", params)
	}
	s.rdb.HIncrBy(ctx, "user:"+subscriber, "FactsReceived", 1)
	response, _ := json.Marshal(*rsp)
	s.logger.Log("method", "sendFact", "user", r.GetPhone(), "response", string(response))
}
