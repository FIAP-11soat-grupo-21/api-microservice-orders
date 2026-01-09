package use_cases

import (
	"testing"
	"time"

	"microservice/internal/adapters/dtos"
	"microservice/internal/adapters/gateways"
	"microservice/internal/domain/entities"
	"microservice/internal/domain/exceptions"
)

func TestUpdateOrderUseCase_ValidatesUUID(t *testing.T) {
	invalidIDs := []string{
		"",
		"invalid-uuid",
		"123",
	}

	for _, id := range invalidIDs {
		t.Run(id, func(t *testing.T) {
			err := entities.ValidateID(id)
			if err == nil {
				t.Errorf("ValidateID(%v) expected error, got nil", id)
			}
		})
	}
}

func TestUpdateOrderUseCase_DTOStructure(t *testing.T) {
	dto := dtos.UpdateOrderDTO{
		ID:       "550e8400-e29b-41d4-a716-446655440000",
		StatusID: "status-2",
	}

	if dto.ID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("UpdateOrderDTO.ID = %v, want 550e8400-e29b-41d4-a716-446655440000", dto.ID)
	}
	if dto.StatusID != "status-2" {
		t.Errorf("UpdateOrderDTO.StatusID = %v, want status-2", dto.StatusID)
	}
}

func TestUpdateOrderUseCase_StatusNotFoundExceptionMessage(t *testing.T) {
	err := &exceptions.OrderStatusNotFoundException{}
	if err.Error() != "Order Status not found" {
		t.Errorf("OrderStatusNotFoundException.Error() = %v, want 'Order Status not found'", err.Error())
	}
}

func TestUpdateOrderUseCase_OrderCanUpdateStatus(t *testing.T) {
	customerID := "customer-123"
	status1, _ := entities.NewOrderStatus("status-1", "Pending")
	status2, _ := entities.NewOrderStatus("status-2", "Confirmed")
	item, _ := entities.NewOrderItem("item-1", "product-1", "order-1", 2, 10.0)
	now := time.Now()

	order, _ := entities.NewOrderWithItems(
		"550e8400-e29b-41d4-a716-446655440000",
		&customerID,
		20.0,
		*status1,
		[]entities.OrderItem{*item},
		now,
		nil,
	)

	order.Status = *status2
	updateTime := time.Now()
	order.UpdatedAt = &updateTime

	if order.Status.ID != "status-2" {
		t.Errorf("Order.Status.ID = %v, want status-2", order.Status.ID)
	}
	if order.UpdatedAt == nil {
		t.Error("Order.UpdatedAt should not be nil after update")
	}
}

func TestUpdateOrderUseCase_PreservesOrderData(t *testing.T) {
	customerID := "customer-123"
	status, _ := entities.NewOrderStatus("status-1", "Pending")
	item, _ := entities.NewOrderItem("item-1", "product-1", "order-1", 2, 10.0)
	createdAt := time.Now().Add(-24 * time.Hour)

	order, _ := entities.NewOrderWithItems(
		"550e8400-e29b-41d4-a716-446655440000",
		&customerID,
		20.0,
		*status,
		[]entities.OrderItem{*item},
		createdAt,
		nil,
	)

	if order.ID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("Order.ID = %v, want 550e8400-e29b-41d4-a716-446655440000", order.ID)
	}
	if *order.CustomerID != customerID {
		t.Errorf("Order.CustomerID = %v, want %v", *order.CustomerID, customerID)
	}
	if order.Amount.Value() != 20.0 {
		t.Errorf("Order.Amount = %v, want 20.0", order.Amount.Value())
	}
	if len(order.Items) != 1 {
		t.Errorf("Order.Items length = %v, want 1", len(order.Items))
	}
}

func TestNewUpdateOrderUseCase(t *testing.T) {
	_ = NewUpdateOrderUseCase
}
func TestUpdateOrderUseCase_NewUpdateOrderUseCase(t *testing.T) {
	orderGateway := gateways.OrderGateway{}
	statusGateway := gateways.OrderStatusGateway{}
	uc := NewUpdateOrderUseCase(orderGateway, statusGateway)

	if uc == nil {
		t.Error("Expected use case to be created")
	}
}

func TestUpdateOrderUseCase_Execute_InvalidID(t *testing.T) {
	orderGateway := gateways.OrderGateway{}
	statusGateway := gateways.OrderStatusGateway{}
	uc := NewUpdateOrderUseCase(orderGateway, statusGateway)

	dto := dtos.UpdateOrderDTO{
		ID:       "invalid-id",
		StatusID: "status-1",
	}

	_, err := uc.Execute(dto)
	if err == nil {
		t.Error("Expected error for invalid ID")
	}

	if _, ok := err.(*exceptions.InvalidOrderDataException); !ok {
		t.Errorf("Expected InvalidOrderDataException, got %T", err)
	}
}

func TestUpdateOrderUseCase_Execute_EmptyID(t *testing.T) {
	orderGateway := gateways.OrderGateway{}
	statusGateway := gateways.OrderStatusGateway{}
	uc := NewUpdateOrderUseCase(orderGateway, statusGateway)

	dto := dtos.UpdateOrderDTO{
		ID:       "",
		StatusID: "status-1",
	}

	_, err := uc.Execute(dto)
	if err == nil {
		t.Error("Expected error for empty ID")
	}
}

func TestUpdateOrderUseCase_Execute_ValidIDFormat(t *testing.T) {
	validID := "550e8400-e29b-41d4-a716-446655440000"
	err := entities.ValidateID(validID)
	if err != nil {
		t.Errorf("Expected no error for valid UUID, got %v", err)
	}
}

func TestUpdateOrderUseCase_Execute_ReturnsEmptyOrderOnError(t *testing.T) {
	orderGateway := gateways.OrderGateway{}
	statusGateway := gateways.OrderStatusGateway{}
	uc := NewUpdateOrderUseCase(orderGateway, statusGateway)

	dto := dtos.UpdateOrderDTO{
		ID:       "invalid-id",
		StatusID: "status-1",
	}

	order, err := uc.Execute(dto)
	if err == nil {
		t.Error("Expected error for invalid ID")
	}

	if !order.IsEmpty() {
		t.Error("Expected empty order on error")
	}
}

func TestUpdateOrderDTO_Structure(t *testing.T) {
	dto := dtos.UpdateOrderDTO{
		ID:       "order-123",
		StatusID: "status-456",
	}

	if dto.ID != "order-123" {
		t.Errorf("Expected ID 'order-123', got '%s'", dto.ID)
	}

	if dto.StatusID != "status-456" {
		t.Errorf("Expected StatusID 'status-456', got '%s'", dto.StatusID)
	}
}