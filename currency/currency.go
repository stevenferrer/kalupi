package currency

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Currency is a currency
type Currency int

// List of supported currencies
const (
	USD Currency = iota + 1
	// EUR
)

// String implements Stringer
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

// Value implements driver.Valuer interface
func (c Currency) Value() (driver.Value, error) {
	return c.String(), nil
}

// Scan implements the sql.Scanner interface
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

// strToCurrency takes a string and returns the Currency
func strToCurrency(s string) Currency {
	switch s {
	case "USD":
		return USD
		// case "EUR":
		// 	return EUR
	}

	return Currency(0)
}

// MarshalJSON implements the json.Marshaler interface
func (c Currency) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer([]byte(`"`))
	buf.WriteString(c.String())
	buf.WriteByte(byte('"'))
	return buf.Bytes(), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (c *Currency) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	*c = strToCurrency(s)
	return nil
}
