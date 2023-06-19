package types

import "github.com/shopspring/decimal"

// API
type AnalysisResponse struct {
	TotalAtPresentValue         string                 `json:"totalAtPresentValue"`
	InstallmentsAtPresentValue  []*InstallmentResponse `json:"installmentsAtPresentValue"`
	SavingsFromFixedTermDeposit string                 `json:"savingsFromFixedTermDeposit"`
}

type InstallmentResponse struct {
	Number int    `json:"number"`
	Amount string `json:"amount"`
}

// Domain
type Analysis struct {
	TotalAtPresentValue         decimal.Decimal
	InstallmentsAtPresentValue  []*Installment
	SavingsFromFixedTermDeposit decimal.Decimal
}

type Installment struct {
	Number int
	Amount decimal.Decimal
}

func (a *Analysis) ToAnalysisResponse() *AnalysisResponse {
	installments := []*InstallmentResponse{}

	for _, installment := range a.InstallmentsAtPresentValue {
		installments = append(installments, installment.toInstallmentResponse())
	}

	return &AnalysisResponse{
		TotalAtPresentValue:         a.TotalAtPresentValue.StringFixedBank(2),
		InstallmentsAtPresentValue:  installments,
		SavingsFromFixedTermDeposit: a.SavingsFromFixedTermDeposit.StringFixedBank(2),
	}
}

func (i *Installment) toInstallmentResponse() *InstallmentResponse {
	return &InstallmentResponse{
		Number: i.Number,
		Amount: i.Amount.StringFixedBank(2),
	}
}
