package handlers

import (
	"fmt"
	"fundstransfer/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetUsers(c *gin.Context) {
	users, err := models.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": &users})
}

func GetWallets(c *gin.Context) {
	wallets, err := models.GetAllWallets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, gin.H{"wallets": &wallets})
}

func GetUserDetails(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("user_id"))
	userDetails, err := extractUserDetails(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": userDetails})
}

func extractUserDetails(userId int) (*UserDetails, error) {
	var user models.User
	var wallet models.Wallet
	if err := user.GetById(userId); err != nil {
		return nil, fmt.Errorf("error retrieving user %v", err)
	}
	if err := wallet.GetWalletForUser(user.ID); err != nil {
		return nil, fmt.Errorf("error retrieving wallet %v", err)
	}
	debits, err := getTransactionsForUser(wallet, true)
	if err != nil {
		return nil, fmt.Errorf("error retrieving debits %v", err)
	}

	credits, err := getTransactionsForUser(wallet, false)
	if err != nil {
		return nil, fmt.Errorf("error retrieving credits %v", err)
	}

	userDetails := &UserDetails{
		Id:                   user.ID,
		Name:                 user.Name,
		Balance:              wallet.Balance,
		SentTransactions:     debits,
		ReceivedTransactions: credits,
	}
	return userDetails, nil
}

func getTransactionsForUser(wallet models.Wallet, isSender bool) ([]TransactionDetails, error) {
	transactions, err := models.GetTransactionsByAccount(wallet.ID, isSender)
	if err != nil {
		return nil, err
	}
	transactionDetails := make([]TransactionDetails, 0)
	for i := range transactions {
		transactionDetails = append(transactionDetails, TransactionDetails{
			MobileNumber: wallet.MobileNumber,
			Message:      transactions[i].Message,
			Date:         transactions[i].CreatedAt,
			Amount:       transactions[i].Amount,
		})
	}
	return transactionDetails, nil
}
