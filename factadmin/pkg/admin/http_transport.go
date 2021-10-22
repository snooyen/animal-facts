package admin

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
	ErrBadFactID  = errors.New("could not parse fact id to int")
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

	r.Methods("POST").Path("/admin/approve/{id}").Handler(httptransport.NewServer(
		endpoints.ApproveFactEndpoint,
		decodeHTTPApproveFactRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/admin/delete/{id}").Handler(httptransport.NewServer(
		endpoints.DeleteFactEndpoint,
		decodeHTTPDeleteFactRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/admin/defer/{id}").Handler(httptransport.NewServer(
		endpoints.DeferFactEndpoint,
		decodeHTTPDeferFactRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeHTTPApproveFactRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return struct{}{}, nil
}

func decodeHTTPDeferFactRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return struct{}{}, nil
}

func decodeHTTPHandleSMSRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	req := handleSMSRequest{}
	e := json.NewDecoder(r.Body).Decode(&req)
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
	default:
		return http.StatusInternalServerError
	}
}
