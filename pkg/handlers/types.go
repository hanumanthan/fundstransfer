package handlers

import "time"

type Transaction struct {
	To      int32  `json:"mobile_number" binding:"required"`
	Amount  int    `json:"amount" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type UserDetails struct {
	Id                   uint                 `json:"id"`
	Name                 string               `json:"name"`
	Balance              int                  `json:"balance"`
	SentTransactions     []TransactionDetails `json:"debits"`
	ReceivedTransactions []TransactionDetails `json:"credits"`
}

type TransactionDetails struct {
	MobileNumber int32     `json:"mobile_number"`
	Message      string    `json:"message"`
	Amount       int       `json:"amount"`
	Date         time.Time `json:"transaction_time"`
}
