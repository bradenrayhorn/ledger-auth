package routing

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/bradenrayhorn/ledger-auth/database"
	"github.com/bradenrayhorn/ledger-auth/services"
	"github.com/spf13/viper"
)

type CookieValue struct {
	SessionID string
	Signature []byte
}

func createSession(w http.ResponseWriter, userID string) error {
	sessionService := services.NewSessionService(database.RDB)
	sessionID, err := sessionService.CreateSession(userID)
	if err != nil {
		return err
	}

	hmacService := getHMACService()
	sig, err := hmacService.SignData([]byte(sessionID))
	if err != nil {
		return err
	}

	value, err := json.Marshal(CookieValue{SessionID: sessionID, Signature: sig})
	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:     "session_id",
		Value:    base64.RawURLEncoding.EncodeToString(value),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Domain:   viper.GetString("cookie_domain"),
		Secure:   viper.GetBool("cookie_secure"),
		Path:     "/",
	}

	http.SetCookie(w, &cookie)
	return nil
}

func getSession(cookieValueString string) (string, error) {
	decodedCookie, err := base64.RawURLEncoding.DecodeString(cookieValueString)
	if err != nil {
		return "", err
	}

	var cookieValue CookieValue
	if err = json.Unmarshal(decodedCookie, &cookieValue); err != nil {
		return "", err
	}

	if getHMACService().ValidateSignature([]byte(cookieValue.SessionID), cookieValue.Signature) != nil {
		return "", err
	}

	sessionService := services.NewSessionService(database.RDB)
	return sessionService.GetSession(context.Background(), cookieValue.SessionID)
}

func getHMACService() services.HMACService {
	return services.NewHMACService([]byte(viper.GetString("session_hash_key")))
}
