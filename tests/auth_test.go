package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bradenrayhorn/ledger-auth/database"
	"github.com/bradenrayhorn/ledger-auth/internal/db"
	"github.com/bradenrayhorn/ledger-auth/jwt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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

func TestLogin(t *testing.T) {
	err := db.New(database.DB).CreateUser(context.Background(), db.CreateUserParams{
		ID:       uuid.NewString(),
		Username: "test",
		Password: "$2a$10$naqzJWUaOFm1/512Od.wPO4H8Vh8K38IGAb7rtgFizSflLVhpgMRG",
	})
	assert.Nil(t, err)

	testLogin(t, http.StatusOK, "test", "password")
}

func TestCannotLoginWithInvalidUsername(t *testing.T) {
	testLogin(t, http.StatusUnprocessableEntity, "test-bad", "password")
}

func TestCannotLoginWithInvalidPassword(t *testing.T) {
	err := db.New(database.DB).CreateUser(context.Background(), db.CreateUserParams{
		ID:       uuid.NewString(),
		Username: "test",
		Password: "$2a$10$naqzJWUaOFm1/512Od.wPO4H8Vh8K38IGAb7rtgFizSflLVhpgMRG",
	})
	assert.Nil(t, err)

	testLogin(t, http.StatusUnprocessableEntity, "test", "password-wrong")
}

type GetMeResponse struct {
	Id string `json:"id"`
}

func TestShowMe(t *testing.T) {
	user := db.User{
		ID:        uuid.NewString(),
		Username:  "test",
		Password:  "$2a$10$naqzJWUaOFm1/512Od.wPO4H8Vh8K38IGAb7rtgFizSflLVhpgMRG",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	err := db.New(database.DB).CreateUser(context.Background(), db.CreateUserParams{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
	})
	assert.Nil(t, err)

	token, _ := jwt.CreateToken(user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/me", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var body GetMeResponse
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	assert.Equal(t, user.ID, body.Id)
}

func TestCannotShowMeUnauthenticated(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/me", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCannotShowMeWithExpiredToken(t *testing.T) {
	user := db.User{
		ID:        uuid.NewString(),
		Username:  "test",
		Password:  "$2a$10$naqzJWUaOFm1/512Od.wPO4H8Vh8K38IGAb7rtgFizSflLVhpgMRG",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	err := db.New(database.DB).CreateUser(context.Background(), db.CreateUserParams{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
	})
	assert.Nil(t, err)

	viper.Set("token_expiration", -10*time.Second)
	token, _ := jwt.CreateToken(user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/me", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	viper.Set("token_expiration", 24*time.Hour)
}

func testLogin(t *testing.T, expectedStatus int, username string, password string) {
	w := httptest.NewRecorder()
	reader := strings.NewReader(fmt.Sprintf("username=%s&password=%s", username, password))
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	assert.Equal(t, expectedStatus, w.Code)
}
