package database

import (
	"fmt"
	"fundstransfer/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB

func ConnectDatabase() {
	db, err := gorm.Open("sqlite3", "payments.db")
	if err != nil {
		fmt.Errorf("failed to connect to database")
	}

	db.AutoMigrate(&models.User{})
	DB = db
}
