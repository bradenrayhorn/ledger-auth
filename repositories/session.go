package repositories

import (
	"context"

	"github.com/bradenrayhorn/ledger-auth/database"
	"github.com/bradenrayhorn/ledger-auth/internal/db"
	"github.com/jmoiron/sqlx"
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

func DeleteActiveSessions(ctx context.Context, sessionIDs []string) error {
	query, args, err := sqlx.In("DELETE FROM active_sessions WHERE session_id IN (?);", sessionIDs)
	if err != nil {
		return err
	}

	query = database.DB.Rebind(query)
	_, err = database.DB.ExecContext(ctx, query, args...)
	return err
}
