package handlers

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestGetTransactions(t *testing.T) {
	testRouter.GET("/transactions", GetTransactions)
	t.Run("get all transactions", func(t *testing.T) {
		//Arrange
		now := time.Now()
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "transactions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "from_wallet", "to_wallet", "amount", "message"}).
				AddRow(1, now, 1, 2, 100, "Gung hay fat choy").
				AddRow(2, now, 2, 1, 50, "Happy new year"))
		//Act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/transactions", nil)
		testRouter.ServeHTTP(w, req)
		//Assert
		assert.Equal(t, 200, w.Code)
		response := w.Body.String()
		assert.NotEmpty(t, response)
		data := &Transactions{}
		_ = json.Unmarshal([]byte(response), data)
		assert.Equal(t, len(data.Transactions), 2)
		assert.Equal(t, data.Transactions[0].Message, "Gung hay fat choy")
		assert.Equal(t, data.Transactions[1].Message, "Happy new year")
	})

	t.Run("get all transactions - no transactions available", func(t *testing.T) {
		//Arrange
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "transactions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "from_wallet", "to_wallet", "amount", "message"}))
		//Act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/transactions", nil)
		testRouter.ServeHTTP(w, req)
		//Assert
		assert.Equal(t, 200, w.Code)
		response := w.Body.String()
		assert.NotEmpty(t, response)
		data := &Transactions{}
		_ = json.Unmarshal([]byte(response), data)
		assert.Equal(t, len(data.Transactions), 0)
	})
}


func TestTransact(t *testing.T) {
	testRouter.POST("/api/v1/user/1/transact", Transact)
	t.Run("transact - invalid input", func(t *testing.T) {
		//Arrange
		input := strings.NewReader(`{"mobile_number" : 8888,"message" : "enjoy","amount123": 20}`)
		//Act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/user/1/transact", input)
		testRouter.ServeHTTP(w, req)
		//Assert
		assert.Equal(t, 400, w.Code)
	})

	t.Run("invalid recipient mobile number", func(t *testing.T){
		//Arrange
		input := strings.NewReader(`{"mobile_number" : 8888,"message" : "enjoy","amount": 20}`)
		mock.ExpectQuery(`SELECT \* FROM "wallets"`).
			WillReturnRows(sqlmock.NewRows(
				[]string{"id", "balance", "user_id", "mobile_number"}).
			AddRow(1, 100, 1, 9999))
		mock.ExpectQuery(`SELECT \* FROM "wallets"`).WithArgs(8888).
			WillReturnError(gorm.ErrRecordNotFound)
		//Act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/user/1/transact", input)
		testRouter.ServeHTTP(w, req)
		//Assert
		assert.Equal(t, 400, w.Code)
		assert.Equal(t, w.Body.String(), `{"error":"invalid wallet mobile number"}`)

	})

	t.Run("insufficient balance to transfer", func(t *testing.T){
		//Arrange
		input := strings.NewReader(`{"mobile_number" : 8888,"message" : "enjoy","amount": 100}`)
		mock.ExpectQuery(`SELECT \* FROM "wallets"`).
			WillReturnRows(sqlmock.NewRows(
				[]string{"id", "balance", "user_id", "mobile_number"}).
				AddRow(1, 20, 1, 9999))
		mock.ExpectQuery(`SELECT \* FROM "wallets"`).WithArgs(8888).
			WillReturnRows(sqlmock.NewRows(
			[]string{"id", "balance", "user_id", "mobile_number"}).
			AddRow(2, 50, 2, 8888))
		//Act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/user/1/transact", input)
		testRouter.ServeHTTP(w, req)
		//Assert
		assert.Equal(t, 500, w.Code)
		assert.Equal(t, w.Body.String(), `{"error":"insufficient balance to transfer"}`)
	})

}