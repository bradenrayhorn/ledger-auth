package config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Failed to read config: " + err.Error())
	}

	viper.SetDefault("cookie_secure", true)
	viper.SetDefault("allow_credentials", false)
	viper.SetDefault("session_duration", "24h")
	viper.SetDefault("grpc_port", "9000")

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
}
