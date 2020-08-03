package handlers

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"testing"
	"time"
)

func TestGetUsers(t *testing.T) {
	testRouter.GET("/users", GetUsers)
	t.Run("get all users", func(t *testing.T) {
		//Arrange
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "location"}).
				AddRow(1, "Sherlock Holmes", "Baker street").
				AddRow(2, "John Watson", "baker street"))
		//Act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users", nil)
		testRouter.ServeHTTP(w, req)
		//Assert
		assert.Equal(t, 200, w.Code)
		response := w.Body.String()
		assert.NotEmpty(t, response)
		data := &Users{}
		_ = json.Unmarshal([]byte(response), data)
		assert.Equal(t, len(data.Users), 2)
		assert.Equal(t, data.Users[0].Name, "Sherlock Holmes")
		assert.Equal(t, data.Users[1].Name, "John Watson")
	})

	t.Run("get all users - no users available", func(t *testing.T) {
		//Arrange
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "location"}))
		//Act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users", nil)
		testRouter.ServeHTTP(w, req)
		//Assert
		assert.Equal(t, 200, w.Code)
		response := w.Body.String()
		assert.NotEmpty(t, response)
		data := &Users{}
		_ = json.Unmarshal([]byte(response), data)
		assert.Equal(t, len(data.Users), 0)
	})
}

func TestGetUserDetails(t *testing.T) {
	testRouter.GET("/api/v1/user/:user_id", GetUserDetails)
	t.Run("get user details with no transactions", func(t *testing.T) {
		//Arrange
		expectedUserDetails := UserDetails{
			Id:                   1,
			Name:                 "Sherlock Holmes",
			Balance:              100,
			SentTransactions:     make([]TransactionDetails,0),
			ReceivedTransactions: make([]TransactionDetails,0),
		}
		mock.ExpectQuery(`SELECT \* FROM "users"`).WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "location"}).
				AddRow(1, "Sherlock Holmes", "Baker street"))
		mock.ExpectQuery(`SELECT \* FROM "wallets"`).WithArgs(1).
			WillReturnRows(sqlmock.NewRows(
				[]string{"id", "balance", "user_id", "mobile_number"}).
				AddRow(1, 100, 1, 9999))
		mock.ExpectQuery(`SELECT \* FROM "transactions"`).WithArgs(1).
			WillReturnRows(sqlmock.NewRows(
				[]string{"id", "created_at", "from_wallet", "to_wallet", "amount", "message"}))
		mock.ExpectQuery(`SELECT \* FROM "transactions"`).WithArgs(1).
			WillReturnRows(sqlmock.NewRows(
				[]string{"id", "created_at", "from_wallet", "to_wallet", "amount", "message"}))
		//Act
		w:= httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/user/1", nil)
		testRouter.ServeHTTP(w, req)

		//Assert
		assert.Equal(t, 200, w.Code)
		response := w.Body.String()
		assert.NotEmpty(t, response)
		actualUserDetails := &UserInfo{}
		_ = json.Unmarshal([]byte(response), actualUserDetails)
		assert.True(t, reflect.DeepEqual(actualUserDetails.UserDetails, expectedUserDetails))
	})

	t.Run("get user details with credits and debits", func(t *testing.T) {
		//Arrange
		now := time.Now()
		credits := TransactionDetails{
			MobileNumber: 9999,
			Message:      "Happy new year",
			Amount:       20,
			Date:         now,
		}
		debits := TransactionDetails{
			MobileNumber: 9999,
			Message:      "Happy birthday",
			Amount:       50,
			Date:         now,
		}

		expectedUserDetails := UserDetails{
			Id:                   1,
			Name:                 "Sherlock Holmes",
			Balance:              100,
			SentTransactions: append(make([]TransactionDetails, 0), debits),
			ReceivedTransactions: append(make([]TransactionDetails, 0), credits),
		}
		mock.ExpectQuery(`SELECT \* FROM "users"`).WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "location"}).
				AddRow(1, "Sherlock Holmes", "Baker street"))
		mock.ExpectQuery(`SELECT \* FROM "wallets"`).WithArgs(1).
			WillReturnRows(sqlmock.NewRows(
				[]string{"id", "balance", "user_id", "mobile_number"}).
				AddRow(1, 100, 1, 9999))
		mock.ExpectQuery(`SELECT \* FROM "transactions"`).WithArgs(1).
			WillReturnRows(sqlmock.NewRows(
				[]string{"id", "created_at", "from_wallet", "to_wallet", "amount", "message"}).
				AddRow(2, now, 2, 1, 50, "Happy birthday"))
		mock.ExpectQuery(`SELECT \* FROM "transactions"`).WithArgs(1).
			WillReturnRows(sqlmock.NewRows(
				[]string{"id", "created_at", "from_wallet", "to_wallet", "amount", "message"}).
				AddRow(1, now, 1, 2, 20, "Happy new year"))
		//Act
		w:= httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/user/1", nil)
		testRouter.ServeHTTP(w, req)

		//Assert
		assert.Equal(t, 200, w.Code)
		response := w.Body.String()
		assert.NotEmpty(t, response)
		userInfo := &UserInfo{}
		_ = json.Unmarshal([]byte(response), userInfo)
		assert.Equal(t, len(userInfo.UserDetails.ReceivedTransactions), len(expectedUserDetails.ReceivedTransactions))
		assert.Equal(t, len(userInfo.UserDetails.SentTransactions), len(expectedUserDetails.SentTransactions))
		assert.True(t, reflect.DeepEqual(userInfo.UserDetails.ReceivedTransactions[0].Message, expectedUserDetails.ReceivedTransactions[0].Message))
	})
}