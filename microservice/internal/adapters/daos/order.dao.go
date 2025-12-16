package daos

import "time"

type OrderDAO struct {
	ID         string
	CustomerID *string
	Amount     float64
	Status     OrderStatusDAO
	Items      []OrderItemDAO
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}

type OrderItemDAO struct {
	ID        string
	OrderID   string
	ProductID string
	Quantity  int
	UnitPrice float64
}

type OrderStatusDAO struct {
	ID   string
	Name string
}
