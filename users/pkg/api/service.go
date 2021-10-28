package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-redis/redis/v8"
)

// animal-facts users api service
type Service interface {
	CreateUser(ctx context.Context, user User) (int64, error)
	GetUser(ctx context.Context, uuid int64) (User, error)
	DeleteUser(ctx context.Context, uuid int64) error
}

type User struct {
	ID             int64    `json:"id" redis:"ID"`
	Name           string   `json:"name" redis:"Name"`
	Phone          string   `json:"phone" redis:"Phone"`
	WelcomeMessage string   `json:"welcomeMessage" redis:"WelcomeMessage"`
	Subscriptions  []string `json:"subscriptions" redis:"Subscriptions"`
	Deleted        bool     `json:"deleted" redis:"Deleted"`
}

type service struct {
	rdb *redis.Client
}

var (
	ErrNotFound      = errors.New("user not found")
	ErrAlreadyExists = errors.New("user phone already exists")

	phoneSetKey            = "phones"
	nextUserIDKey          = "next_uid"
	userHashPrefix         = "user:"
	subscriptionsSetPrefix = "subscribers:"
)

func New(redisClient *redis.Client, logger log.Logger) Service {
	var s Service
	{
		s = service{rdb: redisClient}
		s = ServiceLoggingMiddleware(logger)(s)
	}

	return s
}

func (s service) CreateUser(ctx context.Context, user User) (uuid int64, err error) {
	uuid = -1
	// Check if phone already exists
	phoneExists, err := s.rdb.SIsMember(ctx, phoneSetKey, user.Phone).Result()
	if err != nil {
		return
	}
	if phoneExists {
		return -1, ErrAlreadyExists
	}

	// Add phone to set
	err = s.rdb.SAdd(ctx, phoneSetKey, user.Phone).Err()
	if err != nil {
		return
	}
	// Get next user id
	uuid, err = s.rdb.Incr(ctx, nextUserIDKey).Result()
	if err != nil {
		return
	}

	// TODO: Provide a default welcome message
	// TODO: Track last known message
	// Save User
	key := fmt.Sprintf("%s%d", userHashPrefix, uuid)
	if _, err = s.rdb.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, key, "ID", uuid)
		rdb.HSet(ctx, key, "Name", user.Name)
		rdb.HSet(ctx, key, "Phone", user.Phone)
		rdb.HSet(ctx, key, "WelcomeMessage", user.WelcomeMessage)
		rdb.HSet(ctx, key, "Subscriptions", fmt.Sprintf("%v", user.Subscriptions))
		rdb.HSet(ctx, key, "Deleted", false)
		return nil
	}); err != nil {
		return
	}

	// Add user to subscriptions
	for _, sub := range user.Subscriptions {
		key = fmt.Sprintf("%s%s", subscriptionsSetPrefix, sub)
		err = s.rdb.SAdd(ctx, key, uuid).Err()
		if err != nil {
			return
		}
	}

	return
}

// GetUser retrieves a single user from redis by its id
func (s service) GetUser(ctx context.Context, uuid int64) (User, error) {
	return User{}, fmt.Errorf("not implemented")
}

// DeleteUser deletes a user given its id
func (s service) DeleteUser(ctx context.Context, uuid int64) error {
	return fmt.Errorf("not implemented")
}
