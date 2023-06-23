package services

import (
	"github.com/guidogimeno/smartpay-be/scrapper"
	"github.com/guidogimeno/smartpay-be/types"
	"github.com/guidogimeno/smartpay-be/utils"
	"github.com/shopspring/decimal"
)

func PaymentAnalysis(payment *types.Payment) (*types.Analysis, error) {
	installmentWithInterest := calculateInstallmentWithInterest(payment)

	installmentsAtPresentValue, err := calculateInstallmentsWithInflation(
		installmentWithInterest,
		payment.NumberOfInstallments,
	)
	if err != nil {
		return nil, err
	}

	totalAtPresentValue := sumInstallments(installmentsAtPresentValue)

	savingsFromFixedTermDeposit, err := calculateSavingsFromFixedTermDeposit(
		payment.Amount,
		installmentWithInterest,
		payment.NumberOfInstallments,
	)
	if err != nil {
		return nil, err
	}

	return &types.Analysis{
		TotalAtPresentValue:         totalAtPresentValue,
		InstallmentsAtPresentValue:  installmentsAtPresentValue,
		SavingsFromFixedTermDeposit: savingsFromFixedTermDeposit,
	}, nil
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
) ([]*types.Installment, error) {
	startDate := utils.NewClock().AddMonths(-2).Format()
	finishDate := utils.NewClock().AddMonths(-1).Format()

	financialData, err := scrapper.ScrapInflation(startDate, finishDate)
	if err != nil {
		return nil, err
	}

	inflation := financialData[0].Index

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

	return installments, nil
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
) (decimal.Decimal, error) {
	yesterday := utils.NewClock().AddDays(-1).Format()

	financialData, err := scrapper.ScrapTNA(yesterday, yesterday)
	if err != nil {
		return decimal.NewFromInt(0), err
	}

	tna := financialData[0].Index
	tnm := tna.Div(decimal.NewFromInt(12))

	savings := paymentAmount.Copy()

	for i := 1; i <= numOfInstallments; i++ {
		fixedDepositInterest := savings.Mul(tnm)
		savings = savings.Add(fixedDepositInterest).Sub(installmentAmount)
	}

	return savings, nil
}
