package scraper

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-redis/redis/v8"
	"github.com/gocolly/colly"
)

const (
	animalSetKey        = "animals"
	factsSetKey         = "facts"
	nextFIDKey          = "next_fid"
	animalFactSetPrefix = "facts:"
	factHashPrefix      = "fact:"
)

var (
	ErrAnimalUnsupported = errors.New("Unsupported Animal")
)

// Service describes a service that scrapes the web for animal-facts
type Service interface {
	Scrape(ctx context.Context, animal string) ([]string, error)
}

type service struct {
	logger     log.Logger
	animalURLs map[string]string
	rdb        *redis.Client
}

// ServiceMiddleware is a chainable behavior modifier for Service.
type ServiceMiddleware func(Service) Service

func New(animalURLs map[string]string, redisClient *redis.Client, logger log.Logger) Service {
	return service{
		animalURLs: animalURLs,
		rdb:        redisClient,
		logger:     logger,
	}
}

func (s service) Scrape(ctx context.Context, animal string) (visited []string, err error) {

	// Check if animal is supported
	url, ok := s.animalURLs[animal]
	if !ok {
		err = ErrAnimalUnsupported
		return
	}

	// store animal name in animal set
	err = s.rdb.SAdd(ctx, animalSetKey, animal).Err()
	if err != nil {
		return
	}

	switch animal {
	case "elephant-seal":
		c := colly.NewCollector(
			colly.AllowedDomains("elephantseal.org"), //FIXME: Extract from url
		)

		// For every h2 tag with entry-title class
		//		visit its child a element with href attribute
		c.OnHTML("*.entry-title", func(e *colly.HTMLElement) {
			t := e.ChildAttr("a", "href")
			c.Visit(t)
		})

		c.OnHTML("div.et_pb_text_inner", func(e *colly.HTMLElement) {
			factText := e.ChildText("p")

			// REDIS: Check if Fact Exists or if facts set is empty
			factSetCard, err := s.rdb.SCard(ctx, factsSetKey).Result()
			if err != nil {
				return
			}
			factExists, err := s.rdb.SIsMember(ctx, factsSetKey, factText).Result()
			if err != nil {
				return
			}
			if factExists == false || factSetCard == 0 {
				s.addFact(ctx, factText, animal)
			}
		})

		c.OnRequest(func(r *colly.Request) {
			visited = append(visited, fmt.Sprintf("%s", r.URL))
		})

		c.Visit(url)
	default:
	}

	return visited, nil
}

func (s service) addFact(ctx context.Context, factText string, animal string) (err error) {
	// Add fact to master fact set
	err = s.rdb.SAdd(ctx, factsSetKey, factText).Err()
	if err != nil {
		return
	}

	// get next fact id
	thisFID, err := s.rdb.Incr(ctx, nextFIDKey).Result()
	if err != nil {
		return
	}
	// store fact in facts hash
	key := fmt.Sprintf("%s%d", factHashPrefix, thisFID)
	hashFields := map[string]interface{}{
		"Animal": animal,
		"Fact":   factText,
	}
	err = s.rdb.HSet(ctx, key, hashFields).Err()
	if err != nil {
		return
	}

	// add fact id to animal fact set
	key = fmt.Sprintf("%s%s", animalFactSetPrefix, animal)
	err = s.rdb.SAdd(ctx, key, thisFID).Err()
	if err != nil {
		return
	}

	s.logger.Log(
		"method", "scrape",
		"animal", animal,
		"msg", "new fact added",
		"fid", thisFID,
		"fact", factText,
	)
	return
}
