package types

import "github.com/shopspring/decimal"

type FinancialData struct {
	Date  string
	Index float64
}

func (f *FinancialData) GetIndex() decimal.Decimal {
	return decimal.NewFromFloat(f.Index)
}
