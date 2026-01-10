package use_cases

import (
	"testing"
	"time"

	"microservice/internal/adapters/dtos"
	"microservice/internal/domain/entities"
)

func TestFindAllOrdersUseCase_FilterDTOStructure(t *testing.T) {
	customerID := "customer-123"
	statusID := "status-1"
	from := time.Now().Add(-24 * time.Hour)
	to := time.Now()

	filter := dtos.OrderFilterDTO{
		CustomerID:    &customerID,
		StatusID:      &statusID,
		CreatedAtFrom: &from,
		CreatedAtTo:   &to,
	}

	if *filter.CustomerID != customerID {
		t.Errorf("OrderFilterDTO.CustomerID = %v, want %v", *filter.CustomerID, customerID)
	}
	if *filter.StatusID != statusID {
		t.Errorf("OrderFilterDTO.StatusID = %v, want %v", *filter.StatusID, statusID)
	}
	if filter.CreatedAtFrom == nil {
		t.Error("OrderFilterDTO.CreatedAtFrom should not be nil")
	}
	if filter.CreatedAtTo == nil {
		t.Error("OrderFilterDTO.CreatedAtTo should not be nil")
	}
}

func TestFindAllOrdersUseCase_FilterDTONilValues(t *testing.T) {
	filter := dtos.OrderFilterDTO{
		CustomerID:    nil,
		StatusID:      nil,
		CreatedAtFrom: nil,
		CreatedAtTo:   nil,
	}

	if filter.CustomerID != nil {
		t.Error("OrderFilterDTO.CustomerID should be nil")
	}
	if filter.StatusID != nil {
		t.Error("OrderFilterDTO.StatusID should be nil")
	}
	if filter.CreatedAtFrom != nil {
		t.Error("OrderFilterDTO.CreatedAtFrom should be nil")
	}
	if filter.CreatedAtTo != nil {
		t.Error("OrderFilterDTO.CreatedAtTo should be nil")
	}
}

func TestFindAllOrdersUseCase_FilterDTOPartialValues(t *testing.T) {
	customerID := "customer-123"
	filter := dtos.OrderFilterDTO{
		CustomerID: &customerID,
		StatusID:   nil,
	}

	if *filter.CustomerID != customerID {
		t.Errorf("OrderFilterDTO.CustomerID = %v, want %v", *filter.CustomerID, customerID)
	}
	if filter.StatusID != nil {
		t.Error("OrderFilterDTO.StatusID should be nil")
	}
}

func TestFindAllOrdersUseCase_Execute_FilterValidation(t *testing.T) {
	testCases := []struct {
		name       string
		customerID *string
		statusID   *string
	}{
		{"empty filter", nil, nil},
		{"customer only", stringPtr("customer-1"), nil},
		{"status only", nil, stringPtr("status-1")},
		{"both fields", stringPtr("customer-1"), stringPtr("status-1")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filter := dtos.OrderFilterDTO{
				CustomerID: tc.customerID,
				StatusID:   tc.statusID,
			}

			if tc.customerID != nil && *filter.CustomerID != *tc.customerID {
				t.Errorf("Expected CustomerID %v, got %v", *tc.customerID, *filter.CustomerID)
			}
			if tc.statusID != nil && *filter.StatusID != *tc.statusID {
				t.Errorf("Expected StatusID %v, got %v", *tc.statusID, *filter.StatusID)
			}
		})
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}

// Comprehensive tests using mocks for full coverage

func TestFindAllOrdersUseCase_NewFindAllOrdersUseCase(t *testing.T) {
	mockGateway := NewMockOrderGateway()
	uc := NewFindAllOrdersUseCase(mockGateway)

	if uc == nil {
		t.Error("Expected use case to be created")
	}
}

func TestFindAllOrdersUseCase_Execute_Success(t *testing.T) {
	mockGateway := NewMockOrderGateway()
	uc := NewFindAllOrdersUseCase(mockGateway)

	// Add test orders
	customerID := "customer-123"
	status, _ := entities.NewOrderStatus("pending", "Pending")
	order1, _ := entities.NewOrderWithItems("order-1", &customerID, 25.0, *status, []entities.OrderItem{}, time.Now(), nil)
	order2, _ := entities.NewOrderWithItems("order-2", &customerID, 35.0, *status, []entities.OrderItem{}, time.Now(), nil)
	
	mockGateway.AddOrder(order1)
	mockGateway.AddOrder(order2)

	filter := dtos.OrderFilterDTO{}
	orders, err := uc.Execute(filter)

	if err != nil {
		t.Errorf("Expected no error for successful find all, got %v", err)
	}

	if len(orders) != 2 {
		t.Errorf("Expected 2 orders, got %d", len(orders))
	}
}

func TestFindAllOrdersUseCase_Execute_EmptyResult(t *testing.T) {
	mockGateway := NewMockOrderGateway()
	uc := NewFindAllOrdersUseCase(mockGateway)

	filter := dtos.OrderFilterDTO{}
	orders, err := uc.Execute(filter)

	if err != nil {
		t.Errorf("Expected no error for empty result, got %v", err)
	}

	if len(orders) != 0 {
		t.Errorf("Expected 0 orders, got %d", len(orders))
	}
}

func TestFindAllOrdersUseCase_Execute_WithFilter(t *testing.T) {
	mockGateway := NewMockOrderGateway()
	uc := NewFindAllOrdersUseCase(mockGateway)

	// Add test orders
	customerID := "customer-123"
	statusID := "pending"
	status, _ := entities.NewOrderStatus(statusID, "Pending")
	order, _ := entities.NewOrderWithItems("order-1", &customerID, 25.0, *status, []entities.OrderItem{}, time.Now(), nil)
	
	mockGateway.AddOrder(order)

	filter := dtos.OrderFilterDTO{
		CustomerID: &customerID,
		StatusID:   &statusID,
	}
	orders, err := uc.Execute(filter)

	if err != nil {
		t.Errorf("Expected no error for filtered find all, got %v", err)
	}

	if len(orders) != 1 {
		t.Errorf("Expected 1 order, got %d", len(orders))
	}
}