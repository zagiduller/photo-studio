package orders

import (
	"context"
	"fmt"
	"github.com/zagiduller/photo-studio/components"
	"gorm.io/gorm"
)

// @project photo-studio
// @created 14.08.2022

type Service struct {
	components.Default
	db *gorm.DB
}

func New() *Service {
	return &Service{
		Default: components.DefaultComponent("orders"),
	}
}

func (s *Service) Configure(ctx context.Context) error {
	s.Default.Ctx = ctx
	s.db = components.GetDB()
	if s.db == nil {
		return fmt.Errorf("orders.Configure: %w ", components.ErrorCodeDbIsNil)
	}
	// migrate model
	if err := s.db.AutoMigrate(&Order{}); err != nil {
		return fmt.Errorf("orders.Configure: %w ", err)
	}
	return nil
}

func (s *Service) GetAll() ([]*Order, error) {
	var orders []*Order
	if err := s.db.Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("orders.GetAll: %w ", err)
	}
	return orders, nil
}

func (s *Service) Create(description string) (*Order, error) {
	order := &Order{
		db:          s.db,
		Description: description,
		Status:      OrderStatusNew,
	}
	if err := order.Validate(); err != nil {
		return nil, fmt.Errorf("orders.Create: %w ", err)
	}
	if err := order.Save(); err != nil {
		return nil, fmt.Errorf("orders.Create: %w ", err)
	}
	return order, nil
}
