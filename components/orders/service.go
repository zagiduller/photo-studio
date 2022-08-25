package orders

import (
	"gorm.io/gorm"
	"photostudio/components"
)

// @project photo-studio
// @created 14.08.2022

type Service struct {
	db *gorm.DB
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Configure() error {
	s.db = components.GetDB()
	return nil
}
