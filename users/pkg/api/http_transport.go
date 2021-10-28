package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadUserID  = errors.New("could not parse user id to int")
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func NewHTTPHandler(endpoints Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	options := []httptransport.ServerOption{}

	/*
		POST /users/	create a new user
		GET /users/{id} get a user by its id
		GET	/animals/	get list of known animals
	*/

	r.Methods("POST").Path("/users/").Handler(httptransport.NewServer(
		endpoints.CreateUserEndpoint,
		decodeHTTPCreateUserRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/users/{id}").Handler(httptransport.NewServer(
		endpoints.GetUserEndpoint,
		decodeHTTPGetUserRequest,
		encodeResponse,
		options...,
	))

	r.Methods("DELETE").Path("/users/{id}").Handler(httptransport.NewServer(
		endpoints.DeleteUserEndpoint,
		decodeHTTPDeleteUserRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeHTTPCreateUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req createUserRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeHTTPGetUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	id, err := decodeRequestID(r)
	if err != nil {
		return nil, err
	}

	req := getUserRequest{ID: id}
	e := json.NewDecoder(r.Body).Decode(&request)
	if e == nil {
		return nil, e
	}
	return req, nil
}

func decodeHTTPDeleteUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	id, err := decodeRequestID(r)
	if err != nil {
		return nil, err
	}

	req := deleteUserRequest{ID: id}
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
		return -1, ErrBadUserID
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
