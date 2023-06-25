package db

import "github.com/guidogimeno/smartpay-be/types"

type Storer interface {
	Create(*types.FinancialData) error
	Read() (*types.FinancialData, error)
}
