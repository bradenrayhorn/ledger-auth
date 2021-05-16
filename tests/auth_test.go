package tests

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/bradenrayhorn/ledger-auth/database"
	"github.com/bradenrayhorn/ledger-auth/internal/db"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthSuite struct {
	suite.Suite
}

func (s *AuthSuite) TearDownTest() {
	database.DB.MustExec("truncate table users")
	database.RDB.FlushDB(context.Background())
}

func (s *AuthSuite) TestRegister() {
	w := httptest.NewRecorder()
	reader := strings.NewReader("username=test&password=password")
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	s.Assert().Equal(http.StatusOK, w.Code)
}

func (s *AuthSuite) TestCannotRegisterTwice() {
	_ = makeUser(s.T())

	w := httptest.NewRecorder()
	reader := strings.NewReader("username=test&password=password")
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	s.Assert().Equal(http.StatusUnprocessableEntity, w.Code)
}

func (s *AuthSuite) TestCannotRegisterWithNoData() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	s.Assert().Equal(http.StatusUnprocessableEntity, w.Code)
}

func (s *AuthSuite) TestLogin() {
	_ = makeUser(s.T())

	testLogin(s.T(), http.StatusOK, "test", "password")
}

func (s *AuthSuite) TestCannotLoginWithInvalidUsername() {
	testLogin(s.T(), http.StatusUnprocessableEntity, "test-bad", "password")
}

func (s *AuthSuite) TestCannotLoginWithInvalidPassword() {
	_ = makeUser(s.T())

	testLogin(s.T(), http.StatusUnprocessableEntity, "test", "password-wrong")
}

type StaticReader struct {
}

func (r StaticReader) Read(p []byte) (n int, err error) {
	p = append(p, 1)
	return 1, nil
}

func (s *AuthSuite) TestCannotLoginIfGeneratedSessionIDExists() {
	_ = makeUser(s.T())
	oldReader := rand.Reader
	rand.Reader = StaticReader{}
	defer func() {
		rand.Reader = oldReader
	}()

	database.RDB.Set(context.Background(), "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "id", time.Minute)

	testLogin(s.T(), http.StatusInternalServerError, "test", "password")
}

type GetMeResponse struct {
	Id string `json:"id"`
}

func (s *AuthSuite) TestShowMe() {
	user := makeUser(s.T())
	w := httptest.NewRecorder()
	reader := strings.NewReader(fmt.Sprintf("username=%s&password=%s", "test", "password"))
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	s.Require().Equal(http.StatusOK, w.Code)
	s.Require().Len(w.Result().Cookies(), 1)
	s.Require().Equal("session_id", w.Result().Cookies()[0].Name)

	sessionID := w.Result().Cookies()[0].Value
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/me", nil)
	req.Header.Add("Cookie", "session_id="+sessionID)
	r.ServeHTTP(w, req)

	s.Require().Equal(http.StatusOK, w.Code)

	var body GetMeResponse
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	s.Require().Equal(user.ID, body.Id)
}

func (s *AuthSuite) TestCannotShowMeWithExpiredSession() {
	makeUser(s.T())
	viper.Set("session_duration", "1s")
	w := httptest.NewRecorder()
	reader := strings.NewReader(fmt.Sprintf("username=%s&password=%s", "test", "password"))
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	s.Require().Equal(http.StatusOK, w.Code)
	s.Require().Len(w.Result().Cookies(), 1)
	s.Require().Equal("session_id", w.Result().Cookies()[0].Name)

	sessionID := w.Result().Cookies()[0].Value

	time.Sleep(time.Second * 2)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/me", nil)
	req.Header.Add("Cookie", "session_id="+sessionID)
	r.ServeHTTP(w, req)

	s.Require().Equal(http.StatusUnauthorized, w.Code)
}

func (s *AuthSuite) TestCannotShowMeUnauthenticated() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/me", nil)
	r.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}

func testLogin(t *testing.T, expectedStatus int, username string, password string) {
	w := httptest.NewRecorder()
	reader := strings.NewReader(fmt.Sprintf("username=%s&password=%s", username, password))
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	assert.Equal(t, expectedStatus, w.Code)
}

func makeUser(t *testing.T) db.User {
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
	return user
}
