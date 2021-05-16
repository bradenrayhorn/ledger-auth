package server

import (
	"context"
	"errors"

	"github.com/bradenrayhorn/ledger-auth/services"
	"github.com/bradenrayhorn/ledger-protos/session"
	"github.com/go-redis/redis/v8"
)

type SessionAuthenticatorServer struct {
	sessionService services.SessionService
	session.UnimplementedSessionAuthenticatorServer
}

func NewSessionAuthenticatorServer(client *redis.Client) SessionAuthenticatorServer {
	return SessionAuthenticatorServer{
		sessionService: services.NewSessionService(client),
	}
}

func (s SessionAuthenticatorServer) Authenticate(ctx context.Context, req *session.SessionAuthenticateRequest) (*session.SessionAuthenticateResponse, error) {
	response := &session.SessionAuthenticateResponse{}

	userID, err := s.sessionService.GetSession(ctx, req.GetSessionID())
	if err != nil {
		return response, err
	}
	if len(userID) == 0 {
		return response, errors.New("invalid session")
	}

	response.Session = &session.Session{
		SessionID: req.GetSessionID(),
		UserID:    userID,
	}

	return response, nil
}
