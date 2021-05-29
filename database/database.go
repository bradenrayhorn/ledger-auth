package database

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sqlx.DB
var RDB *redis.Client

func Setup() {
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		viper.GetString("mysql_username"),
		viper.GetString("mysql_password"),
		viper.GetString("mysql_host"),
		viper.GetString("mysql_port"),
		viper.GetString("mysql_database"),
	))

	if err != nil {
		zap.S().Panic(err.Error())
	}

	DB = db
}

func SetupRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis_addr"),
		Password: viper.GetString("redis_password"),
		DB:       viper.GetInt("redis_db"),
		Username: viper.GetString("redis_username"),
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		zap.S().Panic(err)
	}

	RDB = rdb
}
