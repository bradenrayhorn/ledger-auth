package config

import (
	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic("Failed to read config: " + err.Error())
	}

	// mysql
	_ = viper.BindEnv("mysql_host", "MYSQL_HOST")
	_ = viper.BindEnv("mysql_port", "MYSQL_PORT")
	_ = viper.BindEnv("mysql_username", "MYSQL_USERNAME")
	_ = viper.BindEnv("mysql_password", "MYSQL_PASSWORD")
	_ = viper.BindEnv("mysql_database", "MYSQL_DATABASE")
}
