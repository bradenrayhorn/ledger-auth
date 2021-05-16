package routing

import (
	"context"
	"net/http"

	"github.com/bradenrayhorn/ledger-auth/database"
	"github.com/bradenrayhorn/ledger-auth/services"
	"github.com/spf13/viper"
)

func createSession(w http.ResponseWriter, userID string) error {
	sessionService := services.NewSessionService(database.RDB)
	sessionID, err := sessionService.CreateSession(userID)
	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Domain:   viper.GetString("cookie_domain"),
		Secure:   viper.GetBool("cookie_secure"),
		Path:     "/",
	}

	http.SetCookie(w, &cookie)
	return nil
}

func getSession(sessionID string) (string, error) {
	sessionService := services.NewSessionService(database.RDB)
	return sessionService.GetSession(context.Background(), sessionID)
}
