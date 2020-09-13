package repository

import (
	"database/sql"
	"github.com/Andronovdima/job-avito-assignment/internal/models"
	"net/http"
)

type BalanceRepository struct {
	db *sql.DB
}

func NewBalanceRepository(thisDB *sql.DB) *BalanceRepository {
	balanceRep := &BalanceRepository{
		db: thisDB,
	}
	return balanceRep
}

func (b *BalanceRepository) GetBalance(user *models.User) (*models.Balance, error) {
	balance := new(models.Balance)

	err := b.db.QueryRow(
		"SELECT id, user_id, balance " +
			"FROM balances "+
			"WHERE user_id = $1",
		user.ID,
	).Scan(&balance.ID, &balance.UserID, &balance.Sum)

	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (b *BalanceRepository) IsExistBill(userID int64) bool {
	row := b.db.QueryRow(
		"SELECT user_id "+
			"FROM balances "+
			"WHERE user_id = $1",
		userID,
	)
	if row.Scan(&userID) != nil {
		return false
	}

	return true
}

func (b *BalanceRepository) CreateBill(userID int64) error {
	var id int64
	err := b.db.QueryRow(
		"INSERT INTO balances (user_id, balance) "+
			"VALUES ($1 , $2) " +
			"RETURNING id ",
		userID,
		0,
	).Scan(&id)

	if err != nil {
		return err
	}

	return nil
}


func (b *BalanceRepository) ChangeBalance(operation *models.Balance) error {
	HttpError := new(models.HttpError)
	tx, err := b.db.Begin()
	if err != nil {
		HttpError.StringErr = "Transaction error"
		HttpError.StatusCode = http.StatusInternalServerError
		return HttpError
	}
	var balance int64

	err = b.db.QueryRow(
		"SELECT balance " +
			"FROM balances "+
			"WHERE user_id = $1 ",
		operation.UserID,
	).Scan(&balance)

	if err != nil {
		HttpError.StringErr = "No bill for this user_id"
		HttpError.StatusCode = http.StatusInternalServerError
		return HttpError
	}

	if balance + operation.Sum < 0 {
		_ = tx.Commit()
		HttpError.StringErr = "Insufficient funds in the account"
		HttpError.StatusCode = http.StatusForbidden
		return HttpError
	}

	_ , err = b.db.Exec(
		"UPDATE balances " +
			"SET balance = balance + $1 " +
			"WHERE user_id = $2 AND ((balance + $1) > -1)",
		operation.Sum,
		operation.UserID,
	)
	if err != nil {
		HttpError.StringErr = "Transaction error"
		HttpError.StatusCode = http.StatusInternalServerError
		_ = tx.Rollback()
		return HttpError
	}

	err = tx.Commit()
	if err != nil {
		HttpError.StringErr = "Transaction error"
		HttpError.StatusCode = http.StatusInternalServerError
		return HttpError
	}

	return nil
}

func (b *BalanceRepository) Transaction (transaction *models.Transaction) error {
	HttpError := new(models.HttpError)
	tx, err := b.db.Begin()
	if err != nil {
		HttpError.StringErr = "Transaction Begin error"
		HttpError.StatusCode = http.StatusInternalServerError
		return HttpError
	}
	var balance int64

	err = b.db.QueryRow(
		"SELECT balance " +
			"FROM balances "+
			"WHERE user_id = $1 ",
		transaction.Sender,
	).Scan(&balance)

	if err != nil {
		HttpError.StringErr = "No bill for this user_id"
		HttpError.StatusCode = http.StatusInternalServerError
		return HttpError
	}

	if balance - transaction.Sum < 0 {
		_ = tx.Commit()
		HttpError.StringErr = "Insufficient funds in the sender account"
		HttpError.StatusCode = http.StatusForbidden
		return HttpError
	}

	_ , err = b.db.Exec(
		"UPDATE balances " +
			"SET balance = balance - $1 " +
			"WHERE user_id = $2 AND ((balance - $1) > -1)",
		transaction.Sum,
		transaction.Sender,
	)

	if err != nil {
		HttpError.StringErr = "Transaction Exec error"
		HttpError.StatusCode = http.StatusInternalServerError
		_ = tx.Rollback()
		return HttpError
	}

	_ , err = b.db.Exec(
		"UPDATE balances " +
			"SET balance = balance + $1 " +
			"WHERE user_id = $2",
		transaction.Sum,
		transaction.Receiver,
	)

	if err != nil {
		HttpError.StringErr = "Transaction Exec error"
		HttpError.StatusCode = http.StatusInternalServerError
		_ = tx.Rollback()
		return HttpError
	}

	err = tx.Commit()
	if err != nil {
		HttpError.StringErr = "Transaction Commit error"
		HttpError.StatusCode = http.StatusInternalServerError
		return HttpError
	}

	return nil
}