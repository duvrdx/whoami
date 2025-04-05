package models

import (
	"gorm.io/gorm"
)

type Config struct {
	gorm.Model
	Key   string `json:"key" gorm:"unique"`
	Value string `json:"value"`
}

type Superuser struct {
	gorm.Model
	Identifier string `json:"identifier" gorm:"unique"`
	Password   string `json:"password"`
	IsActive   bool   `gorm:"default: true" json:"is_active"`
}
