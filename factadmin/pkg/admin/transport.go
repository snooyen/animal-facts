package admin

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

type ApprovalRequest struct {
	Action string
	Animal string
	Fact   string
}

type ApprovalResponse struct {
	Action string `json:"action"`
	Animal string `json:"animal"`
	Fact   string `json:"fact"`
	Msg    string `json:"msg,omitempty"`
	Err    string `json:"err,omitempty"`
}

func NewHTTPHandler(s Service) http.Handler {
	m := http.NewServeMux()

	m.Handle("/approve", httptransport.NewServer(
		MakeApproveEndpoint(s),
		decodeApproveRequest,
		encodeResponse,
	))

	return m
}

func decodeApproveRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request ApprovalRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
