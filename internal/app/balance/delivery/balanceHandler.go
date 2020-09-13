package delivery

import (
	"encoding/json"
	BalanceUC "github.com/Andronovdima/job-avito-assignment/internal/app/balance/usecase"
	"github.com/Andronovdima/job-avito-assignment/internal/app/respond"
	"github.com/Andronovdima/job-avito-assignment/internal/models"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
)

type BalanceHandler struct {
	BalanceUsecase BalanceUC.BalanceUsecase
	logger         *zap.SugaredLogger
}

func NewBalanceHandler(m *mux.Router, balance BalanceUC.BalanceUsecase, logger *zap.SugaredLogger) {
	handler := &BalanceHandler{
		BalanceUsecase: balance,
		logger:         logger,
	}

	m.HandleFunc("/balance/change", handler.HandleChangeBalance).Methods(http.MethodPost)
	m.HandleFunc("/balance/get", handler.HandleGetBalance).Methods(http.MethodGet)
	m.HandleFunc("/balance/transaction", handler.HandleTransaction).Methods(http.MethodPost)

}

func (c *BalanceHandler) HandleGetBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "ERROR : HandleGetBalance <- Body.Close")
			c.logger.Info(err.Error())
			respond.Error(w, r, http.StatusInternalServerError, err)
			return
		}
	}()

	thisUser := new(models.User)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(thisUser)
	if err != nil {
		err = errors.Wrapf(err, "ERROR : HandleGetBalance <- Decode: ")
		c.logger.Info(err.Error())
		err = errors.New("Invalid format data, example: '{'user_id': 1 }'")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	if thisUser.ID <= 0 {
		err = errors.New("Invalid format data, example: '{'user_id': 1 }'")
		c.logger.Info(err.Error())
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	balance, err := c.BalanceUsecase.GetBalance(thisUser)
	if err != nil {
		rerr := err.(*models.HttpError)
		c.logger.Info("ERROR : HandleGetBalance <- GetBalance: " + rerr.Error())
		respond.Error(w, r, rerr.StatusCode, rerr)
		return
	}

	c.logger.Info("/balance/get || HTTP: 200")
	respond.Respond(w, r, http.StatusOK, balance)
	return
}

func (c *BalanceHandler) HandleChangeBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "ERROR : HandleChangeBalance <- Body.Close")
			c.logger.Info(err.Error())
			respond.Error(w, r, http.StatusInternalServerError, err)
			return
		}
	}()

	operation := new(models.Balance)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(operation)
	if err != nil {
		err = errors.Wrapf(err, "ERROR : HandleChangeBalance <- Decode: ")
		c.logger.Info(err.Error())
		err = errors.New("Invalid format data, example: '{'user_id': 1, 'sum' : 100 }' , note: sum can't be 0")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	if operation.UserID <= 0 || operation.Sum == 0 {
		err = errors.New("Invalid format data, example: '{'user_id': 1, 'sum' : 100 }', note: sum can't be 0")
		c.logger.Info(err.Error())
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	err = c.BalanceUsecase.ChangeBalance(operation)
	if err != nil {
		rerr := err.(*models.HttpError)
		c.logger.Info("ERROR : HandleChangeBalance <- ChangeBalance: " + rerr.Error())
		respond.Error(w, r, rerr.StatusCode, rerr)
		return
	}

	c.logger.Info("/balance/change || HTTP: 200")
	respond.Respond(w, r, http.StatusOK, nil)
	return
}

func (c *BalanceHandler) HandleTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "ERROR : HandleTransaction <- Body.Close")
			c.logger.Info(err.Error())
			respond.Error(w, r, http.StatusInternalServerError, err)
			return
		}
	}()

	transaction := new(models.Transaction)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(transaction)
	if err != nil {
		err = errors.Wrapf(err, "ERROR : HandleTransaction <- Decode: ")
		c.logger.Info(err.Error())
		err = errors.New("Invalid format data, example: '{'sender': 1, 'receiver' : 2, 'sum': 100 }' , note: sum can't be <1")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	if transaction.Sender <= 0 || transaction.Receiver <= 0 || transaction.Sum <= 0 {
		err = errors.New("Invalid format data, example: '{'sender': 1, 'receiver' : 2, 'sum': 100 }' , note: sum can't be <1")
		c.logger.Info(err.Error())
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	err = c.BalanceUsecase.Transaction(transaction)
	if err != nil {
		rerr := err.(*models.HttpError)
		c.logger.Info("ERROR : HandleTransaction <- Transaction: " + rerr.Error())
		respond.Error(w, r, rerr.StatusCode, rerr)
		return
	}

	c.logger.Info("/balance/transaction || HTTP: 200")
	respond.Respond(w, r, http.StatusOK, nil)
	return
}
