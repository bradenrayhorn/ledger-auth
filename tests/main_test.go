package tests

import (
	"os"
	"testing"

	"github.com/bradenrayhorn/ledger-auth/config"
	"github.com/bradenrayhorn/ledger-auth/database"
	"github.com/bradenrayhorn/ledger-auth/routing"
	"github.com/bradenrayhorn/ledger-auth/services"
	"github.com/gin-gonic/gin"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
)

var r *gin.Engine

type mockMailClient struct {
	mock.Mock
}

func (c *mockMailClient) Send(message *mail.SGMailV3) (*rest.Response, error) {
	args := c.Called(message)
	return args.Get(0).(*rest.Response), args.Error(1)
}

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
	services.ServiceMailClient = new(mockMailClient)

	return m.Run()
}
