package tests

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net"
	"testing"
	"time"

	"github.com/bradenrayhorn/ledger-auth/database"
	"github.com/bradenrayhorn/ledger-auth/routing"
	"github.com/bradenrayhorn/ledger-auth/server"
	"github.com/bradenrayhorn/ledger-auth/services"
	"github.com/bradenrayhorn/ledger-protos/session"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

type SessionSuite struct {
	lis *bufconn.Listener
	suite.Suite
}

func (s *SessionSuite) SetupTest() {
	s.lis = bufconn.Listen(bufSize)
	sv := grpc.NewServer()
	session.RegisterSessionAuthenticatorServer(sv, server.NewSessionAuthenticatorServer(database.RDB))
	go func() {
		if err := sv.Serve(s.lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func (s *SessionSuite) TearDownTest() {
	database.DB.MustExec("truncate table users")
	database.DB.MustExec("truncate table active_sessions")
	database.RDB.FlushDB(context.Background())
}

func (s *SessionSuite) bufDialer(context.Context, string) (net.Conn, error) {
	return s.lis.Dial()
}

func (s *SessionSuite) TestCanGetActiveSession() {
	ctx := context.Background()
	userID := uuid.Must(uuid.NewRandom()).String()
	database.RDB.HSet(ctx, "1234", map[string]interface{}{
		"user_id":       userID,
		"ip":            "18.8.9.1",
		"user_agent":    "TestAgent",
		"last_accessed": time.Now().Add(time.Minute * -10),
	})
	hmacService := services.NewHMACService([]byte(viper.GetString("session_hash_key")))
	sig, err := hmacService.SignData([]byte("1234"))
	s.Require().Nil(err)

	sessionValue, err := json.Marshal(routing.CookieValue{SessionID: "1234", Signature: sig})
	s.Require().Nil(err)

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(s.bufDialer), grpc.WithInsecure())
	s.Require().Nil(err)
	defer conn.Close()

	client := session.NewSessionAuthenticatorClient(conn)
	resp, err := client.Authenticate(ctx, &session.SessionAuthenticateRequest{SessionID: base64.RawURLEncoding.EncodeToString(sessionValue), UserAgent: "NewAgent", IP: "1.1.1.1"})

	s.Require().Nil(err)
	s.Require().NotNil(resp)
	s.Require().NotNil(resp.Session)
	s.Require().Equal(userID, resp.Session.UserID)
	s.Require().Equal("1234", resp.Session.SessionID)

	res, err := database.RDB.HGetAll(ctx, "1234").Result()
	s.Require().Nil(err)
	s.Require().Equal("1.1.1.1", res["ip"])
	s.Require().Equal("NewAgent", res["user_agent"])
	lastAccess, err := time.Parse(time.RFC3339, res["last_accessed"])
	s.Require().Nil(err)
	s.Require().True(lastAccess.After(time.Now().Add(time.Minute * -5)))
}

func (s *SessionSuite) TestCannotGetNonExistantSession() {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(s.bufDialer), grpc.WithInsecure())
	s.Require().Nil(err)
	defer conn.Close()

	client := session.NewSessionAuthenticatorClient(conn)
	resp, err := client.Authenticate(ctx, &session.SessionAuthenticateRequest{SessionID: "1234"})
	s.Assert().NotNil(err)
	s.Assert().Nil(resp)
}

func TestSessionSuite(t *testing.T) {
	suite.Run(t, new(SessionSuite))
}
