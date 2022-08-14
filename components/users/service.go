package users

import (
	"gorm.io/gorm"
)

// @project photo-studio
// @created 10.08.2022

type Service struct {
	db *gorm.DB
}

func New() *Service {
	return &Service{}
}
func (s *Service) Configure() error {

	return nil
}

func (s *Service) NewUser() (*User, error) {
	user := &User{
		db:     s.db,
		Status: UserStatusActive,
	}

	return user, nil
}
