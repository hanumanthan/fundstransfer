package handlers

import (
	"fmt"
	"fundstransfer/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Transact(c *gin.Context) {
	var t Transaction
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId, _ := strconv.Atoi(c.Param("user_id"))
	if err := t.CreateTransaction(userId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "Transaction processed."})
}

func (t *Transaction) CreateTransaction(userId int) error {
	var from, to models.Wallet
	if err := from.GetWalletForUser(userId); err != nil {
		return err
	}
	if err := to.GetWalletForMobileNumber(t.To); err != nil {
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

func GetTransactions(c *gin.Context) {
	transactions, err := models.GetAllTransactions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, gin.H{"transactions": &transactions})
}
