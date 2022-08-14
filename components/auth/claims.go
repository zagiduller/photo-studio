package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"photostudio/components/users"
)

// @project photo-studio
// @created 11.08.2022

type Claims struct {
	*jwt.RegisteredClaims
	User users.Claims
}
