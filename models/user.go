package models

import (
	"fundstransfer/database"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name" binding:"required"`
	Location string `json:"location" binding:"required"`
}

func (u *User) GetById(id int) error {
	if err := database.DB.Where("id = ?", id).First(&u).Error; err != nil {
		return err
	}
	return nil
}

func GetAllUsers() ([]User, error) {
	var users []User
	if err := database.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
