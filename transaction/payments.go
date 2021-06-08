package transaction

import (
	"github.com/sf9v/kalupi/account"
	"github.com/shopspring/decimal"
)

// Payment is maps to transfer transaction
type Payment struct {
	XactNo    XactNo            `json:"xact_no"`
	Account   account.AccountID `json:"account"`
	Amount    decimal.Decimal   `json:"amount"`
	Direction string            `json:"direction"`

	ToAccount   account.AccountID `json:"to_account,omitempty"`
	FromAccount account.AccountID `json:"from_account,omitempty"`
}

func xactsToPayments(xacts []*Transaction) []*Payment {
	// NOTE: Send and receive always comes in pair
	payments := []*Payment{}
	for i := 0; i < len(xacts); i += 2 {
		sndXact := xacts[i]
		rcvXact := xacts[i+1]

		payments = append(payments, &Payment{
			XactNo:    sndXact.XactNo,
			Account:   sndXact.AccountID,
			Amount:    sndXact.Amount,
			ToAccount: rcvXact.AccountID,
			Direction: "outgoing",
		})

		payments = append(payments, &Payment{
			XactNo:      rcvXact.XactNo,
			Account:     rcvXact.AccountID,
			Amount:      rcvXact.Amount,
			FromAccount: sndXact.AccountID,
			Direction:   "incoming",
		})
	}

	return payments
}
