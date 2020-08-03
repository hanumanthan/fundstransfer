package handlers

import (
	"encoding/json"
	"fundstransfer/pkg/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"
)

type Data struct {
	Transactions []models.Transaction `json:"transactions"`
}

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
		data := &Data{}
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
		data := &Data{}
		_ = json.Unmarshal([]byte(response), data)
		assert.Equal(t, len(data.Transactions), 0)
	})
}
