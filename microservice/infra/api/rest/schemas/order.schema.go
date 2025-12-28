package schemas

import "time"

type CreateOrderItemSchema struct {
	ProductID string  `json:"product_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
	Price     float64 `json:"price" binding:"required,min=0"`
}

type CreateOrderSchema struct {
	CustomerID *string                 `json:"customer_id"`
	Items      []CreateOrderItemSchema `json:"items" binding:"required,min=1"`
}

type UpdateOrderSchema struct {
	StatusID string `json:"status_id" binding:"required"`
}

type OrderItemResponseSchema struct {
	ID        string  `json:"id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

type OrderResponseSchema struct {
	ID         string                    `json:"id"`
	CustomerID *string                   `json:"customer_id"`
	Amount     float64                   `json:"amount"`
	Status     string                    `json:"status"`
	Items      []OrderItemResponseSchema `json:"items"`
	CreatedAt  time.Time                 `json:"created_at"`
	UpdatedAt  *time.Time                `json:"updated_at"`
}

type OrderStatusResponseSchema struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
