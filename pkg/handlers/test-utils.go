package handlers

import (
	"database/sql"
	"fundstransfer/pkg/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var mock sqlmock.Sqlmock
var testDb *sql.DB
var testRouter *gin.Engine

func init() {
	testRouter = setupRouter()
	testDb, mock, _ = sqlmock.New()
	models.DB, _ = gorm.Open("sqlite3", testDb)
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	return r
}
