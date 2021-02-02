package routing

import (
	"context"
	"fmt"
	"github.com/bradenrayhorn/ledger-auth/database"
	"github.com/bradenrayhorn/ledger-auth/internal/db"
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	applyRoutes(router)
	return router
}

func applyRoutes(router *gin.Engine) {
	router.GET("/api/v1/health-check", func(context *gin.Context) {
		context.String(http.StatusOK, "ok")
	})

	router.POST("/api/v1/register", func(_ *gin.Context) {
		err := db.New(database.DB).CreateUser(context.Background(), db.CreateUserParams{
			ID:       uuid.New().String(),
			Username: "test user",
			Password: "my password",
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	})
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
