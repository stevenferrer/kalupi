package ledger

import (
	"github.com/stevenferrer/kalupi/currency"
)

// List of cash ledger account numbers
const (
	CashUSDLedgerNo LedgerNo = "100"
	// CashEURLedgerNo LedgerNo = "110"
)

// cashLedgers contains the list of cash ledgers
var cashLedgers = [...]Ledger{
	cashUSD,
	// cashEUR,
}

// List of cash ledgers
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

//GetCashLedgerNo retrieives the cash ledger number for the given currency
func GetCashLedgerNo(curr currency.Currency) (LedgerNo, error) {
	switch curr {
	case currency.USD:
		return CashUSDLedgerNo, nil
	}

	return "", currency.ErrUnsupportedCurrency
}
