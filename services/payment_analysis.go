package services

import (
	"github.com/guidogimeno/smartpay-be/types"
	"github.com/shopspring/decimal"
)

const (
	inflationRate = 0.078
	tnaRate       = 0.97
)

func PaymentAnalysis(payment *types.Payment) *types.Analysis {
	installmentWithInterest := calculateInstallmentWithInterest(payment)
	installmentsAtPresentValue := calculateInstallmentsWithInflation(
		installmentWithInterest,
		payment.NumberOfInstallments,
	)
	totalAtPresentValue := sumInstallments(installmentsAtPresentValue)
	savingsFromFixedTermDeposit := calculateSavingsFromFixedTermDeposit(
		payment.Amount,
		installmentWithInterest,
		payment.NumberOfInstallments,
	)

	return &types.Analysis{
		TotalAtPresentValue:         totalAtPresentValue,
		InstallmentsAtPresentValue:  installmentsAtPresentValue,
		SavingsFromFixedTermDeposit: savingsFromFixedTermDeposit,
	}
}

func calculateInstallmentWithInterest(payment *types.Payment) decimal.Decimal {
	numberOfInstallments := decimal.NewFromInt(int64(payment.NumberOfInstallments))
	installmentAmount := payment.Amount.Div(numberOfInstallments)
	interest := installmentAmount.Mul(payment.InterestRate)
	return installmentAmount.Add(interest)
}

func calculateInstallmentsWithInflation(
	installmentAmount decimal.Decimal,
	numOfInstallments int,
) []*types.Installment {
	inflation := decimal.NewFromFloat(inflationRate)

	installments := []*types.Installment{}
	for i := 1; i <= numOfInstallments; i++ {
		period := decimal.NewFromInt(int64(i))
		dividend := inflation.Add(decimal.NewFromInt(1)).Pow(period)
		presentValue := installmentAmount.Div(dividend)

		installment := &types.Installment{
			Number: i,
			Amount: presentValue,
		}

		installments = append(installments, installment)
	}

	return installments
}

func sumInstallments(installments []*types.Installment) decimal.Decimal {
	total := decimal.NewFromFloat(0)

	for _, installment := range installments {
		total = total.Add(installment.Amount)
	}

	return total
}

func calculateSavingsFromFixedTermDeposit(
	paymentAmount decimal.Decimal,
	installmentAmount decimal.Decimal,
	numOfInstallments int,
) decimal.Decimal {
	tna := decimal.NewFromFloat(tnaRate)
	tnm := tna.Div(decimal.NewFromInt(12))

	savings := paymentAmount.Copy()

	for i := 1; i <= numOfInstallments; i++ {
		fixedDepositInterest := savings.Mul(tnm)
		savings = savings.Add(fixedDepositInterest).Sub(installmentAmount)
	}

	return savings
}
