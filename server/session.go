package server

import (
	"context"

	"github.com/bradenrayhorn/ledger-auth/routing"
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

	sessionID, userID, err := routing.GetSessionFromCookie(req.GetSessionID(), req.GetIP(), req.GetUserAgent())
	if err != nil {
		return response, err
	}

	response.Session = &session.Session{
		SessionID: sessionID,
		UserID:    userID,
	}

	return response, nil
}
