package users

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zagiduller/photo-studio/components"
)

// @project photo-studio
// @created 10.08.2022

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"

	RoleAdmin = "admin"
	RoleUser  = "user"
)

// User was described
type User struct {
	components.Model

	Status UserStatus `gorm:"type:varchar(12)" json:"status"`
	Role   string     `gorm:"type:varchar(12)" json:"role"`
}

var (
	ValidateErrorCodeInvalidStatus = errors.New("Status is invalid ")
	ValidateErrorNilUser           = errors.New("User is nil ")
	ValidateErrorCodeInvalidRole   = errors.New("Role is invalid ")

	ErrorCodeUserNotFound = errors.New("User not found ")
	ErrorCodeUserExists   = errors.New("User already exists ")
)

func (u *User) Validate() error {
	if u == nil {
		return fmt.Errorf("User.Validate: [%w] ", ValidateErrorNilUser)
	}
	if u.GetDB() == nil {
		return fmt.Errorf("User.Validate: [%w] ", components.ErrorCodeDbIsNil)
	}
	if u.Status == "" {
		return fmt.Errorf("User.Validate: [%w] ", ValidateErrorCodeInvalidStatus)
	}
	if u.Role == "" {
		return fmt.Errorf("User.Validate: [%w] ", ValidateErrorCodeInvalidRole)
	}
	return nil
}

func (u *User) Save() error {
	if err := u.Validate(); err != nil {
		return fmt.Errorf("User.Save: [%w] ", err)
	}
	if err := u.GetDB().Save(u).Error; err != nil {
		return fmt.Errorf("User.Save: [%w] ", err)
	}
	log.Infof("User.Save: id(%d) ", u.ID)
	return nil
}

func (u *User) GetClaims() Claims {
	return Claims{
		ID:   u.ID,
		Role: u.Role,
	}
}
