package account

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewHandler(s Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{}

	createAccountHandler := kithttp.NewServer(
		newCreateAccountEndpoint(s),
		decodeCreateAccountRequest,
		encodeResponse,
		opts...,
	)

	getAccountHandler := kithttp.NewServer(
		newGetAccountEndpoint(s),
		decodeGetAccountRequest,
		encodeResponse,
		opts...,
	)

	listAccountsHandler := kithttp.NewServer(
		newListAccountsEndpoint(s),
		decodeListAccountsRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/accounts", createAccountHandler).Methods(http.MethodPost)
	r.Handle("/accounts/{accountId}", getAccountHandler).Methods(http.MethodGet)
	r.Handle("/accounts", listAccountsHandler).Methods(http.MethodGet)

	return r
}

func decodeCreateAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func decodeGetAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)

	accountID, ok := vars["accountId"]
	if !ok {
		return nil, errors.New("bad route")
	}

	return getAccountRequest{AccountID: AccountID(accountID)}, nil
}

func decodeListAccountsRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return listAccountsRequest{}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
}
