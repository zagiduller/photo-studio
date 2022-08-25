package users

import (
	"fmt"
	"gorm.io/gorm"
	"photostudio/components"
)

// @project photo-studio
// @created 10.08.2022

type Service struct {
	components.Default
	db *gorm.DB
}

func New() *Service {
	return &Service{
		Default: components.New("users"),
	}
}
func (s *Service) Configure() error {
	s.db = components.GetDB()
	if s.db == nil {
		return fmt.Errorf("users.Configure: %w ", components.ErrorCodeDbIsNil)
	}
	// migrate model
	if err := s.db.AutoMigrate(&User{}); err != nil {
		return fmt.Errorf("users.Configure: %w ", err)
	}
	return nil
}

func (s *Service) NewUser() (*User, error) {
	user := &User{
		db:     s.db,
		Status: UserStatusActive,
	}

	return user, nil
}
