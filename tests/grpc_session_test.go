package tests

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/bradenrayhorn/ledger-auth/database"
	"github.com/bradenrayhorn/ledger-auth/server"
	"github.com/bradenrayhorn/ledger-protos/session"
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
	database.RDB.FlushDB(context.Background())
}

func (s *SessionSuite) bufDialer(context.Context, string) (net.Conn, error) {
	return s.lis.Dial()
}

func (s *SessionSuite) TestCanGetActiveSession() {
	ctx := context.Background()
	database.RDB.Set(ctx, "1234", "my user id", time.Hour)

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(s.bufDialer), grpc.WithInsecure())
	s.Require().Nil(err)
	defer conn.Close()

	client := session.NewSessionAuthenticatorClient(conn)
	resp, err := client.Authenticate(ctx, &session.SessionAuthenticateRequest{SessionID: "1234"})

	s.Require().Nil(err)
	s.Require().NotNil(resp)
	s.Require().NotNil(resp.Session)
	s.Require().Equal("my user id", resp.Session.UserID)
	s.Require().Equal("1234", resp.Session.SessionID)
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
