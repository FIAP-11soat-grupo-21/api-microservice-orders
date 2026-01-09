package daos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOrderDAO_Structure(t *testing.T) {
	orderDAO := OrderDAO{
		ID:         "order-123",
		CustomerID: stringPtr("customer-456"),
		Amount:     99.99,
		Status: OrderStatusDAO{
			ID:   "status-1",
			Name: "pending",
		},
		Items: []OrderItemDAO{
			{
				ID:        "item-1",
				OrderID:   "order-123",
				ProductID: "product-1",
				Quantity:  2,
				UnitPrice: 49.99,
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: timePtr(time.Now()),
	}

	assert.Equal(t, "order-123", orderDAO.ID)
	assert.Equal(t, "customer-456", *orderDAO.CustomerID)
	assert.Equal(t, 99.99, orderDAO.Amount)
	assert.Equal(t, "status-1", orderDAO.Status.ID)
	assert.Equal(t, "pending", orderDAO.Status.Name)
	assert.Len(t, orderDAO.Items, 1)
	assert.Equal(t, "item-1", orderDAO.Items[0].ID)
	assert.Equal(t, "order-123", orderDAO.Items[0].OrderID)
	assert.Equal(t, "product-1", orderDAO.Items[0].ProductID)
	assert.Equal(t, 2, orderDAO.Items[0].Quantity)
	assert.Equal(t, 49.99, orderDAO.Items[0].UnitPrice)
	assert.NotNil(t, orderDAO.CreatedAt)
	assert.NotNil(t, orderDAO.UpdatedAt)
}

func TestOrderItemDAO_Structure(t *testing.T) {
	itemDAO := OrderItemDAO{
		ID:        "item-123",
		OrderID:   "order-456",
		ProductID: "product-789",
		Quantity:  3,
		UnitPrice: 25.50,
	}

	assert.Equal(t, "item-123", itemDAO.ID)
	assert.Equal(t, "order-456", itemDAO.OrderID)
	assert.Equal(t, "product-789", itemDAO.ProductID)
	assert.Equal(t, 3, itemDAO.Quantity)
	assert.Equal(t, 25.50, itemDAO.UnitPrice)
}

func TestOrderStatusDAO_Structure(t *testing.T) {
	statusDAO := OrderStatusDAO{
		ID:   "status-123",
		Name: "confirmed",
	}

	assert.Equal(t, "status-123", statusDAO.ID)
	assert.Equal(t, "confirmed", statusDAO.Name)
}

func TestOrderDAO_WithNilValues(t *testing.T) {
	orderDAO := OrderDAO{
		ID:         "order-123",
		CustomerID: nil,
		Amount:     0.0,
		Status: OrderStatusDAO{
			ID:   "",
			Name: "",
		},
		Items:     []OrderItemDAO{},
		CreatedAt: time.Time{},
		UpdatedAt: nil,
	}

	assert.Equal(t, "order-123", orderDAO.ID)
	assert.Nil(t, orderDAO.CustomerID)
	assert.Equal(t, 0.0, orderDAO.Amount)
	assert.Equal(t, "", orderDAO.Status.ID)
	assert.Equal(t, "", orderDAO.Status.Name)
	assert.Empty(t, orderDAO.Items)
	assert.True(t, orderDAO.CreatedAt.IsZero())
	assert.Nil(t, orderDAO.UpdatedAt)
}

func TestOrderDAO_MultipleItems(t *testing.T) {
	items := []OrderItemDAO{
		{
			ID:        "item-1",
			OrderID:   "order-123",
			ProductID: "product-1",
			Quantity:  1,
			UnitPrice: 10.00,
		},
		{
			ID:        "item-2",
			OrderID:   "order-123",
			ProductID: "product-2",
			Quantity:  2,
			UnitPrice: 15.00,
		},
	}

	orderDAO := OrderDAO{
		ID:         "order-123",
		CustomerID: stringPtr("customer-456"),
		Amount:     40.00,
		Status: OrderStatusDAO{
			ID:   "status-1",
			Name: "pending",
		},
		Items:     items,
		CreatedAt: time.Now(),
		UpdatedAt: timePtr(time.Now()),
	}

	assert.Len(t, orderDAO.Items, 2)
	assert.Equal(t, "item-1", orderDAO.Items[0].ID)
	assert.Equal(t, "item-2", orderDAO.Items[1].ID)
	assert.Equal(t, 1, orderDAO.Items[0].Quantity)
	assert.Equal(t, 2, orderDAO.Items[1].Quantity)
}

func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}