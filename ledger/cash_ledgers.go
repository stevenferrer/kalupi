package ledger

import (
	"github.com/pkg/errors"

	"github.com/sf9v/kalupi/currency"
)

// Ledger account numbers
const (
	CashUSDLedgerNo LedgerNo = "100"
	// CashEURLedgerNo LedgerNo = "110"
)

var cashLedgers = [...]Ledger{
	cashUSD,
	// cashEUR,
}

var (
	cashUSD = Ledger{
		LedgerNo:    CashUSDLedgerNo,
		AccountType: AccountTypeLiability,
		Currency:    currency.USD,
		Name:        "Cash USD",
	}

	// cashEUR = Ledger{
	// 	LedgerNo:    CashEURLedgerNo,
	// 	AccountType: AccountTypeLiability,
	// 	Currency:    currency.USD,
	// 	Name:        "Cash EUR",
	// }
)

func GetCashLedgerNo(curr currency.Currency) (LedgerNo, error) {
	switch curr {
	case currency.USD:
		return CashUSDLedgerNo, nil
	}

	return "", errors.Errorf("currency %s is not supported", curr)
}
