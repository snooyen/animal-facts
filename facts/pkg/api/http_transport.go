package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadFactID  = errors.New("could not parse fact id to int")
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func NewHTTPHandler(endpoints Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	options := []httptransport.ServerOption{}

	/*
		POST /facts/	create a new fact
		GET /facts/{id} get a fact by its id
		GET	/animals/	get list of known animals
	*/

	r.Methods("POST").Path("/facts/").Handler(httptransport.NewServer(
		endpoints.CreateFactEndpoint,
		decodeHTTPCreateFactRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/facts/{id}").Handler(httptransport.NewServer(
		endpoints.GetFactEndpoint,
		decodeHTTPGetFactRequest,
		encodeResponse,
		options...,
	))

	r.Methods("DELETE").Path("/facts/{id}").Handler(httptransport.NewServer(
		endpoints.DeleteFactEndpoint,
		decodeHTTPDeleteFactRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/animals/").Handler(httptransport.NewServer(
		endpoints.GetAnimalsEndpoint,
		decodeHTTPGetAnimalsRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/random/{animal}").Handler(httptransport.NewServer(
		endpoints.GetRandAnimalFactEndpoint,
		decodeHTTPGetRandAnimalFactRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeHTTPCreateFactRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req createFactRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeHTTPGetAnimalsRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return nil, nil
}

func decodeHTTPGetFactRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	id, err := decodeRequestID(r)
	if err != nil {
		return nil, err
	}

	req := getFactRequest{ID: id}
	e := json.NewDecoder(r.Body).Decode(&request)
	if e == nil {
		return nil, e
	}
	return req, nil
}

func decodeHTTPDeleteFactRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	id, err := decodeRequestID(r)
	if err != nil {
		return nil, err
	}

	req := deleteFactRequest{ID: id}
	e := json.NewDecoder(r.Body).Decode(&request)
	if e == nil {
		return nil, e
	}
	return req, nil
}

func decodeHTTPGetRandAnimalFactRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	animal, ok := vars["animal"]
	if !ok {
		return -1, ErrBadRouting
	}

	req := getRandAnimalFactRequest{Animal: animal}
	e := json.NewDecoder(r.Body).Decode(&request)
	if e == nil {
		return nil, e
	}
	return req, nil
}

func decodeRequestID(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	idVar, ok := vars["id"]
	if !ok {
		return -1, ErrBadRouting
	}

	id, err := strconv.Atoi(idVar)
	if err != nil {
		return -1, ErrBadFactID
	}
	return int64(id), nil
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
	case ErrNotFound:
		return http.StatusNotFound
	case ErrAlreadyExists:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
