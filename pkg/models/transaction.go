package models

import (
	"time"
)

type Transaction struct {
	ID         uint      `gorm:"primary_key" json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	FromWallet uint      `json:"from_wallet"`
	ToWallet   uint      `json:"to_wallet"`
	Amount     int       `json:"amount"`
	Message    string    `json:"message"`
}

func (t *Transaction) Save() error {
	if err := DB.Save(t).Error; err != nil {
		return err
	}
	return nil
}

func GetAllTransactions() ([]Transaction, error) {
	var transactions []Transaction
	if err := DB.Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func GetTransactionsByAccount(accId uint, isSender bool) ([]Transaction, error) {
	var query string
	if isSender {
		query = "from_wallet = ?"
	} else {
		query = "to_wallet = ?"
	}
	var transactions []Transaction
	if err := DB.Where(query, accId).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
