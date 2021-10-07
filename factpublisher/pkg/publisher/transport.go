package publisher

import (
	"context"
	"encoding/json"
	"net/http"
)

type PublishRequest struct {
	Animal string
}

type PublishResponse struct {
	Fact  string  `json:"fact"`
	Score float64 `json:"score"`
	Err   string  `json:"err,omitempty"`
}

func DecodePublishRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request PublishRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
