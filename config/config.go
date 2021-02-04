package config

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var RsaPrivate *rsa.PrivateKey
var RsaPublic *rsa.PublicKey

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Failed to read config: " + err.Error())
	}

	viper.SetDefault("token_expiration", time.Hour*24)
	viper.SetDefault("rsa_path", "jwt_rsa")

	// mysql
	_ = viper.BindEnv("mysql_host", "MYSQL_HOST")
	_ = viper.BindEnv("mysql_port", "MYSQL_PORT")
	_ = viper.BindEnv("mysql_username", "MYSQL_USERNAME")
	_ = viper.BindEnv("mysql_password", "MYSQL_PASSWORD")
	_ = viper.BindEnv("mysql_database", "MYSQL_DATABASE")
	// other
	_ = viper.BindEnv("rsa_path", "RSA_PATH")

	loadRsaKeys()
}

func loadRsaKeys() {
	privateKey, err := readKey(false)
	if err != nil {
		log.Fatalf("failed to load private rsa key: %s", err)
	}
	publicKey, err := readKey(true)
	if err != nil {
		log.Fatalf("failed to load public rsa key: %s", err)
	}
	RsaPrivate = privateKey.(*rsa.PrivateKey)
	RsaPublic = publicKey.(*rsa.PublicKey)
}

func readKey(public bool) (interface{}, error) {
	filePath := viper.GetString("rsa_path")
	if public {
		filePath += ".pub"
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	keyBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var rsaKey interface{}
	if public {
		rsaKey, err = jwt.ParseRSAPublicKeyFromPEM(keyBytes)
	} else {
		rsaKey, err = jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	}
	return rsaKey, err
}
