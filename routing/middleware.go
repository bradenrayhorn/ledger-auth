package routing

import (
	"github.com/bradenrayhorn/ledger-auth/jwt"
	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"
)

type RequestError interface {
	Error() string
	Code() int
}

func ErrorReporter() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		errors := c.Errors.ByType(gin.ErrorTypeAny)

		if len(errors) > 0 {
			err := errors[0].Err
			requestErr, ok := err.(RequestError)
			statusCode := 500
			if ok {
				statusCode = requestErr.Code()
			}

			c.IndentedJSON(statusCode, map[string]interface{}{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

	}
}

func getToken(header string) string {
	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := getToken(c.GetHeader("Authorization"))
		if len(tokenString) == 0 {
			tokenString, _ = c.GetQuery("auth")
			tokenString, _ = url.QueryUnescape(tokenString)
		}

		token, err := jwt.ParseToken(tokenString)

		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "invalid api token",
			})
			c.Abort()
			return
		}

		claims := token.Claims.(jwtGo.MapClaims)

		c.Set("user_id", claims["user_id"])
		c.Set("user_username", claims["user_username"])
		c.Next()
	}
}
