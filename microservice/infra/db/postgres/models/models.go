package models

import "time"

type OrderModel struct {
	ID         string           `gorm:"primaryKey;size:36"`
	CustomerID *string          `gorm:"size:36"`
	Amount     float64          `gorm:"not null"`
	StatusID   string           `gorm:"not null;size:36"`
	Status     OrderStatusModel `gorm:"foreignKey:StatusID;references:ID"`
	Items      []OrderItemModel `gorm:"foreignKey:OrderID;references:ID"`
	CreatedAt  time.Time        `gorm:"not null"`
	UpdatedAt  *time.Time
}

func (OrderModel) TableName() string {
	return "orders"
}

type OrderItemModel struct {
	ID        string  `gorm:"primaryKey;size:36"`
	OrderID   string  `gorm:"not null;size:36"`
	ProductID string  `gorm:"not null;size:36"`
	Quantity  int     `gorm:"not null"`
	UnitPrice float64 `gorm:"not null"`
}

func (OrderItemModel) TableName() string {
	return "order_items"
}

type OrderStatusModel struct {
	ID   string `gorm:"primaryKey;size:36"`
	Name string `gorm:"not null;size:100"`
}

func (OrderStatusModel) TableName() string {
	return "order_status"
}
