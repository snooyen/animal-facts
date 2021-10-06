package scraper

import (
	"context"
	"errors"
	"fmt"

	"github.com/gocolly/colly"
	"github.com/go-redis/redis/v8"
)

var (
	ErrAnimalUnsupported = errors.New("Unsupported Animal")
)

// Service describes a service that scrapes the web for animal-facts
type Service interface {
	Scrape(ctx context.Context, animal string) ([]string, error)
}

type service struct {
	animalURLs map[string]string
	rdb      *redis.Client
}

// ServiceMiddleware is a chainable behavior modifier for Service.
type ServiceMiddleware func(Service) Service

func New(animalURLs map[string]string, redisClient *redis.Client) Service {
	return service {
		animalURLs: animalURLs,
		rdb: redisClient,
	}
}

func (s service) Scrape(ctx context.Context, animal string) ([]string, error) {
	var visited []string

	url, ok := s.animalURLs[animal]
	if !ok {
		return visited, ErrAnimalUnsupported
	}

	key := fmt.Sprintf("facts:%s", animal)

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
			tP := e.ChildText("p")
			s.rdb.ZAdd(ctx, key, &redis.Z{0, tP})
		})

		c.OnRequest(func(r *colly.Request) {
			visited = append(visited, fmt.Sprintf("%s", r.URL))
		})

		c.Visit(url)
	default:
	}

	return visited, nil
}
