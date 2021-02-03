package routing

import (
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
