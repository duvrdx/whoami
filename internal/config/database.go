package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPostgres() {

	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		Config.Database.Host, Config.Database.User, Config.Database.Password, Config.Database.Name, Config.Database.Port)

	database, err := gorm.Open(
		postgres.Open(connectionString),
		&gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	DB = database
}

func ConnectSQLite() {
	database, err := gorm.Open(sqlite.Open("./whoami.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	DB = database
}

func Connect() {
	if Config.Database.Dialect == "postgres" {
		ConnectPostgres()
	} else if Config.Database.Dialect == "sqlite" {
		ConnectSQLite()
	} else {
		panic("Invalid database dialect!")
	}
}

func GetDB() *gorm.DB {
	return DB
}

func MigrateDB(models ...interface{}) {
	DB.AutoMigrate(models...)
}
