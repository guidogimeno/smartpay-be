package types

import "github.com/shopspring/decimal"

// API
type PaymentRequest struct {
	Amount               float64 `json:"amount"`
	NumberOfInstallments int     `json:"numberOfInstallments"`
	InterestRate         float32 `json:"interestRate"`
}

// Domain
type Payment struct {
	Amount               decimal.Decimal
	NumberOfInstallments int
	InterestRate         decimal.Decimal
}

func (p *PaymentRequest) ToPayment() *Payment {
	return &Payment{
		Amount:               decimal.NewFromFloat(p.Amount),
		NumberOfInstallments: p.NumberOfInstallments,
		InterestRate:         decimal.NewFromFloat32(p.InterestRate),
	}
}
