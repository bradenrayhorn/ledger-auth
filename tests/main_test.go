package tests

import (
	"os"
	"testing"

	"github.com/bradenrayhorn/ledger-auth/config"
	"github.com/bradenrayhorn/ledger-auth/database"
	"github.com/bradenrayhorn/ledger-auth/routing"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var r *gin.Engine

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	viper.AddConfigPath("../")
	config.LoadConfig()
	viper.Set("trusted_proxies", []string{"0.0.0.0/0"})

	database.Setup()
	database.SetupRedis()

	r = routing.MakeRouter()

	return m.Run()
}
