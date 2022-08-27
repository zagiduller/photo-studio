package auth

import (
	"gorm.io/gorm"
	"photostudio/components/users"
)

// @project photo-studio
// @created 10.08.2022

type Auth struct {
	gorm.Model
	db *gorm.DB

	Token  string
	UserID uint
	User   users.User
}
