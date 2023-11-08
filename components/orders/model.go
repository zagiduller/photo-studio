package orders

import (
	"errors"
	"fmt"
	"github.com/zagiduller/photo-studio/components"
	"github.com/zagiduller/photo-studio/components/users"
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
	components.Model

	Status       OrderStatus `gorm:"type:varchar(16)" json:"status"`
	User         *users.User `gorm:"foreignKey:UserID" json:"user"`
	UserID       uint        `json:"user_id"`
	ManagerID    uint        `json:"manager_id"`
	Description  string      `gorm:"type:varchar(255)" json:"description"`
	CustomerName string      `gorm:"type:varchar(255)" json:"customer_name"`

	//Participants []
}

var (
	ValidateErrorOwnerIsNil              = errors.New("orders: owner is nil")
	ValidateErrorCodeOrderIsNil          = errors.New("orders: order is nil")
	ValidateErrorCodeUnsupportedStatus   = errors.New("orders: status is not supported")
	ValidateErrorCodeEmailOrNameRequired = errors.New("orders: email or name is required")
)

func (o *Order) Validate() error {
	if o == nil {
		return fmt.Errorf("Order.Validate: [%w]", ValidateErrorCodeOrderIsNil)
	}
	if o.GetDB() == nil {
		return fmt.Errorf("Order.Validate: [%w]", components.ErrorCodeDbIsNil)
	}
	if !o.CheckStatusIsValid() {
		return fmt.Errorf("Order.Validate: [%w]", ValidateErrorCodeUnsupportedStatus)
	}

	return nil
}

func (o *Order) CheckStatusIsValid() bool {
	for _, sups := range supportedStatuses {
		if o.Status == sups {
			return true
		}
	}
	return false
}

func (o *Order) Save() error {
	if err := o.Validate(); err != nil {
		return fmt.Errorf("Order.Save: [%w] ", err)
	}
	if err := o.GetDB().Save(o).Error; err != nil {
		return fmt.Errorf("Order.Save: [%w] ", err)
	}
	return nil
}
