package usecase

import (
	BalanceRep "github.com/Andronovdima/job-avito-assignment/internal/app/balance/repository"
	"github.com/Andronovdima/job-avito-assignment/internal/models"
	"net/http"
)

type BalanceUsecase struct {
	BalanceRep *BalanceRep.BalanceRepository
}

func NewBalanceUsecase(c *BalanceRep.BalanceRepository) *BalanceUsecase {
	BalanceUsecase := &BalanceUsecase{
		BalanceRep: c,
	}
	return BalanceUsecase
}

func (c *BalanceUsecase) GetBalance(user *models.User) (*models.Balance, error) {
	err := new(models.HttpError)

	balance, cerr := c.BalanceRep.GetBalance(user)
	if cerr != nil {
		err.StatusCode = http.StatusBadRequest
		err.StringErr = "No information about this account yet"
		return nil, err
	}

	return balance, nil
}

func (c *BalanceUsecase) ChangeBalance(operation *models.Balance) error {
	err := new(models.HttpError)

	isExistBill := c.BalanceRep.IsExistBill(operation.UserID)
	if !isExistBill {
		_ = c.BalanceRep.CreateBill(operation.UserID)
	}

	cerr := c.BalanceRep.ChangeBalance(operation)
	if cerr != nil {
		err = cerr.(*models.HttpError)
		return err
	}

	return nil
}

func (c *BalanceUsecase) Transaction(transaction *models.Transaction) error {
	err := new(models.HttpError)

	isExistBill := c.BalanceRep.IsExistBill(transaction.Sender)
	if !isExistBill {
		err.StatusCode = http.StatusOK
		err.StringErr = "Sender account doesn't have enough money"
		return err
	}

	isExistBill = c.BalanceRep.IsExistBill(transaction.Receiver)
	if !isExistBill {
		_ = c.BalanceRep.CreateBill(transaction.Receiver)
	}

	cerr := c.BalanceRep.Transaction(transaction)
	if cerr != nil {
		err = cerr.(*models.HttpError)
		return err
	}

	return nil
}
