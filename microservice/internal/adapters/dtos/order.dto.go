package dtos

import "time"

type OrderDTO struct {
	ID         string
	CustomerID *string
	Amount     float64
	Items      []OrderItemDTO
}

type OrderItemDTO struct {
	ID        string
	ProductID string
	OrderID   string
	Quantity  int
	UnitPrice float64
}

type CreateOrderDTO struct {
	CustomerID *string
	Items      []CreateOrderItemDTO
}

type CreateOrderItemDTO struct {
	ProductID string
	Quantity  int
	Price     float64
}

type UpdateOrderDTO struct {
	ID       string
	StatusID string
}

type OrderFilterDTO struct {
	CreatedAtFrom *time.Time
	CreatedAtTo   *time.Time
	StatusID      *string
	CustomerID    *string
}

type OrderStatusDTO struct {
	ID   string
	Name string
}

type OrderResponseDTO struct {
	ID         string
	CustomerID *string
	Amount     float64
	Status     OrderStatusDTO
	Items      []OrderItemDTO
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}

type OrderStatusResponseDTO struct {
	ID   string
	Name string
}
