package users

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"photostudio/components"
)

// @project photo-studio
// @created 10.08.2022

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"

	RoleAdmin    = "admin"
	RoleCustomer = "customer"
)

// User was described
type User struct {
	gorm.Model
	db *gorm.DB

	Status UserStatus `gorm:"type:varchar(12)" json:"status"`
	Role   string     `gorm:"type:varchar(12)" json:"role"`
	Login  string     `gorm:"unique;type:varchar(255)" json:"login"`
	Email  string     `gorm:"unique;type:varchar(255)" json:"email"`
}

var (
	ValidateErrorCodeInvalidLogin  = errors.New("Login is invalid ")
	ValidateErrorCodeInvalidEmail  = errors.New("Email is invalid ")
	ValidateErrorCodeInvalidStatus = errors.New("Status is invalid ")
	ValidateErrorNilUser           = errors.New("User is nil ")
	ValidateErrorCodeInvalidRole   = errors.New("Role is invalid ")

	ErrorCodeUserNotFound = errors.New("User not found ")
	ErrorCodeUserExists   = errors.New("User already exists ")
)

func (u *User) Validate() error {
	if u == nil {
		return fmt.Errorf("Validate: %w ", ValidateErrorNilUser)
	}
	if u.db == nil {
		return fmt.Errorf("Validate: %w ", components.ErrorCodeDbIsNil)
	}
	if u.Login == "" {
		return fmt.Errorf("Validate: %w ", ValidateErrorCodeInvalidLogin)
	}
	if u.Status == "" {
		return fmt.Errorf("Validate: %w ", ValidateErrorCodeInvalidStatus)
	}
	if u.Role == "" {
		return fmt.Errorf("Validate: %w ", ValidateErrorCodeInvalidRole)
	}
	return nil
}

func (u *User) Save() error {
	if err := u.Validate(); err != nil {
		return fmt.Errorf("Save: %w ", err)
	}
	if err := u.db.Save(u).Error; err != nil {
		return fmt.Errorf("Save: %w ", err)
	}
	log.Infof("Save: User saved: %d %s ", u.ID, u.Login)
	return nil
}

func (u *User) GetClaims() Claims {
	return Claims{
		ID:   u.ID,
		Role: u.Role,
	}
}
