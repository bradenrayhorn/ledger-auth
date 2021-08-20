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
	"github.com/bradenrayhorn/ledger-auth/repositories"
	"github.com/bradenrayhorn/ledger-auth/routing"
	"github.com/bradenrayhorn/ledger-auth/services"
	"github.com/google/uuid"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AuthSuite struct {
	suite.Suite
}

func (s *AuthSuite) TearDownTest() {
	database.DB.MustExec("truncate table users")
	database.DB.MustExec("truncate table active_sessions")
	database.RDB.FlushAll(context.Background())
	services.ServiceMailClient = new(mockMailClient)
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
	testLogin(s.T(), http.StatusUnauthorized, "test-bad", "password")
}

func (s *AuthSuite) TestCannotLoginWithInvalidPassword() {
	_ = makeUser(s.T())

	testLogin(s.T(), http.StatusUnauthorized, "test", "password-wrong")
}

func (s *AuthSuite) TestLoginSendsEmailOnNewDeviceButNotKnownDevice() {
	user := makeUser(s.T())
	database.DB.MustExec("UPDATE users SET email = $1 WHERE id = $2", "test@test.com", user.ID.String())

	mockClient := new(mockMailClient)
	mockClient.On("Send", mock.MatchedBy(func(message *mail.SGMailV3) bool {
		return message.Personalizations[0].To[0].Address == "test@test.com"
	})).Return(&rest.Response{StatusCode: 200, Body: ""}, nil)
	services.ServiceMailClient = mockClient

	w := httptest.NewRecorder()
	reader := strings.NewReader(fmt.Sprintf("username=%s&password=%s", "test", "password"))
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	s.Require().Equal(http.StatusOK, w.Code)
	var deviceCookie *http.Cookie
	for _, c := range w.Result().Cookies() {
		if c.Name == "device_id" {
			deviceCookie = c
		}
	}
	s.Require().NotNil(deviceCookie)
	mockClient.AssertNumberOfCalls(s.T(), "Send", 1)

	w = httptest.NewRecorder()
	reader = strings.NewReader(fmt.Sprintf("username=%s&password=%s", "test", "password"))
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", "device_id="+deviceCookie.Value)
	r.ServeHTTP(w, req)

	s.Require().Equal(http.StatusOK, w.Code)

	mockClient.AssertNumberOfCalls(s.T(), "Send", 1)
	mockClient.AssertExpectations(s.T())
}

func (s *AuthSuite) TestAuthRouteIsRateLimited() {
	_ = makeUser(s.T())
	viper.Set("rate_limit_auth", "1")

	oldR := r
	r = routing.MakeRouter()
	defer func() {
		r = oldR
	}()

	testLogin(s.T(), http.StatusOK, "test", "password")
	testLogin(s.T(), http.StatusTooManyRequests, "test", "password")
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
	sessionID := getSessionID(&s.Suite, user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/me", nil)
	req.Header.Add("Cookie", "session_id="+sessionID)
	r.ServeHTTP(w, req)

	s.Require().Equal(http.StatusOK, w.Code)

	var body GetMeResponse
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	s.Require().Equal(user.ID.String(), body.Id)
}

func (s *AuthSuite) TestCannotShowMeWithExpiredSession() {
	user := makeUser(s.T())
	viper.Set("session_duration", "1s")

	sessionID := getSessionID(&s.Suite, user)

	time.Sleep(time.Second * 2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/me", nil)
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

func (s *AuthSuite) TestLogout() {
	user := makeUser(s.T())
	sessionID := getSessionID(&s.Suite, user)
	activeSessions, err := repositories.GetActiveSessions(context.Background(), user.ID)
	s.Require().Nil(err)
	s.Require().Len(activeSessions, 1)
	activeSessionIDs := make([]string, 0)
	for _, s := range activeSessions {
		activeSessionIDs = append(activeSessionIDs, s.SessionID)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/logout", nil)
	req.Header.Add("Cookie", "session_id="+sessionID)
	r.ServeHTTP(w, req)

	s.Require().Equal(http.StatusOK, w.Code)
	s.Require().Len(w.Result().Cookies(), 1)
	s.Require().Equal("session_id", w.Result().Cookies()[0].Name)
	s.Require().True(w.Result().Cookies()[0].Expires.Before(time.Now()))

	exists, err := database.RDB.Exists(context.Background(), activeSessionIDs...).Result()
	s.Require().Nil(err)
	s.Require().Equal(int64(0), exists)
}

func (s *AuthSuite) TestRevokeSessions() {
	user := makeUser(s.T())
	sessionID := getSessionID(&s.Suite, user)
	_ = getSessionID(&s.Suite, user)
	activeSessions, err := repositories.GetActiveSessions(context.Background(), user.ID)
	s.Require().Nil(err)
	s.Require().Len(activeSessions, 2)
	activeSessionIDs := make([]string, 0)
	for _, s := range activeSessions {
		activeSessionIDs = append(activeSessionIDs, s.SessionID)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/revoke", nil)
	req.Header.Add("Cookie", "session_id="+sessionID)
	r.ServeHTTP(w, req)

	s.Require().Equal(http.StatusOK, w.Code)

	exists, err := database.RDB.Exists(context.Background(), activeSessionIDs...).Result()
	s.Require().Nil(err)
	s.Require().Equal(int64(0), exists)

	activeSessions, err = repositories.GetActiveSessions(context.Background(), user.ID)
	s.Require().Nil(err)
	s.Require().Len(activeSessions, 0)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/me", nil)
	req.Header.Add("Cookie", "session_id="+sessionID)
	r.ServeHTTP(w, req)

	s.Require().Equal(http.StatusUnauthorized, w.Code)
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
		ID:        uuid.Must(uuid.NewRandom()),
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

func getSessionID(s *suite.Suite, user db.User) string {
	w := httptest.NewRecorder()
	reader := strings.NewReader(fmt.Sprintf("username=%s&password=%s", user.Username, "password"))
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	s.Require().Equal(http.StatusOK, w.Code)
	var cookie *http.Cookie
	for _, c := range w.Result().Cookies() {
		if c.Name == "session_id" {
			cookie = c
		}
	}
	s.Require().NotNil(cookie)

	return cookie.Value
}
