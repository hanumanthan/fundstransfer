package service

import (
	"fmt"
	"fundstransfer/models"
)

type CreateTransaction struct {
	From    uint    `json:"from" binding:"required"`
	To      uint    `json:"to" binding:"required"`
	Amount  int    `json:"amount" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type GetTransactions struct {
	Transactions []models.Transaction
}

func (t *CreateTransaction) CreateTransaction() error {
	var from, to models.Account
	if err := from.GetAccountByUser(t.From); err != nil {
		return err
	}
	if err := to.GetAccountByUser(t.To); err != nil {
		return err
	}
	if from.Balance < t.Amount {
		return fmt.Errorf("insufficient balance to transfer")
	}
	transaction := models.Transaction{
		FromAccountId: from.ID,
		ToAccountId:   to.ID,
		Amount: t.Amount,
		Message:     t.Message}

	from.Balance -= t.Amount
	to.Balance += t.Amount
	if err := transaction.Save(); err != nil {
		return fmt.Errorf("error saving transaction %v", err.Error())
	}
	if err := from.Save(); err != nil {
		return fmt.Errorf("error updating account %v", err.Error())
	}
	if err := to.Save(); err != nil {
		return fmt.Errorf("error updating account %v", err.Error())
	}
	return nil
}
