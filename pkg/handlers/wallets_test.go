package handlers

import (
	"database/sql"
	"fundstransfer/pkg/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

var mock sqlmock.Sqlmock
var db  *sql.DB

func init() {
	db, mock, _ = sqlmock.New()
	models.DB, _ = gorm.Open("sqlite3", db)
}

func setupRouter() *gin.Engine{
	r := gin.Default()
	r.GET("/wallets", GetWallets)
	return r
}

func TestGetWallets(t *testing.T) {
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "wallets"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "balance", "user_id", "mobile_number"}).
			AddRow(1, 1, 1, 1))

	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/wallets", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"wallets":[{"ID":1,"Balance":1,"UserId":1,"MobileNumber":1}]}`, w.Body.String())
}

