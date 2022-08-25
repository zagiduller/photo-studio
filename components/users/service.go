package users

import (
	"gorm.io/gorm"
	"photostudio/components"
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
	s.db = components.GetDB()
	return nil
}

func (s *Service) NewUser() (*User, error) {
	user := &User{
		db:     s.db,
		Status: UserStatusActive,
	}

	return user, nil
}
