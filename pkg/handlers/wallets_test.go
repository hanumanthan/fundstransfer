package handlers

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestGetWallets(t *testing.T) {
	testRouter.GET("/wallets", GetWallets)
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "wallets"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "balance", "user_id", "mobile_number"}).
			AddRow(1, 100, 1, 1234))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/wallets", nil)
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"wallets":[{"id":1,"balance":100,"user_id":1,"mobile_number":1234}]}`, w.Body.String())
}
