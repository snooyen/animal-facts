package api

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-kit/log"
	"github.com/go-redis/redis/v8"
)

// animal-facts api service
type Service interface {
	CreateFact(ctx context.Context, animal string, fact string) (int64, error)
	GetFact(ctx context.Context, ufid int64) (Fact, error)
	DeleteFact(ctx context.Context, ufid int64) error
	GetAnimals(ctx context.Context) ([]string, error)
	GetRandAnimalFact(ctx context.Context, animal string) (Fact, error)
	PublishFact(ctx context.Context, animal string) (fact Fact, err error)
}

type Fact struct {
	Animal  string `json:"animal" redis:"Animal"`
	Fact    string `json:"fact" redis:"Fact"`
	ID      int64  `json:"id"`
	Deleted bool   `json:"deleted" redis:"Deleted"`
}

type service struct {
	rdb *redis.Client
}

var (
	ErrAnimalFactNotFound = errors.New("facts not found for this animal")
	ErrNotFound           = errors.New("fact not found")
	ErrAlreadyExists      = errors.New("fact already exists")

	animalsSetKey    = "animals"
	masterFactSetKey = "facts"
	nextFactIDKey    = "next_fid"
	factHashPrefix   = "fact:"
)

func New(redisClient *redis.Client, logger log.Logger) Service {
	var s Service
	{
		s = service{rdb: redisClient}
		s = ServiceLoggingMiddleware(logger)(s)
	}

	return s
}

func (s service) CreateFact(ctx context.Context, animal string, factText string) (ufid int64, err error) {
	// Check if Fact Already Exists
	factExists, err := s.rdb.SIsMember(ctx, masterFactSetKey, factText).Result()
	if err != nil {
		return
	}
	if factExists {
		return -1, ErrAlreadyExists
	}

	// Add fact to master fact set
	err = s.rdb.SAdd(ctx, masterFactSetKey, factText).Err()
	if err != nil {
		return
	}

	// Get next fact id
	ufid, err = s.rdb.Incr(ctx, nextFactIDKey).Result()
	if err != nil {
		return
	}
	// Store fact in facts hash
	key := fmt.Sprintf("%s%d", factHashPrefix, ufid)
	hashFields := map[string]interface{}{
		"Animal":  animal,
		"Fact":    factText,
		"Deleted": false,
	}
	err = s.rdb.HSet(ctx, key, hashFields).Err()
	if err != nil {
		return
	}

	// Add fact id to animal fact sorted set
	z := redis.Z{
		Member: ufid,
	}
	err = s.rdb.ZAdd(ctx, animal, &z).Err()
	if err != nil {
		return
	}

	err = s.rdb.SAdd(ctx, animalsSetKey, animal).Err()
	if err != nil {
		return
	}

	return
}

// GetFact retrieves a single fact from redis by its id
func (s service) GetFact(ctx context.Context, ufid int64) (Fact, error) {
	fact := Fact{ID: ufid}
	key := fmt.Sprintf("%s%d", factHashPrefix, ufid)
	err := s.rdb.HGetAll(ctx, key).Scan(&fact)
	if fact.Fact == "" {
		return fact, ErrNotFound
	}

	return fact, err
}

// DeleteFact deletes a fact given its id
func (s service) DeleteFact(ctx context.Context, ufid int64) error {
	fact, err := s.GetFact(ctx, ufid)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s%d", factHashPrefix, fact.ID)
	// Update Fact.Deleted
	err = s.rdb.HSet(ctx, key, "Deleted", 1).Err()
	if err != nil {
		return err
	}

	// Remove the fact from the animal's fact set
	err = s.rdb.ZRem(ctx, fact.Animal, fact.ID).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetAnimals returns a list of known animals
func (s service) GetAnimals(ctx context.Context) ([]string, error) {
	return s.rdb.SMembers(ctx, animalsSetKey).Result()
}

// GetRandAnimalFact returns a random fact for a given animal
func (s service) GetRandAnimalFact(ctx context.Context, animal string) (Fact, error) {
	var err error
	fact := Fact{Deleted: true}

	for {
		fid, err := s.getRandFactID(ctx, animal)
		if err != nil {
			return fact, err
		}
		fact, err = s.GetFact(ctx, fid)
		if err != nil {
			return fact, err
		}
		if !fact.Deleted {
			break
		}
	}

	return fact, err
}

func (s service) PublishFact(ctx context.Context, animal string) (fact Fact, err error) {
	fact, err = s.GetRandAnimalFact(ctx, animal)
	if err != nil {
		return
	}
	approvalChan := fmt.Sprintf("approvals:%s", animal)
	approvalMsg := fmt.Sprintf("%d:%s", fact.ID, fact.Fact)
	err = s.rdb.Publish(ctx, approvalChan, approvalMsg).Err()
	return
}

func (s service) getRandFactID(ctx context.Context, animal string) (int64, error) {
	result, err := s.rdb.ZRandMember(ctx, animal, 1, false).Result()
	if err != nil {
		return -1, err
	}
	if len(result) != 1 {
		return -1, ErrAnimalFactNotFound
	}
	factID, err := strconv.Atoi(result[0])
	if err != nil {
		return -1, err
	}
	return int64(factID), nil
}
