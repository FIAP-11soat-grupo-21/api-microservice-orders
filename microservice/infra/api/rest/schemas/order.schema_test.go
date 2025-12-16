package schemas

import (
	"testing"
	"time"
)

func TestCreateOrderItemSchema_Structure(t *testing.T) {
	schema := CreateOrderItemSchema{
		ProductID: "product-123",
		Quantity:  2,
		Price:     25.50,
	}

	if schema.ProductID != "product-123" {
		t.Errorf("CreateOrderItemSchema.ProductID = %v, want product-123", schema.ProductID)
	}
	if schema.Quantity != 2 {
		t.Errorf("CreateOrderItemSchema.Quantity = %v, want 2", schema.Quantity)
	}
	if schema.Price != 25.50 {
		t.Errorf("CreateOrderItemSchema.Price = %v, want 25.50", schema.Price)
	}
}

func TestCreateOrderSchema_Structure(t *testing.T) {
	customerID := "customer-123"
	schema := CreateOrderSchema{
		CustomerID: &customerID,
		Items: []CreateOrderItemSchema{
			{ProductID: "product-1", Quantity: 1, Price: 10.0},
		},
	}

	if *schema.CustomerID != customerID {
		t.Errorf("CreateOrderSchema.CustomerID = %v, want %v", *schema.CustomerID, customerID)
	}
	if len(schema.Items) != 1 {
		t.Errorf("CreateOrderSchema.Items length = %v, want 1", len(schema.Items))
	}
}

func TestCreateOrderSchema_NilCustomerID(t *testing.T) {
	schema := CreateOrderSchema{
		CustomerID: nil,
		Items: []CreateOrderItemSchema{
			{ProductID: "product-1", Quantity: 1, Price: 10.0},
		},
	}

	if schema.CustomerID != nil {
		t.Errorf("CreateOrderSchema.CustomerID = %v, want nil", schema.CustomerID)
	}
}

func TestUpdateOrderSchema_Structure(t *testing.T) {
	schema := UpdateOrderSchema{
		StatusID: "status-2",
	}

	if schema.StatusID != "status-2" {
		t.Errorf("UpdateOrderSchema.StatusID = %v, want status-2", schema.StatusID)
	}
}

func TestOrderItemResponseSchema_Structure(t *testing.T) {
	schema := OrderItemResponseSchema{
		ID:        "item-1",
		ProductID: "product-1",
		Quantity:  2,
		UnitPrice: 10.0,
	}

	if schema.ID != "item-1" {
		t.Errorf("OrderItemResponseSchema.ID = %v, want item-1", schema.ID)
	}
	if schema.ProductID != "product-1" {
		t.Errorf("OrderItemResponseSchema.ProductID = %v, want product-1", schema.ProductID)
	}
	if schema.Quantity != 2 {
		t.Errorf("OrderItemResponseSchema.Quantity = %v, want 2", schema.Quantity)
	}
	if schema.UnitPrice != 10.0 {
		t.Errorf("OrderItemResponseSchema.UnitPrice = %v, want 10.0", schema.UnitPrice)
	}
}

func TestOrderResponseSchema_Structure(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()
	updatedAt := now.Add(time.Hour)

	schema := OrderResponseSchema{
		ID:         "order-1",
		CustomerID: &customerID,
		Amount:     100.0,
		Status:     "Pending",
		Items: []OrderItemResponseSchema{
			{ID: "item-1", ProductID: "product-1", Quantity: 2, UnitPrice: 50.0},
		},
		CreatedAt: now,
		UpdatedAt: &updatedAt,
	}

	if schema.ID != "order-1" {
		t.Errorf("OrderResponseSchema.ID = %v, want order-1", schema.ID)
	}
	if *schema.CustomerID != customerID {
		t.Errorf("OrderResponseSchema.CustomerID = %v, want %v", *schema.CustomerID, customerID)
	}
	if schema.Amount != 100.0 {
		t.Errorf("OrderResponseSchema.Amount = %v, want 100.0", schema.Amount)
	}
	if schema.Status != "Pending" {
		t.Errorf("OrderResponseSchema.Status = %v, want Pending", schema.Status)
	}
	if len(schema.Items) != 1 {
		t.Errorf("OrderResponseSchema.Items length = %v, want 1", len(schema.Items))
	}
	if schema.UpdatedAt == nil {
		t.Error("OrderResponseSchema.UpdatedAt should not be nil")
	}
}

func TestOrderResponseSchema_NilCustomerID(t *testing.T) {
	now := time.Now()

	schema := OrderResponseSchema{
		ID:         "order-1",
		CustomerID: nil,
		Amount:     100.0,
		Status:     "Pending",
		Items:      []OrderItemResponseSchema{},
		CreatedAt:  now,
		UpdatedAt:  nil,
	}

	if schema.CustomerID != nil {
		t.Errorf("OrderResponseSchema.CustomerID = %v, want nil", schema.CustomerID)
	}
	if schema.UpdatedAt != nil {
		t.Errorf("OrderResponseSchema.UpdatedAt = %v, want nil", schema.UpdatedAt)
	}
}

func TestOrderStatusResponseSchema_Structure(t *testing.T) {
	schema := OrderStatusResponseSchema{
		ID:   "status-1",
		Name: "Pending",
	}

	if schema.ID != "status-1" {
		t.Errorf("OrderStatusResponseSchema.ID = %v, want status-1", schema.ID)
	}
	if schema.Name != "Pending" {
		t.Errorf("OrderStatusResponseSchema.Name = %v, want Pending", schema.Name)
	}
}
