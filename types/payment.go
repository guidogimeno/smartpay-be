package types

import (
	"errors"

	"github.com/shopspring/decimal"
)

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

func (p *PaymentRequest) ToPayment() (*Payment, error) {
	err := p.isValid()
	if err != nil {
		return nil, err
	}

	return &Payment{
		Amount:               decimal.NewFromFloat(p.Amount),
		NumberOfInstallments: p.NumberOfInstallments,
		InterestRate:         decimal.NewFromFloat32(p.InterestRate),
	}, nil
}

func (p *PaymentRequest) isValid() error {
	if p.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	if p.InterestRate < 0 {
		return errors.New("interest rate must be positive")
	}

	if p.NumberOfInstallments < 1 {
		return errors.New("number of installments must be positive")
	}

	return nil
}
