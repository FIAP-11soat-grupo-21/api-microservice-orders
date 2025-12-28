package use_cases

import (
	"testing"

	"microservice/internal/adapters/dtos"
	"microservice/internal/domain/entities"
)

func TestCreateOrderUseCase_ValidatesItemQuantity(t *testing.T) {
	_, err := entities.NewOrderItem("item-1", "product-1", "order-1", 0, 10.0)
	if err == nil {
		t.Error("Expected error for zero quantity, got nil")
	}

	_, err = entities.NewOrderItem("item-1", "product-1", "order-1", -1, 10.0)
	if err == nil {
		t.Error("Expected error for negative quantity, got nil")
	}
}

func TestCreateOrderUseCase_ValidatesItemPrice(t *testing.T) {
	_, err := entities.NewOrderItem("item-1", "product-1", "order-1", 1, 0.0)
	if err == nil {
		t.Error("Expected error for zero price, got nil")
	}

	_, err = entities.NewOrderItem("item-1", "product-1", "order-1", 1, -10.0)
	if err == nil {
		t.Error("Expected error for negative price, got nil")
	}
}

func TestCreateOrderUseCase_ValidatesProductID(t *testing.T) {
	_, err := entities.NewOrderItem("item-1", "", "order-1", 1, 10.0)
	if err == nil {
		t.Error("Expected error for empty product ID, got nil")
	}
}

func TestCreateOrderUseCase_CalculatesTotalAmount(t *testing.T) {
	customerID := "customer-123"
	order, _ := entities.NewOrder("550e8400-e29b-41d4-a716-446655440000", &customerID)

	item1, _ := entities.NewOrderItem("item-1", "product-1", order.ID, 2, 10.0)
	item2, _ := entities.NewOrderItem("item-2", "product-2", order.ID, 3, 15.0)

	order.AddItem(*item1)
	order.AddItem(*item2)

	err := order.CalcTotalAmount()
	if err != nil {
		t.Errorf("CalcTotalAmount() unexpected error: %v", err)
	}

	expectedTotal := (2 * 10.0) + (3 * 15.0)
	if order.Amount.Value() != expectedTotal {
		t.Errorf("CalcTotalAmount() = %v, want %v", order.Amount.Value(), expectedTotal)
	}
}

func TestCreateOrderUseCase_FailsWithEmptyItems(t *testing.T) {
	customerID := "customer-123"
	order, _ := entities.NewOrder("550e8400-e29b-41d4-a716-446655440000", &customerID)

	err := order.CalcTotalAmount()
	if err == nil {
		t.Error("CalcTotalAmount() with empty items expected error, got nil")
	}
}

func TestCreateOrderUseCase_AllowsNilCustomerID(t *testing.T) {
	order, err := entities.NewOrder("550e8400-e29b-41d4-a716-446655440000", nil)
	if err != nil {
		t.Errorf("NewOrder() with nil customerID unexpected error: %v", err)
	}
	if order.CustomerID != nil {
		t.Errorf("NewOrder() CustomerID = %v, want nil", order.CustomerID)
	}
}

func TestCreateOrderItemDTO_Structure(t *testing.T) {
	dto := dtos.CreateOrderItemDTO{
		ProductID: "product-123",
		Quantity:  2,
		Price:     25.50,
	}

	if dto.ProductID != "product-123" {
		t.Errorf("CreateOrderItemDTO.ProductID = %v, want product-123", dto.ProductID)
	}
	if dto.Quantity != 2 {
		t.Errorf("CreateOrderItemDTO.Quantity = %v, want 2", dto.Quantity)
	}
	if dto.Price != 25.50 {
		t.Errorf("CreateOrderItemDTO.Price = %v, want 25.50", dto.Price)
	}
}

func TestCreateOrderDTO_Structure(t *testing.T) {
	customerID := "customer-123"
	dto := dtos.CreateOrderDTO{
		CustomerID: &customerID,
		Items: []dtos.CreateOrderItemDTO{
			{ProductID: "product-1", Quantity: 1, Price: 10.0},
		},
	}

	if *dto.CustomerID != customerID {
		t.Errorf("CreateOrderDTO.CustomerID = %v, want %v", *dto.CustomerID, customerID)
	}
	if len(dto.Items) != 1 {
		t.Errorf("CreateOrderDTO.Items length = %v, want 1", len(dto.Items))
	}
}

func TestInitialOrderStatusID_IsValid(t *testing.T) {
	if INITIAL_ORDER_STATUS_ID == "" {
		t.Error("INITIAL_ORDER_STATUS_ID should not be empty")
	}
	expectedID := "56d3b3c3-1801-49cd-bae7-972c78082001"
	if INITIAL_ORDER_STATUS_ID != expectedID {
		t.Errorf("INITIAL_ORDER_STATUS_ID = %v, want %v", INITIAL_ORDER_STATUS_ID, expectedID)
	}
}

func TestNewCreateOrderUseCase(t *testing.T) {
	if INITIAL_ORDER_STATUS_ID == "" {
		t.Error("INITIAL_ORDER_STATUS_ID should be defined")
	}
}
