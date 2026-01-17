package use_cases

import (
	"context"
	"testing"

	"microservice/internal/adapters/brokers"
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

func TestCreateOrderUseCase_ValidatesItemIDs(t *testing.T) {
	// Test that valid IDs work
	_, err := entities.NewOrderItem("item-1", "product-1", "order-1", 1, 10.0)
	if err != nil {
		t.Errorf("Expected no error for valid IDs, got %v", err)
	}
}

func TestCreateOrderUseCase_OrderItemCreation(t *testing.T) {
	item, err := entities.NewOrderItem("item-1", "product-1", "order-1", 2, 15.0)
	if err != nil {
		t.Errorf("NewOrderItem() unexpected error: %v", err)
	}

	if item.ID != "item-1" {
		t.Errorf("OrderItem.ID = %v, want item-1", item.ID)
	}
	if item.ProductID.Value() != "product-1" {
		t.Errorf("OrderItem.ProductID = %v, want product-1", item.ProductID.Value())
	}
	if item.OrderID != "order-1" {
		t.Errorf("OrderItem.OrderID = %v, want order-1", item.OrderID)
	}
	if item.Quantity.Value() != 2 {
		t.Errorf("OrderItem.Quantity = %v, want 2", item.Quantity.Value())
	}
	if item.UnitPrice.Value() != 15.0 {
		t.Errorf("OrderItem.UnitPrice = %v, want 15.0", item.UnitPrice.Value())
	}
}

func TestCreateOrderUseCase_OrderCreation(t *testing.T) {
	customerID := "customer-123"
	order, err := entities.NewOrder("order-1", &customerID)
	if err != nil {
		t.Errorf("NewOrder() unexpected error: %v", err)
	}

	if order.ID != "order-1" {
		t.Errorf("Order.ID = %v, want order-1", order.ID)
	}
	if *order.CustomerID != customerID {
		t.Errorf("Order.CustomerID = %v, want %v", *order.CustomerID, customerID)
	}
}

func TestCreateOrderUseCase_DTOStructure(t *testing.T) {
	customerID := "customer-123"
	dto := dtos.CreateOrderDTO{
		CustomerID: &customerID,
		Items: []dtos.CreateOrderItemDTO{
			{
				ProductID: "product-1",
				Quantity:  2,
				Price:     10.0,
			},
		},
	}

	if *dto.CustomerID != customerID {
		t.Errorf("CreateOrderDTO.CustomerID = %v, want %v", *dto.CustomerID, customerID)
	}
	if len(dto.Items) != 1 {
		t.Errorf("CreateOrderDTO.Items length = %v, want 1", len(dto.Items))
	}
}

func TestCreateOrderUseCase_ItemDTOStructure(t *testing.T) {
	dto := dtos.CreateOrderItemDTO{
		ProductID: "product-1",
		Quantity:  2,
		Price:     15.0,
	}

	if dto.ProductID != "product-1" {
		t.Errorf("CreateOrderItemDTO.ProductID = %v, want product-1", dto.ProductID)
	}
	if dto.Quantity != 2 {
		t.Errorf("CreateOrderItemDTO.Quantity = %v, want 2", dto.Quantity)
	}
	if dto.Price != 15.0 {
		t.Errorf("CreateOrderItemDTO.Price = %v, want 15.0", dto.Price)
	}
}

func TestCreateOrderUseCase_OrderAmountCalculation(t *testing.T) {
	customerID := "customer-123"
	order, _ := entities.NewOrder("order-1", &customerID)

	item1, _ := entities.NewOrderItem("item-1", "product-1", "order-1", 2, 10.0)
	item2, _ := entities.NewOrderItem("item-2", "product-2", "order-1", 1, 15.0)

	order.AddItem(*item1)
	order.AddItem(*item2)

	err := order.CalcTotalAmount()
	if err != nil {
		t.Errorf("CalcTotalAmount() unexpected error: %v", err)
	}

	expectedAmount := 2*10.0 + 1*15.0 // 35.0
	if order.Amount.Value() != expectedAmount {
		t.Errorf("Order.Amount = %v, want %v", order.Amount.Value(), expectedAmount)
	}
}

// Mock message broker for testing
type MockMessageBroker struct{}

func (m *MockMessageBroker) ConsumeOrderUpdates(ctx context.Context, handler brokers.OrderUpdateHandler) error {
	return nil
}

func (m *MockMessageBroker) ConsumeOrderError(ctx context.Context, handler brokers.OrderErrorHandler) error {
	return nil
}

func (m *MockMessageBroker) PublishOnTopic(ctx context.Context, topic string, message interface{}) error {
	return nil
}

func (m *MockMessageBroker) Close() error {
	return nil
}

// Comprehensive tests using mocks for full coverage

func TestCreateOrderUseCase_NewCreateOrderUseCase(t *testing.T) {
	mockOrderGateway := NewMockOrderGateway()
	mockStatusGateway := NewMockOrderStatusGateway()
	mockBroker := &MockMessageBroker{}

	uc := NewCreateOrderUseCase(mockOrderGateway, mockStatusGateway, mockBroker)

	if uc == nil {
		t.Error("Expected use case to be created")
	}
}

func TestCreateOrderUseCase_Execute_Success(t *testing.T) {
	mockOrderGateway := NewMockOrderGateway()
	mockStatusGateway := NewMockOrderStatusGateway()
	mockBroker := &MockMessageBroker{}

	// Add initial status
	initialStatus, _ := entities.NewOrderStatus(INITIAL_ORDER_STATUS_ID, "Pending")
	mockStatusGateway.AddStatus(initialStatus)

	uc := NewCreateOrderUseCase(mockOrderGateway, mockStatusGateway, mockBroker)

	customerID := "customer-123"
	items := []dtos.CreateOrderItemDTO{
		{
			ProductID: "product-1",
			Quantity:  2,
			Price:     10.0,
		},
	}

	order, err := uc.Execute(&customerID, items)
	if err != nil {
		t.Errorf("Expected no error for successful create, got %v", err)
	}

	if order.IsEmpty() {
		t.Error("Expected non-empty order")
	}

	if *order.CustomerID != customerID {
		t.Errorf("Expected customer ID %s, got %s", customerID, *order.CustomerID)
	}

	if len(order.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(order.Items))
	}

	if order.Amount.Value() != 20.0 {
		t.Errorf("Expected amount 20.0, got %f", order.Amount.Value())
	}
}

func TestCreateOrderUseCase_Execute_StatusNotFound(t *testing.T) {
	mockOrderGateway := NewMockOrderGateway()
	mockStatusGateway := NewMockOrderStatusGateway()
	mockBroker := &MockMessageBroker{}

	// Don't add the initial status to simulate not found
	mockStatusGateway.SetShouldFailFindByID(true)

	uc := NewCreateOrderUseCase(mockOrderGateway, mockStatusGateway, mockBroker)

	customerID := "customer-123"
	items := []dtos.CreateOrderItemDTO{
		{
			ProductID: "product-1",
			Quantity:  2,
			Price:     10.0,
		},
	}

	order, err := uc.Execute(&customerID, items)
	if err == nil {
		t.Error("Expected error when status not found")
	}

	if !order.IsEmpty() {
		t.Error("Expected empty order when status not found")
	}
}

func TestCreateOrderUseCase_Execute_CreateError(t *testing.T) {
	mockOrderGateway := NewMockOrderGateway()
	mockStatusGateway := NewMockOrderStatusGateway()
	mockBroker := &MockMessageBroker{}

	// Add initial status
	initialStatus, _ := entities.NewOrderStatus(INITIAL_ORDER_STATUS_ID, "Pending")
	mockStatusGateway.AddStatus(initialStatus)

	// Make create fail
	mockOrderGateway.SetShouldFailCreate(true)

	uc := NewCreateOrderUseCase(mockOrderGateway, mockStatusGateway, mockBroker)

	customerID := "customer-123"
	items := []dtos.CreateOrderItemDTO{
		{
			ProductID: "product-1",
			Quantity:  2,
			Price:     10.0,
		},
	}

	order, err := uc.Execute(&customerID, items)
	if err == nil {
		t.Error("Expected error when create fails")
	}

	if !order.IsEmpty() {
		t.Error("Expected empty order when create fails")
	}
}
