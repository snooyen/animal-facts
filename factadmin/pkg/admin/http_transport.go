package admin

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
)

var (
	// Set a Decoder instance as a package global, because it caches
	// meta-data about structs, and an instance can be shared safely.
	schemaDecoder = schema.NewDecoder()
	ErrBadFactID  = errors.New("could not parse fact id to int")
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func NewHTTPHandler(endpoints Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	options := []httptransport.ServerOption{}

	/*
		POST /admin/approve/{id}	create a new fact
		POST /admin/delete/{id}		get a fact by its id
		POST /admin/defer/{id}		get list of known animals
	*/

	r.Methods("POST").Path("/admin/sms").Handler(httptransport.NewServer(
		endpoints.HandleSMSEndpoint,
		decodeHTTPHandleSMSRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeHTTPHandleSMSRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	err = r.ParseForm()
	if err != nil {
		return nil, err
	}

	req := handleSMSRequest{}
	err = schemaDecoder.Decode(&req, r.PostForm)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	default:
		return http.StatusInternalServerError
	}
}
