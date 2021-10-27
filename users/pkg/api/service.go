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
	ErrAlreadyExists = errors.New("user already exists")

	masterUserSetKey = "user"
	nextUserIDKey    = "next_uid"
	userHashPrefix   = "user:"
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
	return -1, fmt.Errorf("not implemented")
}

// GetUser retrieves a single user from redis by its id
func (s service) GetUser(ctx context.Context, uuid int64) (User, error) {
	return User{}, fmt.Errorf("not implemented")
}

// DeleteUser deletes a user given its id
func (s service) DeleteUser(ctx context.Context, uuid int64) error {
	return fmt.Errorf("not implemented")
}
