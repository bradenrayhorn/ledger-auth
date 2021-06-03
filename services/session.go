package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/bradenrayhorn/ledger-auth/repositories"
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

	err = repositories.CreateActiveSession(context.Background(), userID, sessionID)
	if err != nil {
		return "", err
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

func (s SessionService) GetActiveSessions(ctx context.Context, userID string) ([]map[string]interface{}, error) {
	formattedSessions := []map[string]interface{}{}
	sessions, err := repositories.GetActiveSessions(ctx, userID)
	if err != nil {
		return formattedSessions, nil
	}

	var results []*redis.StringStringMapCmd

	_, err = s.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, session := range sessions {
			results = append(results, pipe.HGetAll(ctx, session.SessionID))
		}
		return nil
	})

	if err != nil {
		return formattedSessions, nil
	}

	expiredSessions := []string{}

	for i, v := range results {
		if len(v.Val()) > 0 {
			formattedSessions = append(formattedSessions, map[string]interface{}{
				"ip":            v.Val()["ip"],
				"user_agent":    v.Val()["user_agent"],
				"last_accessed": v.Val()["last_accessed"],
				"created_at":    sessions[i].CreatedAt.Format(time.RFC3339),
			})
		} else {
			expiredSessions = append(expiredSessions, sessions[i].SessionID)
		}
	}

	if len(expiredSessions) > 0 {
		err = repositories.DeleteActiveSessions(ctx, expiredSessions)
		if err != nil {
			return formattedSessions, err
		}
	}

	return formattedSessions, nil
}

func makeSessionHash(userID string, data SessionData) map[string]interface{} {
	return map[string]interface{}{
		"user_id":       userID,
		"ip":            data.IP,
		"user_agent":    data.UserAgent,
		"last_accessed": time.Now().Format(time.RFC3339),
	}
}
