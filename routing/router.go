package routing

import (
	"io"
	"net/http"
	"time"

	"github.com/bradenrayhorn/ledger-auth/database"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/ulule/limiter/v3"
	lgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/redis"
	"go.uber.org/zap"
)

func MakeRouter() *gin.Engine {
	gin.DefaultWriter = makeZapWriter()
	router := gin.New()
	router.TrustedProxies = viper.GetStringSlice("trusted_proxies")
	router.RemoteIPHeaders = []string{"X-Forwarded-For"}

	// middleware
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

	// rate limiting
	store, err := redis.NewStore(database.RDB)
	if err != nil {
		zap.S().Error(err)
	} else {
		authRateLimit := limiter.Rate{
			Period: time.Minute,
			Limit:  viper.GetInt64("rate_limit_auth"),
		}
		standardRateLimit := limiter.Rate{
			Period: time.Minute,
			Limit:  viper.GetInt64("rate_limit_standard"),
		}
		authMiddleware := lgin.NewMiddleware(limiter.New(store, authRateLimit, limiter.WithTrustForwardHeader(true)))
		standardMiddleware := lgin.NewMiddleware(limiter.New(store, standardRateLimit, limiter.WithTrustForwardHeader(true)))
		authApi.Use(authMiddleware)
		api.Use(standardMiddleware)
	}

	authApi.POST("/register", Register)
	authApi.POST("/login", Login)

	api.GET("/me", Me)
	api.GET("/sessions", GetSessions)
	api.POST("/auth/logout", Logout)
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
