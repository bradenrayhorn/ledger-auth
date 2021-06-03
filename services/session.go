package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

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

type SessionData struct {
	IP        string
	UserAgent string
}

func (s SessionService) CreateSession(userID string, data SessionData) (string, error) {
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

	_, err = s.rdb.HSet(context.Background(), sessionID, makeSessionHash(userID, data)).Result()
	if err != nil {
		return "", err
	}

	err = s.rdb.Expire(context.Background(), sessionID, viper.GetDuration("session_duration")).Err()
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (s SessionService) GetSession(ctx context.Context, sessionID string, data SessionData) (string, error) {
	sessionData, err := s.rdb.HGetAll(ctx, sessionID).Result()
	if err != nil {
		return "", err
	}
	if len(sessionData) == 0 || len(sessionData["user_id"]) == 0 {
		return "", errors.New("invalid session")
	}

	err = s.rdb.HSet(context.Background(), sessionID, makeSessionHash(sessionData["user_id"], data)).Err()
	if err != nil {
		return "", err
	}

	return sessionData["user_id"], nil
}

func (s SessionService) DeleteSession(ctx context.Context, sessionID string) error {
	return s.rdb.Del(ctx, sessionID).Err()
}

func makeSessionHash(userID string, data SessionData) map[string]interface{} {
	return map[string]interface{}{
		"user_id":       userID,
		"ip":            data.IP,
		"user_agent":    data.UserAgent,
		"last_accessed": time.Now().Format(time.RFC3339),
	}
}
