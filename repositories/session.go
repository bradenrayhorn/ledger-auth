package repositories

import (
	"context"

	"github.com/bradenrayhorn/ledger-auth/database"
	"github.com/bradenrayhorn/ledger-auth/internal/db"
)

func CreateActiveSession(ctx context.Context, userID string, sessionID string) error {
	return db.New(database.DB).CreateActiveSession(ctx, db.CreateActiveSessionParams{
		UserID:    userID,
		SessionID: sessionID,
	})
}

func GetActiveSessions(ctx context.Context, userID string) ([]db.ActiveSession, error) {
	return db.New(database.DB).GetActiveSessions(ctx, userID)
}
