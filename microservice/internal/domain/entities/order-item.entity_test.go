package entities

import (
	"testing"
)

func TestNewOrderItem_ValidItem(t *testing.T) {
	item, err := NewOrderItem("item-1", "product-123", "order-123", 2, 25.50)

	if err != nil {
		t.Errorf("NewOrderItem() unexpected error: %v", err)
	}
	if item == nil {
		t.Fatal("NewOrderItem() returned nil")
	}
	if item.ID != "item-1" {
		t.Errorf("NewOrderItem() ID = %v, want item-1", item.ID)
	}
	if item.OrderID != "order-123" {
		t.Errorf("NewOrderItem() OrderID = %v, want order-123", item.OrderID)
	}
	if item.ProductID.Value() != "product-123" {
		t.Errorf("NewOrderItem() ProductID = %v, want product-123", item.ProductID.Value())
	}
	if item.Quantity.Value() != 2 {
		t.Errorf("NewOrderItem() Quantity = %v, want 2", item.Quantity.Value())
	}
	if item.UnitPrice.Value() != 25.50 {
		t.Errorf("NewOrderItem() UnitPrice = %v, want 25.50", item.UnitPrice.Value())
	}
}

func TestNewOrderItem_InvalidProductID(t *testing.T) {
	_, err := NewOrderItem("item-1", "", "order-123", 2, 25.50)

	if err == nil {
		t.Error("NewOrderItem() with empty productID expected error, got nil")
	}
}

func TestNewOrderItem_InvalidQuantity(t *testing.T) {
	tests := []struct {
		name     string
		quantity int
	}{
		{"zero quantity", 0},
		{"negative quantity", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewOrderItem("item-1", "product-123", "order-123", tt.quantity, 25.50)
			if err == nil {
				t.Errorf("NewOrderItem() with quantity %d expected error, got nil", tt.quantity)
			}
		})
	}
}

func TestNewOrderItem_InvalidUnitPrice(t *testing.T) {
	tests := []struct {
		name      string
		unitPrice float64
	}{
		{"zero price", 0.0},
		{"negative price", -10.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewOrderItem("item-1", "product-123", "order-123", 2, tt.unitPrice)
			if err == nil {
				t.Errorf("NewOrderItem() with unitPrice %v expected error, got nil", tt.unitPrice)
			}
		})
	}
}

func TestOrderItem_GetTotal(t *testing.T) {
	tests := []struct {
		name      string
		quantity  int
		unitPrice float64
		expected  float64
	}{
		{"single item", 1, 10.0, 10.0},
		{"multiple items", 3, 15.50, 46.50},
		{"large quantity", 100, 5.0, 500.0},
		{"decimal price", 2, 19.99, 39.98},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := NewOrderItem("item-1", "product-123", "order-123", tt.quantity, tt.unitPrice)
			if err != nil {
				t.Fatalf("NewOrderItem() unexpected error: %v", err)
			}

			total := item.GetTotal()
			if total != tt.expected {
				t.Errorf("GetTotal() = %v, want %v", total, tt.expected)
			}
		})
	}
}
