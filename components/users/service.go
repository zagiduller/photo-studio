package users

import (
	"context"
	"fmt"
	"github.com/zagiduller/photo-studio/components"
	"gorm.io/gorm"
)

// @project photo-studio
// @created 10.08.2022

type Service struct {
	components.Default
	db *gorm.DB
}

func New() *Service {
	return &Service{
		Default: components.DefaultComponent("users"),
	}
}

func (s *Service) Configure(ctx context.Context) error {
	s.Default.Ctx = ctx
	s.db = components.GetDB()
	if s.db == nil {
		return fmt.Errorf("users.Configure: [%w] ", components.ErrorCodeDbIsNil)
	}
	// migrate model
	if err := s.db.AutoMigrate(&User{}); err != nil {
		return fmt.Errorf("users.Configure: [%w] ", err)
	}
	return nil
}

func (s *Service) NewUser() (*User, error) {
	user := &User{
		Status: UserStatusActive,
		Role:   RoleUser,
	}

	return user, nil
}

func (s *Service) FindByLogin(login string) (*User, error) {
	u := &User{}
	if err := s.db.Model(&User{}).Where("login = ?", login).First(u).Error; err != nil {
		return nil, fmt.Errorf("FindByLogin: [%w] ", err)
	}
	if u == nil {
		return nil, fmt.Errorf("FindByLogin: [%w] ", components.ErrModelNotFound)
	}
	return u, nil
}
