package auth

import (
	"github.com/zagiduller/photo-studio/components/users"
	"gorm.io/gorm"
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
