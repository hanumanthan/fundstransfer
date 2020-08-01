package models

import (
	"fundstransfer/database"
	"github.com/jinzhu/gorm"
)

type Account struct {
	gorm.Model
	Balance int
	UserId uint
}

func (a *Account) GetAccountByUser(userId uint) error {
	if err := database.DB.Where("user_id = ?", userId).First(&a).Error; err != nil {
		return err
	}
	return nil
}

func (a *Account) Save() error {
	if err := database.DB.Save(a).Error; err != nil {
		return err
	}
	return nil
}

func GetAllAccounts() ([]Account, error) {
	var accounts []Account
	if err := database.DB.Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}
