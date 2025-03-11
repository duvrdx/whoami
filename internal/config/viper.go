package config

import (
	"log"

	"github.com/spf13/viper"
)

var JWTSecret []byte
var TokenExpiration int

// Inicializa o Viper e carrega as configurações
func InitConfig() {
	// Configurações padrão
	viper.SetDefault("token.secret", "defaultSecret")
	viper.SetDefault("token.expiration", 3600)

	// Configurações do Viper
	viper.SetConfigName("config") // Nome do arquivo (sem extensão)
	viper.SetConfigType("yaml")   // Formato do arquivo
	viper.AddConfigPath(".")      // Diretório atual

	// Tenta ler o arquivo de configuração
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Failed to read config file: %v", err)
	}

	JWTSecret = []byte(viper.GetString("token.secret"))
	TokenExpiration = viper.GetInt("token.expiration")

}

func GetJWTSecret() []byte {
	InitConfig()
	return JWTSecret
}
