package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)
import _ "github.com/go-sql-driver/mysql"

var DB *sqlx.DB

func Setup() {
	fmt.Printf(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		viper.GetString("mysql_username"),
		viper.GetString("mysql_password"),
		viper.GetString("mysql_host"),
		viper.GetString("mysql_port"),
		viper.GetString("mysql_database"),
	))
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
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
