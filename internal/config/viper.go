package config

import (
	"log"

	"github.com/spf13/viper"
)

type TokenConfig struct {
	Secret     []byte
	Expiration int
}

type DatabaseConfig struct {
	Dialect  string
	Host     string
	Port     string
	User     string
	Name     string
	Password string
}

type AppConfig struct {
	Token    TokenConfig
	Database DatabaseConfig
}

var Config AppConfig

func Init() {
	// Configurações padrão
	viper.SetDefault("token.secret", "defaultSecret")
	viper.SetDefault("token.expiration", 3600)
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.dialect", "sqlite")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Failed to read config file: %v", err)
	}

	// Carrega todas as configurações na struct
	Config = AppConfig{
		Token: TokenConfig{
			Secret:     []byte(viper.GetString("token.secret")),
			Expiration: viper.GetInt("token.expiration"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("database.host"),
			Port:     viper.GetString("database.port"),
			Dialect:  viper.GetString("database.dialect"),
			User:     viper.GetString("database.user"),
			Name:     viper.GetString("database.name"),
			Password: viper.GetString("database.password"),
		},
	}
}

func GetConfig() AppConfig {
	return Config
}
