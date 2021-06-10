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
	"github.com/google/uuid"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockMailClient struct {
	mock.Mock
}

func (c *mockMailClient) Send(message *mail.SGMailV3) (*rest.Response, error) {
	args := c.Called(message)
	return args.Get(0).(*rest.Response), args.Error(1)
}

type UserTestSuite struct {
	suite.Suite
	mockMail *mockMailClient
}

func (s *UserTestSuite) TearDownTest() {
	database.DB.MustExec("truncate table users")
	database.DB.MustExec("truncate table active_sessions")
	database.RDB.FlushAll(context.Background())
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
	})).Return(&rest.Response{StatusCode: 200, Body: ""}, nil).Once()

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
}

func (s *UserTestSuite) TestCannotUpdateUserWithSameEmail() {
	user := makeUser(s.T())
	sessionID := getSessionID(&s.Suite, user)

	database.DB.MustExec("INSERT INTO users (id, username, password, email) VALUES (?, ?, ?, ?)", uuid.NewString(), "user2", "x", "test@test.com")

	w := httptest.NewRecorder()
	reader := strings.NewReader("email=test@test.com")
	req, _ := http.NewRequest("POST", "/api/v1/me/email", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", "session_id="+sessionID)
	r.ServeHTTP(w, req)

	s.Assert().Equal(http.StatusInternalServerError, w.Code)
	dbUser, err := repositories.GetUserByID(context.Background(), user.ID)
	s.Require().Nil(err)
	s.Assert().False(dbUser.Email.Valid)
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
