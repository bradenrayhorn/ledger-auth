package services

import (
	"context"
	"fmt"

	"github.com/bradenrayhorn/ledger-auth/internal/db"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

var deviceRedisPrefix = "devices:"

type DeviceService struct {
	client *redis.Client
}

func NewDeviceService(client *redis.Client) DeviceService {
	return DeviceService{client: client}
}

func (s DeviceService) RecognizeDevice(ctx context.Context, userID uuid.UUID, deviceID string) error {
	key := deviceRedisPrefix + userID.String()
	err := s.client.SAdd(ctx, key, deviceID).Err()
	if err != nil {
		return err
	}
	err = s.client.Do(ctx, "EXPIREMEMBER", key, deviceID, 7776000).Err()
	return err
}

func (s DeviceService) DoesRecognizeDevice(ctx context.Context, userID uuid.UUID, deviceID string) (bool, error) {
	key := deviceRedisPrefix + userID.String()
	return s.client.SIsMember(ctx, key, deviceID).Result()
}

func (s DeviceService) NotifyOfNewDevice(user db.User, ip string) error {
	if user.Email.Valid {
		return NewEmailService(ServiceMailClient).SendEmail("Ledger Security Notice", fmt.Sprintf("A login has occurred from a new device with IP %s.", ip), user.Email.String)
	}
	return nil
}
