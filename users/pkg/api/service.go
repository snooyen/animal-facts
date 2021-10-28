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
	subscribersSetPrefix   = "subscribers:"
	subscriptionsSetPrefix = "subscriptions:"
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

	// Save User
	key := fmt.Sprintf("%s%d", userHashPrefix, uuid)
	if _, err = s.rdb.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, key, "ID", uuid)
		rdb.HSet(ctx, key, "Name", user.Name)
		rdb.HSet(ctx, key, "Phone", user.Phone)
		rdb.HSet(ctx, key, "WelcomeMessage", user.WelcomeMessage)
		rdb.HSet(ctx, key, "Deleted", false)
		return nil
	}); err != nil {
		return
	}

	// Add user subscriptions
	subscriptionsKey := fmt.Sprintf("%s%d", subscriptionsSetPrefix, uuid)
	for _, sub := range user.Subscriptions {
		key = fmt.Sprintf("%s%s", subscribersSetPrefix, sub)
		err = s.rdb.SAdd(ctx, key, uuid).Err()
		if err != nil {
			return
		}
		err = s.rdb.SAdd(ctx, subscriptionsKey, sub).Err()
		if err != nil {
			return
		}
	}

	return
}

// GetUser retrieves a single user from redis by its id
func (s service) GetUser(ctx context.Context, uuid int64) (user User, err error) {
	user = User{}
	key := fmt.Sprintf("%s%d", userHashPrefix, uuid)
	err = s.rdb.HGetAll(ctx, key).Scan(&user)
	if user.ID == 0 {
		err = ErrNotFound
	}
	if err != nil {
		return
	}

	key = fmt.Sprintf("%s%d", subscriptionsSetPrefix, uuid)
	user.Subscriptions, err = s.rdb.SMembers(ctx, key).Result()
	if err != nil {
		return User{}, err
	}

	return
}

// DeleteUser deletes a user given its id
func (s service) DeleteUser(ctx context.Context, uuid int64) (err error) {
	user, err := s.GetUser(ctx, uuid)
	if err != nil {
		return
	}

	userKey := fmt.Sprintf("%s%d", userHashPrefix, user.ID)
	subscriptionsKey := fmt.Sprintf("%s%d", subscriptionsSetPrefix, uuid)
	if _, err = s.rdb.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.SRem(ctx, phoneSetKey, user.Phone)
		rdb.Del(ctx, subscriptionsKey)
		for _, sub := range user.Subscriptions {
			key := fmt.Sprintf("%s%s", subscribersSetPrefix, sub)
			rdb.SRem(ctx, key, uuid)
		}
		rdb.HSet(ctx, userKey, "Phone", "")
		rdb.HSet(ctx, userKey, "WelcomeMessage", "")
		rdb.HSet(ctx, userKey, "Deleted", true)
		return nil
	}); err != nil {
		return
	}

	return nil
}
