package routing

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("session_id")

		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "invalid session",
			})
			c.Abort()
			return
		}

		sessionID, userID, err := getSession(cookie.Value)
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "invalid session",
			})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("session_id", sessionID)
		c.Next()
	}
}
