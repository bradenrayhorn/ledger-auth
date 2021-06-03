package routing

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/bradenrayhorn/ledger-auth/database"
	"github.com/bradenrayhorn/ledger-auth/services"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type CookieValue struct {
	SessionID string
	Signature []byte
}

func GetSessions(c *gin.Context) {
	sessionService := services.NewSessionService(database.RDB)
	sessions, err := sessionService.GetActiveSessions(context.Background(), c.GetString("user_id"))

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusOK, map[string]interface{}{
		"sessions": sessions,
	})
}

func createSession(w http.ResponseWriter, userID string, ip string, userAgent string) error {
	sessionService := services.NewSessionService(database.RDB)
	sessionID, err := sessionService.CreateSession(userID, services.SessionData{
		IP:        ip,
		UserAgent: userAgent,
	})
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
		SameSite: http.SameSiteStrictMode,
		Domain:   viper.GetString("cookie_domain"),
		Secure:   viper.GetBool("cookie_secure"),
		Path:     "/",
	}

	http.SetCookie(w, &cookie)
	return nil
}

func getSession(cookieValueString string, ip string, userAgent string) (string, string, error) {
	decodedCookie, err := base64.RawURLEncoding.DecodeString(cookieValueString)
	if err != nil {
		return "", "", err
	}

	var cookieValue CookieValue
	if err = json.Unmarshal(decodedCookie, &cookieValue); err != nil {
		return "", "", err
	}

	if getHMACService().ValidateSignature([]byte(cookieValue.SessionID), cookieValue.Signature) != nil {
		return "", "", err
	}

	sessionService := services.NewSessionService(database.RDB)
	userID, err := sessionService.GetSession(context.Background(), cookieValue.SessionID, services.SessionData{
		IP:        ip,
		UserAgent: userAgent,
	})
	return cookieValue.SessionID, userID, err
}

func deleteSession(w http.ResponseWriter, sessionID string) error {
	sessionService := services.NewSessionService(database.RDB)
	err := sessionService.DeleteSession(context.Background(), sessionID)

	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:     "session_id",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Domain:   viper.GetString("cookie_domain"),
		Secure:   viper.GetBool("cookie_secure"),
		Path:     "/",
	}

	http.SetCookie(w, &cookie)
	return nil
}

func getHMACService() services.HMACService {
	return services.NewHMACService([]byte(viper.GetString("session_hash_key")))
}
