package orders

import (
	"errors"
	"fmt"
	"github.com/zagiduller/photo-studio/components"
	"gorm.io/gorm"
)

// @project photo-studio
// @created 10.08.2022

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "new"
	OrderStatusInProgress OrderStatus = "in_progress"
	OrderStatusPaid       OrderStatus = "paid"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

var supportedStatuses = []OrderStatus{
	OrderStatusNew,
	OrderStatusInProgress,
	OrderStatusPaid,
	OrderStatusCancelled,
}

type Order struct {
	gorm.Model
	db *gorm.DB

	Status       OrderStatus `gorm:"type:varchar(16)" json:"status"`
	ManagerID    uint        `json:"manager_id"`
	Description  string      `gorm:"type:varchar(255)" json:"description"`
	Email        string      `gorm:"type:varchar(255)" json:"email"`
	Phone        string      `gorm:"type:varchar(255)" json:"phone"`
	CustomerName string      `gorm:"type:varchar(255)" json:"customer_name"`
}

var (
	ValidateErrorOwnerIsNil              = errors.New("orders: owner is nil")
	ValidateErrorCodeOrderIsNil          = errors.New("orders: order is nil")
	ValidateErrorCodeUnsupportedStatus   = errors.New("orders: status is not supported")
	ValidateErrorCodeEmailOrNameRequired = errors.New("orders: email or name is required")
)

func (o *Order) Validate() error {
	if o == nil {
		return ValidateErrorCodeOrderIsNil
	}
	if o.db == nil {
		return components.ErrorCodeDbIsNil
	}
	if o.Email == "" && o.Phone == "" {
		return ValidateErrorCodeEmailOrNameRequired
	}
	if !o.CheckStatusIsValid() {
		return ValidateErrorCodeUnsupportedStatus
	}

	return nil
}

func (o *Order) CheckStatusIsValid() bool {
	for _, supportedStatuses := range supportedStatuses {
		if o.Status == supportedStatuses {
			return true
		}
	}
	return false
}

func (o *Order) Save() error {
	if err := o.Validate(); err != nil {
		return fmt.Errorf("orders.Save: %w ", err)
	}
	if err := o.db.Save(o).Error; err != nil {
		return fmt.Errorf("orders.Save: %w ", err)
	}
	return nil
}
