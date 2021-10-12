package scraper

import (
	"context"
	"encoding/json"
	"net/http"
)

func DecodeScrapeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request scrapeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
