package types

import (
	"errors"

	"github.com/shopspring/decimal"
)

type Payment struct {
	Amount               float64 `json:"amount"`
	NumberOfInstallments int     `json:"numberOfInstallments"`
	InterestRate         float32 `json:"interestRate"`
}

func (p *Payment) GetAmount() decimal.Decimal {
	return decimal.NewFromFloat(p.Amount)
}

func (p *Payment) GeInterestRate() decimal.Decimal {
	return decimal.NewFromFloat32(p.InterestRate)
}

func (p *Payment) IsValid() error {
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
