package main

import (
	"log"

	"github.com/duvrdx/whoami/internal/config"
	"github.com/duvrdx/whoami/internal/models"
	"github.com/duvrdx/whoami/internal/routing"

	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigName("config") // Nome do arquivo (sem extensão)
	viper.SetConfigType("yaml")   // Formato do arquivo
	viper.AddConfigPath(".")

	// Tenta ler o arquivo de configuração
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "admin")
	viper.SetDefault("database.password", "secret")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("token.jwt-secret", "youshoudbesecret")
	viper.SetDefault("token.expiration", 3600)

	config.Connect()
	config.MigrateDB(models.User{}, models.Client{}, models.Group{}, models.Token{})

	e := routing.Routing.GetRoutes(routing.Routing{})

	e.Logger.Fatal(e.Start(":8080"))
}
