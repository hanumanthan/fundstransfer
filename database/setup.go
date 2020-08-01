package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB

func ConnectDatabase() {
	db, err := gorm.Open("sqlite3", "payments.db")
	if err != nil {
		fmt.Errorf("failed to connect to database")
	}
	DB = db
}
