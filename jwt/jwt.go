package jwt

import (
	"fmt"
	"github.com/bradenrayhorn/ledger-auth/config"
	"github.com/bradenrayhorn/ledger-auth/internal/db"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

func CreateToken(user db.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"user_username": user.Username,
		"user_id":       user.ID,
		"exp":           time.Now().Add(viper.GetDuration("token_expiration")).Unix(),
	})
	return token.SignedString(config.RsaPrivate)
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return config.RsaPublic, nil
	})
}
