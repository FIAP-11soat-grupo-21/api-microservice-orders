package entities

import (
	"time"

	"microservice/internal/domain/exceptions"
	"microservice/internal/domain/value_objects"
	identityUtils "microservice/utils/identity"
)

type Order struct {
	ID         string
	CustomerID *string
	Amount     value_objects.Amount
	Status     OrderStatus
	Items      []OrderItem
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}

func NewOrder(id string, customerID *string) (*Order, error) {
	return &Order{
		ID:         id,
		CustomerID: customerID,
		Amount:     value_objects.Amount{},
		Items:      []OrderItem{},
		CreatedAt:  time.Now(),
	}, nil
}

func NewOrderWithItems(id string, customerID *string, amount float64, status OrderStatus, items []OrderItem, createdAt time.Time, updatedAt *time.Time) (*Order, error) {
	order, _ := NewOrder(id, customerID)
	order.Items = items
	order.Status = status
	order.CreatedAt = createdAt
	order.UpdatedAt = updatedAt
	var err error
	order.Amount, err = value_objects.NewAmount(amount)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (o *Order) AddItem(item OrderItem) {
	o.Items = append(o.Items, item)
}

func (o *Order) CalcTotalAmount() error {
	total := 0.0
	for _, item := range o.Items {
		total += item.GetTotal()
	}
	amount, err := value_objects.NewAmount(total)
	if err != nil {
		return err
	}
	o.Amount = amount
	return nil
}

func (c *Order) IsEmpty() bool {
	return c.ID == ""
}

func ValidateID(id string) error {
	if !identityUtils.IsValidUUID(id) {
		return &exceptions.InvalidOrderDataException{
			Message: "Invalid order ID",
		}
	}
	return nil
}

func (o *Order) ToMap() map[string]interface{} {
	items := make([]map[string]interface{}, len(o.Items))
	for i, item := range o.Items {
		items[i] = item.ToMap()
	}

	return map[string]interface{}{
		"id":          o.ID,
		"customer_id": o.CustomerID,
		"amount":      o.Amount.Value(),
		"status":      o.Status.ToMap(),
		"items":       items,
		"created_at":  o.CreatedAt,
		"updated_at":  o.UpdatedAt,
	}
}
