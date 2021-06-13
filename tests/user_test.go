package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bradenrayhorn/ledger-auth/database"
	"github.com/bradenrayhorn/ledger-auth/repositories"
	"github.com/bradenrayhorn/ledger-auth/services"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
	mockMail *mockMailClient
}

func (s *UserTestSuite) TearDownTest() {
	database.DB.MustExec("truncate table users")
	database.DB.MustExec("truncate table active_sessions")
	database.RDB.FlushAll(context.Background())
	services.ServiceMailClient = new(mockMailClient)
}

func (s *UserTestSuite) TestCannotUpdateWithInvalidEmail() {
	user := makeUser(s.T())
	sessionID := getSessionID(&s.Suite, user)

	w := httptest.NewRecorder()
	reader := strings.NewReader("email=test")
	req, _ := http.NewRequest("POST", "/api/v1/me/email", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", "session_id="+sessionID)
	r.ServeHTTP(w, req)

	s.Assert().Equal(http.StatusUnprocessableEntity, w.Code)
}

func (s *UserTestSuite) TestCanUpdateUserWithEmail() {
	user := makeUser(s.T())
	sessionID := getSessionID(&s.Suite, user)

	mockClient := new(mockMailClient)
	mockClient.On("Send", mock.MatchedBy(func(message *mail.SGMailV3) bool {
		return message.Personalizations[0].To[0].Address == "test@test.com"
	})).Return(&rest.Response{StatusCode: 200, Body: ""}, nil)

	services.ServiceMailClient = mockClient

	w := httptest.NewRecorder()
	reader := strings.NewReader("email=test@test.com")
	req, _ := http.NewRequest("POST", "/api/v1/me/email", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", "session_id="+sessionID)
	r.ServeHTTP(w, req)

	s.Assert().Equal(http.StatusOK, w.Code)
	user, err := repositories.GetUserByID(context.Background(), user.ID)
	s.Require().Nil(err)
	s.Assert().Equal("test@test.com", user.Email.String)

	mockClient.AssertNumberOfCalls(s.T(), "Send", 1)
	mockClient.AssertExpectations(s.T())
}

func (s *UserTestSuite) TestCanRemoveUserEmail() {
	user := makeUser(s.T())
	sessionID := getSessionID(&s.Suite, user)
	database.DB.MustExec("UPDATE users SET email = ? WHERE id = ?", "test@test.com", user.ID)

	mockClient := new(mockMailClient)
	services.ServiceMailClient = mockClient

	w := httptest.NewRecorder()
	reader := strings.NewReader("email=")
	req, _ := http.NewRequest("POST", "/api/v1/me/email", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", "session_id="+sessionID)
	r.ServeHTTP(w, req)

	s.Assert().Equal(http.StatusOK, w.Code)
	user, err := repositories.GetUserByID(context.Background(), user.ID)
	s.Require().Nil(err)
	s.Assert().Equal(false, user.Email.Valid)

	mockClient.AssertNumberOfCalls(s.T(), "Send", 0)
	mockClient.AssertExpectations(s.T())
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
