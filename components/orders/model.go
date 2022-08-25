package orders

import (
	"errors"
	"gorm.io/gorm"
	"photostudio/components"
	"photostudio/components/users"
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
	db          *gorm.DB
	Description string      `gorm:"type:varchar(255)" json:"description"`
	Status      OrderStatus `gorm:"type:varchar(16)" json:"status"`
	Manager     *users.User `json:"manager"`
}

var (
	ValidateErrorOwnerIsNil            = errors.New("orders: owner is nil")
	ValidateErrorCodeOrderIsNil        = errors.New("orders: order is nil")
	ValidateErrorCodeUnsupportedStatus = errors.New("orders: status is not supported")
)

func (o *Order) Validate() error {
	if o == nil {
		return ValidateErrorCodeOrderIsNil
	}
	if o.db == nil {
		return components.ErrorCodeDbIsNil
	}
	if !o.CheckStatusIsValid() {
		return ValidateErrorCodeUnsupportedStatus
	}
	//if o.Owner == nil {
	//	return ValidateErrorOwnerIsNil
	//}

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
