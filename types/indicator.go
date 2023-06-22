package types

import "github.com/shopspring/decimal"

type Inflation struct {
	Date  string
	Index decimal.Decimal
}
