package tests

import (
	"context"
	"github.com/bradenrayhorn/ledger-auth/database"
	"github.com/bradenrayhorn/ledger-auth/internal/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegister(t *testing.T) {
	w := httptest.NewRecorder()
	reader := strings.NewReader("username=test&password=password")
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCannotRegisterTwice(t *testing.T) {
	err := db.New(database.DB).CreateUser(context.Background(), db.CreateUserParams{
		ID:       uuid.NewString(),
		Username: "test",
		Password: "$2a$10$naqzJWUaOFm1/512Od.wPO4H8Vh8K38IGAb7rtgFizSflLVhpgMRG",
	})
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	reader := strings.NewReader("username=test&password=password")
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestCannotRegisterWithNoData(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}
