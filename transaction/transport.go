package transaction

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

	depositHandler := kithttp.NewServer(
		newDepositEndpoint(s),
		decodeDepositRequest,
		encodeResponse,
		opts...,
	)

	withdrawHandler := kithttp.NewServer(
		newWithdrawalEndpoint(s),
		decodeWithdrawalRequest,
		encodeResponse,
		opts...,
	)

	paymentHandler := kithttp.NewServer(
		newPaymentEndpoint(s),
		decodePaymentRequest,
		encodeResponse,
		opts...,
	)

	mux := chi.NewMux()

	mux.Method(http.MethodPost, "/deposit", depositHandler)
	mux.Method(http.MethodPost, "/withdraw", withdrawHandler)
	mux.Method(http.MethodPost, "/payment", paymentHandler)

	return mux
}

func decodeDepositRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request depositRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func decodeWithdrawalRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request withdrawalRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func decodePaymentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request paymentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
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

	if errors.Is(err, ErrValidation) {
		w.WriteHeader(http.StatusBadRequest)
	} else if errors.Is(err, ErrSendingAccountNotFound) ||
		errors.Is(err, ErrReceivingAccountNotFound) {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
}
