package service

import (
	"fmt"
	"fundstransfer/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type UserDetails struct {
	Id                   uint                 `json:"id"`
	Name                 string               `json:"name"`
	Balance              int                  `json:"balance"`
	SentTransactions     []TransactionDetails `json:"debits"`
	ReceivedTransactions []TransactionDetails `json:"credits"`
}

type TransactionDetails struct {
	ToAccount uint      `json:"to_account"`
	Message   string    `json:"message"`
	Amount    int       `json:"amount"`
	Date      time.Time `json:"transactionTime"`
}

func GetUsers(c *gin.Context) {
	users, err := models.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": &users})
}

func GetAccounts(c *gin.Context) {
	accounts, err := models.GetAllAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, gin.H{"accounts": &accounts})
}

func GetUserDetails(userId int) (*UserDetails, error) {
	var user models.User
	var account models.Account
	if err := user.GetById(userId); err != nil {
		return nil, fmt.Errorf("error retrieving user %v", err)
	}
	if err := account.GetAccountByUser(user.ID); err != nil {
		return nil, fmt.Errorf("error retrieving account %v", err)
	}
	sentTransactions, err := models.GetTransactionsByAccount(account.ID, true)
	if err != nil {
		return nil, fmt.Errorf("error retrieving sentTransactions %v", err)
	}
	debits := make([]TransactionDetails, 0)
	for i := range sentTransactions {
		debits = append(debits, TransactionDetails{
			ToAccount: sentTransactions[i].ToAccountId,
			Message:   sentTransactions[i].Message,
			Date:      sentTransactions[i].CreatedAt,
			Amount:    sentTransactions[i].Amount,
		})
	}

	receivedTransactions, err := models.GetTransactionsByAccount(account.ID, false)
	if err != nil {
		return nil, fmt.Errorf("error retrieving sentTransactions %v", err)
	}
	credits := make([]TransactionDetails, 0)
	for i := range receivedTransactions {
		credits = append(credits, TransactionDetails{
			ToAccount: receivedTransactions[i].ToAccountId,
			Message:   receivedTransactions[i].Message,
			Date:      receivedTransactions[i].CreatedAt,
			Amount:    receivedTransactions[i].Amount,
		})
	}

	userDetails := &UserDetails{
		Id:                   user.ID,
		Name:                 user.Name,
		Balance:              account.Balance,
		SentTransactions:     debits,
		ReceivedTransactions: credits,
	}
	return userDetails, nil
}
