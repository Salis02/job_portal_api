package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv  string
	AppPort string
	DB DBConfig
}

type DBConfig struct {
	Host string
	User string
	Pass string
	Name string
	Port string
	SSL string
}

func LoadConfig()(*Config, error){
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error loading .env files:", err)
	}

	config := &Config{
		AppEnv: viper.GetString("APP_ENV"),
		AppPort: viper.GetString("APP_PORT"),
		DB: DBConfig{
			Host: viper.GetString("DB_HOST"),
			User: viper.GetString("DB_USER"),
			Pass: viper.GetString("DB_PASS"),
			Name: viper.GetString("DB_NAME"),
			Port: viper.GetString("DB_PORT"),
			SSL: viper.GetString("DB_SSL"),
		},
	}

	return config, nil
}
