package routing

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

func createDeviceCookie(w http.ResponseWriter) (string, error) {
	bytes := make([]byte, 64)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	deviceID := base64.RawURLEncoding.EncodeToString(bytes)

	cookie := http.Cookie{
		Name:     "device_id",
		Value:    deviceID,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Domain:   viper.GetString("cookie_domain"),
		Secure:   viper.GetBool("cookie_secure"),
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 438000),
	}

	http.SetCookie(w, &cookie)
	return deviceID, nil
}
