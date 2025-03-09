package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	database, err := gorm.Open(sqlite.Open("./test.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	DB = database
}

func GetDB() *gorm.DB {
	return DB
}

func MigrateDB(models ...interface{}) {
	DB.AutoMigrate(models...)
}
