package tests

import (
	"github.com/bradenrayhorn/ledger-auth/config"
	"github.com/bradenrayhorn/ledger-auth/database"
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
	config.LoadConfig()
	database.Setup()

	r = routing.MakeRouter()

	database.DB.MustExec("TRUNCATE TABLE users")

	return m.Run()
}
