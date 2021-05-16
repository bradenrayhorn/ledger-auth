package routing

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/bradenrayhorn/ledger-auth/database"
	"github.com/spf13/viper"
)

func createSession(w http.ResponseWriter, userID string) error {
	bytes := make([]byte, 64)
	_, err := rand.Read(bytes)
	if err != nil {
		return err
	}
	sessionID := base64.RawURLEncoding.EncodeToString(bytes)

	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Domain:   viper.GetString("cookie_domain"),
		Secure:   viper.GetBool("cookie_secure"),
		Path:     "/",
	}

	_, err = database.RDB.Set(context.Background(), sessionID, userID, viper.GetDuration("session_duration")).Result()
	if err != nil {
		return err
	}

	http.SetCookie(w, &cookie)
	return nil
}

func getSession(sessionID string) (string, error) {
	userID, err := database.RDB.Get(context.Background(), sessionID).Result()
	if err != nil {
		return "", err
	}

	return userID, nil
}
