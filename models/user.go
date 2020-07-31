package models


type User struct {
	ID int `json:"id" gorm:"primary_key"`
	Name string `json:"name" binding:"required"`
	Location string `json:"location" binding:"required"`
}
