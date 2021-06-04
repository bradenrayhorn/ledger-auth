package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bradenrayhorn/ledger-auth/database"
	"github.com/stretchr/testify/suite"
)

type SessionHTTPSuite struct {
	suite.Suite
}

func (s *SessionHTTPSuite) TearDownTest() {
	database.DB.MustExec("truncate table users")
	database.DB.MustExec("truncate table active_sessions")
	database.RDB.FlushDB(context.Background())
}

type GetSessionsResponse struct {
	Sessions []GetSessionsResponseSession `json:"sessions"`
}

type GetSessionsResponseSession struct {
	CreatedAt    string `json:"created_at"`
	IP           string `json:"ip"`
	LastAccessed string `json:"last_accessed"`
	UserAgent    string `json:"user_agent"`
}

func (s *SessionHTTPSuite) TestGetSessions() {
	user := makeUser(s.T())
	sessionID1 := getSessionID(&s.Suite, user)

	database.DB.MustExec("INSERT INTO active_sessions (user_id, session_id) VALUES(?,?);", user.ID, "x")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/sessions", nil)
	req.Header.Add("Cookie", "session_id="+sessionID1)
	req.Header.Add("User-Agent", "TestUserAgent")
	req.RemoteAddr = "1.2.3.4:80"
	r.ServeHTTP(w, req)

	s.Require().Equal(http.StatusOK, w.Code)

	var body GetSessionsResponse
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	s.Require().Len(body.Sessions, 1)
	s.Require().Equal("TestUserAgent", body.Sessions[0].UserAgent)
	s.Require().Equal("1.2.3.4", body.Sessions[0].IP)
}

func TestSessionHTTPSuite(t *testing.T) {
	suite.Run(t, new(SessionHTTPSuite))
}
