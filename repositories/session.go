package repositories

import (
	"context"

	"github.com/bradenrayhorn/ledger-auth/database"
	"github.com/bradenrayhorn/ledger-auth/internal/db"
	"github.com/google/uuid"
)

func CreateActiveSession(ctx context.Context, userID uuid.UUID, sessionID string) error {
	return db.New(database.DB).CreateActiveSession(ctx, db.CreateActiveSessionParams{
		UserID:    userID,
		SessionID: sessionID,
	})
}

func GetActiveSessions(ctx context.Context, userID uuid.UUID) ([]db.ActiveSession, error) {
	return db.New(database.DB).GetActiveSessions(ctx, userID)
}

func DeleteActiveSessions(ctx context.Context, sessionIDs []string) error {
	return db.New(database.DB).DeleteActiveSessionsByID(ctx, sessionIDs)
}

func DeleteActiveSessionsForUser(ctx context.Context, userID uuid.UUID) error {
	return db.New(database.DB).DeleteActiveSessionsForUser(ctx, userID)
}
