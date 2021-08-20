package database

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var DB *sqlx.DB
var RDB *redis.Client

var tlsConfig *tls.Config

func Setup() {
	connConfig, err := pgx.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		viper.GetString("pg_username"),
		viper.GetString("pg_password"),
		viper.GetString("pg_host"),
		viper.GetString("pg_port"),
		viper.GetString("pg_database"),
	))
	if err != nil {
		zap.S().Error(err.Error())
		return
	}
	connConfig.TLSConfig = loadCACert()
	if connConfig.TLSConfig != nil {
		connConfig.TLSConfig.ServerName = "10.0.1.2"
	}
	connString := stdlib.RegisterConnConfig(connConfig)
	db, err := sqlx.Open("pgx", connString)
	if err != nil {
		zap.S().Error(err.Error())
		return
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
