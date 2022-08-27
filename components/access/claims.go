package access

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/zagiduller/photo-studio/components/users"
)

// @project photo-studio
// @created 11.08.2022

type Claims struct {
	*jwt.RegisteredClaims
	User users.Claims
}
