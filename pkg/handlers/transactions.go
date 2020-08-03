package handlers

import (
	"fmt"
	"fundstransfer/pkg/logger"
	"fundstransfer/pkg/metrics"
	"fundstransfer/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func Transact(c *gin.Context) {
	var t Transaction
	if err := c.ShouldBindJSON(&t); err != nil {
		logger.ERROR.Println("input transaction object is invalid")
		metrics.CaptureErrorMetrics(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId, _ := strconv.Atoi(c.Param("user_id"))
	if err := t.CreateTransaction(userId); err != nil {
		metrics.CaptureErrorMetrics(http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "Transaction processed."})
}

func (t *Transaction) CreateTransaction(userId int) error {
	var from, to models.Wallet
	if err := from.GetWalletForUser(userId); err != nil {
		logger.ERROR.Println(fmt.Sprintf("couldnt retrieve wallet for user %v", userId))
		return err
	}
	if err := to.GetWalletForMobileNumber(t.To); err != nil {
		logger.ERROR.Println(fmt.Sprintf("couldnt retrieve wallet for mobile %v", t.To))
		return err
	}
	if from.Balance < t.Amount {
		logger.ERROR.Println("insufficient balance to transfer")
		return fmt.Errorf("insufficient balance to transfer")
	}
	transaction := models.Transaction{
		FromWallet: from.ID,
		ToWallet:   to.ID,
		Amount:     t.Amount,
		CreatedAt:  time.Now(),
		Message:    t.Message}

	from.Balance -= t.Amount
	to.Balance += t.Amount
	if err := transaction.Save(); err != nil {
		logger.ERROR.Println(fmt.Sprintf("error saving transaction %v", err.Error()))
		return fmt.Errorf("error saving transaction %v", err.Error())
	}
	if err := from.Save(); err != nil {
		logger.ERROR.Println(fmt.Sprintf("error updating wallet id %d error - %v", from.ID, err.Error()))
		return fmt.Errorf("error updating wallet %v", err.Error())
	}
	if err := to.Save(); err != nil {
		logger.ERROR.Println(fmt.Sprintf("error updating wallet id %d error - %v", to.ID, err.Error()))
		return fmt.Errorf("error updating wallet %v", err.Error())
	}
	return nil
}

func GetTransactions(c *gin.Context) {
	transactions, err := models.GetAllTransactions()
	if err != nil {
		metrics.CaptureErrorMetrics(http.StatusInternalServerError)
		logger.ERROR.Println(fmt.Sprintf("error getting transactions - %v", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, gin.H{"transactions": &transactions})
}
