package routing

import (
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

func MakeRouter() *gin.Engine {
	gin.DefaultWriter = makeZapWriter()
	router := gin.New()

	router.Use(ginzap.Ginzap(zap.L(), time.RFC3339, false))
	router.Use(ginzap.RecoveryWithZap(zap.L(), true))
	router.Use(ErrorReporter())

	applyRoutes(router)
	return router
}

func applyRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	v1.GET("/health-check", func(context *gin.Context) {
		context.String(http.StatusOK, "ok")
	})

	authApi := v1.Group("/auth")
	authApi.POST("/register", Register)
}

type ZapWriter func([]byte) (int, error)

func (fn ZapWriter) Write(data []byte) (int, error) {
	return fn(data)
}

func makeZapWriter() io.Writer {
	return ZapWriter(func(data []byte) (int, error) {
		zap.S().Debugf("%s", data)
		return 0, nil
	})
}
