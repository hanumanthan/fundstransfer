package models

import (
	"fundstransfer/database"
	"github.com/jinzhu/gorm"
)

type Wallet struct {
	gorm.Model
	Balance      int
	UserId       uint
	MobileNumber int32
}

func (w *Wallet) GetWalletForUser(userId uint) error {
	if err := database.DB.Where("user_id = ?", userId).First(&w).Error; err != nil {
		return err
	}
	return nil
}

func (w *Wallet) Save() error {
	if err := database.DB.Save(w).Error; err != nil {
		return err
	}
	return nil
}

func GetAllWallets() ([]Wallet, error) {
	var wallets []Wallet
	if err := database.DB.Find(&wallets).Error; err != nil {
		return nil, err
	}
	return wallets, nil
}
