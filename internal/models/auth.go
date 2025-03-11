package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Identifier string   `json:"identifier" gorm:"unique"`
	Password   string   `json:"password"`
	IsActive   bool     `gorm:"default: true" json:"is_active"`
	IsAdmin    bool     `gorm:"default: false" json:"is_admin"`
	Metadata   string   `gorm:"default:'{}'" json:"metadata"`
	Groups     []*Group `gorm:"many2many:group_users;"`
}

type Group struct {
	gorm.Model
	Identifier  string  `json:"identifier" gorm:"unique"`
	Description string  `gorm:"default:''" json:"description"`
	Metadata    string  `gorm:"default:'{}'" json:"metadata"`
	IsActive    bool    `gorm:"default: true" json:"is_active"`
	Users       []*User `gorm:"many2many:group_users;"`
}

type Client struct {
	gorm.Model
	Identifier string `json:"identifier" gorm:"unique"`
	Secret     string `json:"secret"`
	IsActive   bool   `gorm:"default: true" json:"is_active"`
	Grant      string `json:"grant"`
}

type Token struct {
	gorm.Model
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	UserID       uint   `json:"user_id"`
	ClientID     uint   `json:"client_id"`

	User   User   `json:"user"`
	Client Client `json:"client"`
}
