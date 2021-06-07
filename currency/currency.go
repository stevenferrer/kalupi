package currency

import (
	"database/sql/driver"
	"errors"
)

// Currency is a currency
type Currency int

const (
	USD Currency = iota + 1
	// EUR
)

func (c Currency) String() string {
	return [...]string{
		"invalid",
		"USD",
		// "EUR",
	}[c]
}

// IsValid returns true if currency is valid
func (c Currency) IsValid() bool {
	return c != Currency(0)
}

func (c Currency) Value() (driver.Value, error) {
	return c.String(), nil
}

func (c *Currency) Scan(src interface{}) error {
	if src == nil {
		*c = Currency(0)
		return nil
	}

	val, ok := src.(string)
	if !ok {
		return errors.New("src is not string")
	}

	*c = strToCurrency(val)
	return nil
}

func strToCurrency(s string) Currency {
	switch s {
	case "USD":
		return USD
		// case "EUR":
		// 	return EUR
	}

	return Currency(0)
}

// func (c Currency) MarshalJSON() ([]byte, error) {
// 	buf := bytes.NewBuffer([]byte(`"`))
// 	buf.WriteString(c.String())
// 	buf.WriteByte(byte('"'))
// 	return buf.Bytes(), nil
// }

// func (c *Currency) UnmarshalJSON(data []byte) error {
// 	var s string
// 	err := json.Unmarshal(data, &s)
// 	if err != nil {
// 		return err
// 	}

// 	*c = strToCurrency(s)
// 	return nil
// }
