package config

import (
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/johanbrandhorst/certify"
	"github.com/johanbrandhorst/certify/issuers/vault"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	zapadapter "logur.dev/adapter/zap"
)

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/vault/secrets/")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Failed to read config: " + err.Error())
	}

	viper.SetDefault("cookie_secure", true)
	viper.SetDefault("allow_credentials", false)
	viper.SetDefault("session_duration", "24h")
	viper.SetDefault("grpc_port", "9000")
	viper.SetDefault("allowed_origins", "http://localhost")
	viper.SetDefault("trusted_proxies", "10.0.0.0/8")
	viper.SetDefault("rate_limit_auth", "6")
	viper.SetDefault("rate_limit_standard", "100")

	// mysql
	_ = viper.BindEnv("mysql_host", "MYSQL_HOST")
	_ = viper.BindEnv("mysql_port", "MYSQL_PORT")
	_ = viper.BindEnv("mysql_username", "MYSQL_USERNAME")
	_ = viper.BindEnv("mysql_password", "MYSQL_PASSWORD")
	_ = viper.BindEnv("mysql_database", "MYSQL_DATABASE")
	// redis
	_ = viper.BindEnv("redis_addr", "REDIS_ADDRESS")
	_ = viper.BindEnv("redis_db", "REDIS_DB")
	_ = viper.BindEnv("redis_username", "REDIS_USERNAME")
	_ = viper.BindEnv("redis_password", "REDIS_PASSWORD")
	// other
	_ = viper.BindEnv("cookie_domain", "COOKIE_DOMAIN")
	_ = viper.BindEnv("allowed_origins", "ALLOWED_ORIGINS")
	_ = viper.BindEnv("allow_credentials", "ALLOW_CREDENTIALS")
	_ = viper.BindEnv("session_duration", "SESSION_DURATION")
	_ = viper.BindEnv("grpc_port", "GRPC_PORT")
	_ = viper.BindEnv("session_hash_key", "SESION_HASH_KEY")
	_ = viper.BindEnv("ca_cert_path", "CA_CERT_PATH")
}

func LoadVaultToken() string {
	tokenBytes, err := ioutil.ReadFile(viper.GetString("vault_token_path"))
	if err != nil {
		log.Println(err)
		return ""
	}

	return strings.TrimSpace(string(tokenBytes))
}

func CreateCertify() (*certify.Certify, error) {
	url, err := url.Parse(viper.GetString("vault_url"))
	if err != nil {
		return nil, err
	}
	issuer := &vault.Issuer{
		URL:        url,
		Mount:      viper.GetString("vault_pki_mount"),
		AuthMethod: &vault.RenewingToken{Initial: LoadVaultToken()},
		Role:       viper.GetString("vault_pki_role"),
		TimeToLive: time.Hour * 24,
	}
	certify := &certify.Certify{
		Issuer:      issuer,
		CommonName:  viper.GetString("vault_pki_cn"),
		Cache:       certify.NewMemCache(),
		RenewBefore: time.Minute * 10,
		Logger:      zapadapter.New(zap.L()),
	}
	return certify, nil
}

var certPool *x509.CertPool

func GetCACertPool() *x509.CertPool {
	if certPool != nil {
		return certPool
	}
	rootCertPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile(viper.GetString("ca_cert_path"))
	if err != nil {
		zap.S().Warn(err.Error())
		return nil
	}

	if ok := rootCertPool.AppendCertsFromPEM([]byte(strings.TrimSpace(string(pem)))); !ok {
		zap.S().Warn("failed to append pem")
		return nil
	}
	certPool = rootCertPool
	return certPool
}
