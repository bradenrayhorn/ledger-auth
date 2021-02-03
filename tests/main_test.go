package tests

import (
	"context"
	"github.com/bradenrayhorn/ledger-auth/config"
	"github.com/bradenrayhorn/ledger-auth/database"
	"github.com/bradenrayhorn/ledger-auth/internal/db"
	"github.com/bradenrayhorn/ledger-auth/routing"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
	"testing"
)

var r *gin.Engine

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	viper.AddConfigPath("../")
	viper.Set("rsa_path", "../jwt_rsa")
	config.LoadConfig()
	database.Setup()

	r = routing.MakeRouter()

	_ = db.New(database.DB).UsersTruncate(context.Background())

	return m.Run()
}
