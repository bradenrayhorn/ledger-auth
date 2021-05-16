package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

type SessionService struct {
	rdb *redis.Client
}

func NewSessionService(client *redis.Client) SessionService {
	return SessionService{
		rdb: client,
	}
}

func (s SessionService) CreateSession(userID string) (string, error) {
	bytes := make([]byte, 64)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	sessionID := base64.RawURLEncoding.EncodeToString(bytes)

	exists, err := s.rdb.Exists(context.Background(), sessionID).Result()
	if err != nil {
		return "", err
	}
	if exists == 1 {
		return "", errors.New("failed to create session")
	}

	_, err = s.rdb.Set(context.Background(), sessionID, userID, viper.GetDuration("session_duration")).Result()
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (s SessionService) GetSession(ctx context.Context, sessionID string) (string, error) {
	userID, err := s.rdb.Get(ctx, sessionID).Result()
	if err != nil {
		return "", err
	}
	if len(userID) == 0 {
		return "", errors.New("invalid session")
	}

	return userID, nil
}
