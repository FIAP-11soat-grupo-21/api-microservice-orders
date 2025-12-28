package entities

import (
	"testing"
	"time"
)

func TestNewOrder_ValidOrder(t *testing.T) {
	customerID := "customer-123"
	order, err := NewOrder("550e8400-e29b-41d4-a716-446655440000", &customerID)

	if err != nil {
		t.Errorf("NewOrder() unexpected error: %v", err)
	}
	if order == nil {
		t.Fatal("NewOrder() returned nil order")
	}
	if order.ID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("NewOrder() ID = %v, want %v", order.ID, "550e8400-e29b-41d4-a716-446655440000")
	}
	if *order.CustomerID != customerID {
		t.Errorf("NewOrder() CustomerID = %v, want %v", *order.CustomerID, customerID)
	}
	if len(order.Items) != 0 {
		t.Errorf("NewOrder() Items length = %v, want 0", len(order.Items))
	}
}

func TestNewOrder_NilCustomerID(t *testing.T) {
	order, err := NewOrder("550e8400-e29b-41d4-a716-446655440000", nil)

	if err != nil {
		t.Errorf("NewOrder() unexpected error: %v", err)
	}
	if order.CustomerID != nil {
		t.Errorf("NewOrder() CustomerID = %v, want nil", order.CustomerID)
	}
}

func TestOrder_AddItem(t *testing.T) {
	customerID := "customer-123"
	order, _ := NewOrder("550e8400-e29b-41d4-a716-446655440000", &customerID)

	item, _ := NewOrderItem("item-1", "product-1", order.ID, 2, 10.0)
	order.AddItem(*item)

	if len(order.Items) != 1 {
		t.Errorf("AddItem() Items length = %v, want 1", len(order.Items))
	}
}

func TestOrder_CalcTotalAmount(t *testing.T) {
	customerID := "customer-123"
	order, _ := NewOrder("550e8400-e29b-41d4-a716-446655440000", &customerID)

	item1, _ := NewOrderItem("item-1", "product-1", order.ID, 2, 10.0)
	item2, _ := NewOrderItem("item-2", "product-2", order.ID, 3, 15.0)

	order.AddItem(*item1)
	order.AddItem(*item2)

	err := order.CalcTotalAmount()
	if err != nil {
		t.Errorf("CalcTotalAmount() unexpected error: %v", err)
	}

	expectedTotal := (2 * 10.0) + (3 * 15.0) // 20 + 45 = 65
	if order.Amount.Value() != expectedTotal {
		t.Errorf("CalcTotalAmount() Amount = %v, want %v", order.Amount.Value(), expectedTotal)
	}
}

func TestOrder_CalcTotalAmount_EmptyItems(t *testing.T) {
	customerID := "customer-123"
	order, _ := NewOrder("550e8400-e29b-41d4-a716-446655440000", &customerID)

	err := order.CalcTotalAmount()
	if err == nil {
		t.Error("CalcTotalAmount() with empty items expected error, got nil")
	}
}

func TestOrder_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		orderID  string
		expected bool
	}{
		{"empty order", "", true},
		{"valid order", "550e8400-e29b-41d4-a716-446655440000", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := &Order{ID: tt.orderID}
			if order.IsEmpty() != tt.expected {
				t.Errorf("IsEmpty() = %v, want %v", order.IsEmpty(), tt.expected)
			}
		})
	}
}

func TestValidateID_ValidUUID(t *testing.T) {
	validUUIDs := []string{
		"550e8400-e29b-41d4-a716-446655440000",
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		"f47ac10b-58cc-4372-a567-0e02b2c3d479",
	}

	for _, uuid := range validUUIDs {
		t.Run(uuid, func(t *testing.T) {
			err := ValidateID(uuid)
			if err != nil {
				t.Errorf("ValidateID(%v) unexpected error: %v", uuid, err)
			}
		})
	}
}

func TestValidateID_InvalidUUID(t *testing.T) {
	invalidUUIDs := []string{
		"",
		"invalid-uuid",
		"123",
		"550e8400-e29b-41d4-a716",
		"not-a-uuid-at-all",
	}

	for _, uuid := range invalidUUIDs {
		t.Run(uuid, func(t *testing.T) {
			err := ValidateID(uuid)
			if err == nil {
				t.Errorf("ValidateID(%v) expected error, got nil", uuid)
			}
		})
	}
}

func TestNewOrderWithItems_ValidOrder(t *testing.T) {
	customerID := "customer-123"
	status, _ := NewOrderStatus("status-1", "Pending")
	item, _ := NewOrderItem("item-1", "product-1", "order-1", 2, 10.0)
	now := time.Now()
	updatedAt := now.Add(time.Hour)

	order, err := NewOrderWithItems(
		"550e8400-e29b-41d4-a716-446655440000",
		&customerID,
		20.0,
		*status,
		[]OrderItem{*item},
		now,
		&updatedAt,
	)

	if err != nil {
		t.Errorf("NewOrderWithItems() unexpected error: %v", err)
	}
	if order == nil {
		t.Fatal("NewOrderWithItems() returned nil")
	}
	if order.ID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("NewOrderWithItems() ID = %v, want 550e8400-e29b-41d4-a716-446655440000", order.ID)
	}
	if *order.CustomerID != customerID {
		t.Errorf("NewOrderWithItems() CustomerID = %v, want %v", *order.CustomerID, customerID)
	}
	if order.Amount.Value() != 20.0 {
		t.Errorf("NewOrderWithItems() Amount = %v, want 20.0", order.Amount.Value())
	}
	if order.Status.ID != "status-1" {
		t.Errorf("NewOrderWithItems() Status.ID = %v, want status-1", order.Status.ID)
	}
	if len(order.Items) != 1 {
		t.Errorf("NewOrderWithItems() Items length = %v, want 1", len(order.Items))
	}
	if order.UpdatedAt == nil {
		t.Error("NewOrderWithItems() UpdatedAt should not be nil")
	}
}

func TestNewOrderWithItems_InvalidAmount(t *testing.T) {
	customerID := "customer-123"
	status, _ := NewOrderStatus("status-1", "Pending")
	now := time.Now()

	_, err := NewOrderWithItems(
		"550e8400-e29b-41d4-a716-446655440000",
		&customerID,
		-10.0, // Invalid amount
		*status,
		[]OrderItem{},
		now,
		nil,
	)

	if err == nil {
		t.Error("NewOrderWithItems() with invalid amount expected error, got nil")
	}
}

func TestNewOrderWithItems_NilCustomerID(t *testing.T) {
	status, _ := NewOrderStatus("status-1", "Pending")
	item, _ := NewOrderItem("item-1", "product-1", "order-1", 1, 10.0)
	now := time.Now()

	order, err := NewOrderWithItems(
		"550e8400-e29b-41d4-a716-446655440000",
		nil,
		10.0,
		*status,
		[]OrderItem{*item},
		now,
		nil,
	)

	if err != nil {
		t.Errorf("NewOrderWithItems() unexpected error: %v", err)
	}
	if order.CustomerID != nil {
		t.Errorf("NewOrderWithItems() CustomerID = %v, want nil", order.CustomerID)
	}
}

func TestNewOrderWithItems_MultipleItems(t *testing.T) {
	customerID := "customer-123"
	status, _ := NewOrderStatus("status-1", "Pending")
	item1, _ := NewOrderItem("item-1", "product-1", "order-1", 2, 10.0)
	item2, _ := NewOrderItem("item-2", "product-2", "order-1", 3, 15.0)
	now := time.Now()

	order, err := NewOrderWithItems(
		"550e8400-e29b-41d4-a716-446655440000",
		&customerID,
		65.0, // 2*10 + 3*15
		*status,
		[]OrderItem{*item1, *item2},
		now,
		nil,
	)

	if err != nil {
		t.Errorf("NewOrderWithItems() unexpected error: %v", err)
	}
	if len(order.Items) != 2 {
		t.Errorf("NewOrderWithItems() Items length = %v, want 2", len(order.Items))
	}
}

func TestNewOrderWithItems_ZeroAmount(t *testing.T) {
	customerID := "customer-123"
	status, _ := NewOrderStatus("status-1", "Pending")
	now := time.Now()

	_, err := NewOrderWithItems(
		"550e8400-e29b-41d4-a716-446655440000",
		&customerID,
		0.0, // Zero amount - invalid
		*status,
		[]OrderItem{},
		now,
		nil,
	)

	if err == nil {
		t.Error("NewOrderWithItems() with zero amount expected error, got nil")
	}
}
