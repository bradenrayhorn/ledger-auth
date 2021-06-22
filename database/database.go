package database

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sqlx.DB
var RDB *redis.Client

var tlsConfig *tls.Config

func Setup() {
	tlsString := ""
	if tls := loadCACert(); tls != nil {
		mysql.RegisterTLSConfig("vault", tls)
		tlsString = "&tls=vault"
	}

	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true%s",
		viper.GetString("mysql_username"),
		viper.GetString("mysql_password"),
		viper.GetString("mysql_host"),
		viper.GetString("mysql_port"),
		viper.GetString("mysql_database"),
		tlsString,
	))

	if err != nil {
		zap.S().Error(err.Error())
	}

	DB = db
}

func SetupRedis() {
	options := &redis.Options{
		Addr:     viper.GetString("redis_addr"),
		Password: viper.GetString("redis_password"),
		DB:       viper.GetInt("redis_db"),
		Username: viper.GetString("redis_username"),
	}

	if tls := loadCACert(); tls != nil {
		options.TLSConfig = tls
	}

	rdb := redis.NewClient(options)

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		zap.S().Error(err)
	}

	RDB = rdb
}

func loadCACert() *tls.Config {
	if tlsConfig != nil {
		return tlsConfig
	}
	if !viper.GetBool("use_database_tls") {
		return nil
	}
	rootCertPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile(viper.GetString("ca_cert_path"))
	if err != nil {
		zap.S().Warn(err.Error())
		return nil
	}

	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		zap.S().Warn("failed to append pem")
		return nil
	}

	tlsConfig = &tls.Config{
		RootCAs: rootCertPool,
	}
	return tlsConfig
}
