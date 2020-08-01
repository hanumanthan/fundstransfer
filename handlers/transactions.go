package handlers

import (
	"fmt"
	"fundstransfer/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Transact(c *gin.Context) {
	var createTransaction CreateTransaction
	if err := c.ShouldBindJSON(&createTransaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := createTransaction.CreateTransaction(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.String(http.StatusOK, "Transaction processed")
}

func (t *CreateTransaction) CreateTransaction() error {
	var from, to models.Wallet
	if err := from.GetWalletForUser(t.From); err != nil {
		return err
	}
	if err := to.GetWalletForUser(t.To); err != nil {
		return err
	}
	if from.Balance < t.Amount {
		return fmt.Errorf("insufficient balance to transfer")
	}
	transaction := models.Transaction{
		FromWallet: from.ID,
		ToWallet:   to.ID,
		Amount:     t.Amount,
		Message:    t.Message}

	from.Balance -= t.Amount
	to.Balance += t.Amount
	if err := transaction.Save(); err != nil {
		return fmt.Errorf("error saving transaction %v", err.Error())
	}
	if err := from.Save(); err != nil {
		return fmt.Errorf("error updating wallet %v", err.Error())
	}
	if err := to.Save(); err != nil {
		return fmt.Errorf("error updating wallet %v", err.Error())
	}
	return nil
}
