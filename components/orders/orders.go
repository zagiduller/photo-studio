package orders

import (
	"gorm.io/gorm"
	"photostudio/components/users"
	"time"
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

type Order struct {
	gorm.Model
	db *gorm.DB

	Status    OrderStatus `json:"status"`
	Owner     *users.User `json:"owner"`
	Manager   *users.User `json:"manager"`
	CreatedAt time.Time   `json:"created_at"`
}
