package types

import "github.com/shopspring/decimal"

type FinancialData struct {
	Date  string
	Index decimal.Decimal
}
