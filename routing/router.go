package routing

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func MakeRouter() *gin.Engine {
	gin.DefaultWriter = makeZapWriter()
	router := gin.New()

	router.Use(ginzap.Ginzap(zap.L(), time.RFC3339, false))
	router.Use(ginzap.RecoveryWithZap(zap.L(), true))
	router.Use(ErrorReporter())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = viper.GetStringSlice("allowed_origins")
	corsConfig.AllowCredentials = viper.GetBool("allow_credentials")
	router.Use(cors.New(corsConfig))

	applyRoutes(router)
	return router
}

func applyRoutes(router *gin.Engine) {
	router.GET("/health-check", func(context *gin.Context) {
		context.String(http.StatusOK, "ok")
	})

	authApi := router.Group("/api/v1/auth")
	api := router.Group("/api/v1")
	api.Use(AuthMiddleware())

	authApi.POST("/register", Register)
	authApi.POST("/login", Login)

	api.GET("/me", Me)
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
