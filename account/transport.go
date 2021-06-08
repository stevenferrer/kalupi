package account

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
)

func NewHandler(s Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

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

	mux := chi.NewMux()

	mux.Method(http.MethodPost, "/", createAccountHandler)
	mux.Method(http.MethodGet, "/", listAccountsHandler)
	mux.Method(http.MethodGet, "/{id}", getAccountHandler)

	return mux
}

var (
	errBadRoute = errors.New("bad route")
)

func decodeCreateAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func decodeGetAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	accountID := chi.URLParam(r, "id")
	if accountID == "" {
		return nil, errBadRoute
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

	// TODO: use 422 for validation errors
	if errors.Is(err, ErrValidation) {
		w.WriteHeader(http.StatusBadRequest)
	} else if errors.Is(err, ErrAccountNotFound) {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
}
