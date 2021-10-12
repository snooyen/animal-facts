package scraper

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func MakeScrapeEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(scrapeRequest)
		v, err := s.Scrape(ctx, req.Animal)
		if err != nil {
			return scrapeResponse{v, err.Error()}, nil
		}
		return scrapeResponse{v, ""}, nil
	}
}

type scrapeRequest struct {
	Animal string
}

type scrapeResponse struct {
	Visited []string `json:visited"`
	Err     string   `json:"err,omitempty"`
}
