package dtos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOrderDTO_Structure(t *testing.T) {
	orderDTO := OrderDTO{
		ID:         "order-123",
		CustomerID: stringPtr("customer-456"),
		Amount:     99.99,
		Items: []OrderItemDTO{
			{
				ID:        "item-1",
				ProductID: "product-1",
				OrderID:   "order-123",
				Quantity:  2,
				UnitPrice: 49.99,
			},
		},
	}

	assert.Equal(t, "order-123", orderDTO.ID)
	assert.Equal(t, "customer-456", *orderDTO.CustomerID)
	assert.Equal(t, 99.99, orderDTO.Amount)
	assert.Len(t, orderDTO.Items, 1)
	assert.Equal(t, "item-1", orderDTO.Items[0].ID)
}

func TestOrderItemDTO_Structure(t *testing.T) {
	itemDTO := OrderItemDTO{
		ID:        "item-123",
		ProductID: "product-789",
		OrderID:   "order-456",
		Quantity:  3,
		UnitPrice: 25.50,
	}

	assert.Equal(t, "item-123", itemDTO.ID)
	assert.Equal(t, "product-789", itemDTO.ProductID)
	assert.Equal(t, "order-456", itemDTO.OrderID)
	assert.Equal(t, 3, itemDTO.Quantity)
	assert.Equal(t, 25.50, itemDTO.UnitPrice)
}

func TestCreateOrderDTO_Structure(t *testing.T) {
	createOrderDTO := CreateOrderDTO{
		CustomerID: stringPtr("customer-123"),
		Items: []CreateOrderItemDTO{
			{
				ProductID: "product-1",
				Quantity:  2,
				Price:     15.99,
			},
		},
	}

	assert.Equal(t, "customer-123", *createOrderDTO.CustomerID)
	assert.Len(t, createOrderDTO.Items, 1)
	assert.Equal(t, "product-1", createOrderDTO.Items[0].ProductID)
	assert.Equal(t, 2, createOrderDTO.Items[0].Quantity)
	assert.Equal(t, 15.99, createOrderDTO.Items[0].Price)
}

func TestCreateOrderItemDTO_Structure(t *testing.T) {
	createItemDTO := CreateOrderItemDTO{
		ProductID: "product-456",
		Quantity:  5,
		Price:     12.50,
	}

	assert.Equal(t, "product-456", createItemDTO.ProductID)
	assert.Equal(t, 5, createItemDTO.Quantity)
	assert.Equal(t, 12.50, createItemDTO.Price)
}

func TestUpdateOrderDTO_Structure(t *testing.T) {
	updateOrderDTO := UpdateOrderDTO{
		ID:       "order-789",
		StatusID: "status-2",
	}

	assert.Equal(t, "order-789", updateOrderDTO.ID)
	assert.Equal(t, "status-2", updateOrderDTO.StatusID)
}

func TestOrderFilterDTO_Structure(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	filterDTO := OrderFilterDTO{
		CreatedAtFrom: &yesterday,
		CreatedAtTo:   &now,
		StatusID:      stringPtr("status-1"),
		CustomerID:    stringPtr("customer-123"),
	}

	assert.Equal(t, yesterday, *filterDTO.CreatedAtFrom)
	assert.Equal(t, now, *filterDTO.CreatedAtTo)
	assert.Equal(t, "status-1", *filterDTO.StatusID)
	assert.Equal(t, "customer-123", *filterDTO.CustomerID)
}

func TestOrderFilterDTO_WithNilValues(t *testing.T) {
	filterDTO := OrderFilterDTO{
		CreatedAtFrom: nil,
		CreatedAtTo:   nil,
		StatusID:      nil,
		CustomerID:    nil,
	}

	assert.Nil(t, filterDTO.CreatedAtFrom)
	assert.Nil(t, filterDTO.CreatedAtTo)
	assert.Nil(t, filterDTO.StatusID)
	assert.Nil(t, filterDTO.CustomerID)
}

func TestOrderStatusDTO_Structure(t *testing.T) {
	statusDTO := OrderStatusDTO{
		ID:   "status-123",
		Name: "confirmed",
	}

	assert.Equal(t, "status-123", statusDTO.ID)
	assert.Equal(t, "confirmed", statusDTO.Name)
}

func TestOrderResponseDTO_Structure(t *testing.T) {
	now := time.Now()
	updatedAt := now.Add(time.Hour)

	responseDTO := OrderResponseDTO{
		ID:         "order-123",
		CustomerID: stringPtr("customer-456"),
		Amount:     150.75,
		Status: OrderStatusDTO{
			ID:   "status-1",
			Name: "pending",
		},
		Items: []OrderItemDTO{
			{
				ID:        "item-1",
				ProductID: "product-1",
				OrderID:   "order-123",
				Quantity:  3,
				UnitPrice: 50.25,
			},
		},
		CreatedAt: now,
		UpdatedAt: &updatedAt,
	}

	assert.Equal(t, "order-123", responseDTO.ID)
	assert.Equal(t, "customer-456", *responseDTO.CustomerID)
	assert.Equal(t, 150.75, responseDTO.Amount)
	assert.Equal(t, "status-1", responseDTO.Status.ID)
	assert.Equal(t, "pending", responseDTO.Status.Name)
	assert.Len(t, responseDTO.Items, 1)
	assert.Equal(t, now, responseDTO.CreatedAt)
	assert.Equal(t, updatedAt, *responseDTO.UpdatedAt)
}

func TestOrderStatusResponseDTO_Structure(t *testing.T) {
	statusResponseDTO := OrderStatusResponseDTO{
		ID:   "status-456",
		Name: "delivered",
	}

	assert.Equal(t, "status-456", statusResponseDTO.ID)
	assert.Equal(t, "delivered", statusResponseDTO.Name)
}

func TestCreateOrderDTO_WithMultipleItems(t *testing.T) {
	items := []CreateOrderItemDTO{
		{
			ProductID: "product-1",
			Quantity:  1,
			Price:     10.00,
		},
		{
			ProductID: "product-2",
			Quantity:  2,
			Price:     15.00,
		},
	}

	createOrderDTO := CreateOrderDTO{
		CustomerID: stringPtr("customer-789"),
		Items:      items,
	}

	assert.Equal(t, "customer-789", *createOrderDTO.CustomerID)
	assert.Len(t, createOrderDTO.Items, 2)
	assert.Equal(t, "product-1", createOrderDTO.Items[0].ProductID)
	assert.Equal(t, "product-2", createOrderDTO.Items[1].ProductID)
}

func TestCreateOrderDTO_WithNilCustomer(t *testing.T) {
	createOrderDTO := CreateOrderDTO{
		CustomerID: nil,
		Items: []CreateOrderItemDTO{
			{
				ProductID: "product-1",
				Quantity:  1,
				Price:     10.00,
			},
		},
	}

	assert.Nil(t, createOrderDTO.CustomerID)
	assert.Len(t, createOrderDTO.Items, 1)
}

func stringPtr(s string) *string {
	return &s
}