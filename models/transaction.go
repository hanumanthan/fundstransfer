package models

import (
	"fundstransfer/database"
	"github.com/jinzhu/gorm"
)

type Transaction struct {
	gorm.Model
	FromWallet uint
	ToWallet   uint
	Amount     int
	Message    string
}

func (t *Transaction) Save() error {
	if err := database.DB.Save(t).Error; err != nil {
		return err
	}
	return nil
}

func GetTransactionsByAccount(accId uint, isSender bool) ([]Transaction, error) {
	var query string
	if isSender {
		query = "from_wallet = ?"
	} else {
		query = "to_wallet = ?"
	}
	var transactions []Transaction
	if err := database.DB.Where(query, accId).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
