package handlers

import (
	"fmt"
	"fundstransfer/pkg/logger"
	"fundstransfer/pkg/metrics"
	"fundstransfer/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"time"
)
// Transact - Extracts the wallet of destination mobile number and
// transfers money from users wallet to destination wallet
// Error validation for wallets and balance
func Transact(c *gin.Context) {
	var t Transaction
	var from, to models.Wallet
	userId, _ := strconv.Atoi(c.Param("user_id"))
	if err := c.ShouldBindJSON(&t); err != nil {
		logger.ERROR.Println("input transaction object is invalid")
		metrics.CaptureErrorMetrics(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := from.GetWalletForUser(userId); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			metrics.CaptureErrorMetrics(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		} else {
			metrics.CaptureErrorMetrics(http.StatusInternalServerError)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	if err := to.GetWalletForMobileNumber(t.To); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			metrics.CaptureErrorMetrics(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wallet mobile number"})
		} else {
			metrics.CaptureErrorMetrics(http.StatusInternalServerError)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	if err := t.CreateTransaction(from, to); err != nil {
		metrics.CaptureErrorMetrics(http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "Transaction processed."})
}

func (t *Transaction) CreateTransaction(from, to models.Wallet) error {
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

// Gets all the transaction in the system
func GetTransactions(c *gin.Context) {
	transactions, err := models.GetAllTransactions()
	if err != nil {
		metrics.CaptureErrorMetrics(http.StatusInternalServerError)
		logger.ERROR.Println(fmt.Sprintf("error getting transactions - %v", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, &Transactions{Transactions:transactions} )
}