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

			s.logger.Log("msg", factText)
		})

		c.OnRequest(func(r *colly.Request) {
			visited = append(visited, fmt.Sprintf("%s", r.URL))
		})

		c.Visit(url)
	default:
	}

	return visited, nil
}
