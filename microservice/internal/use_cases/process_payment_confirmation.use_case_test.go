package use_cases

import (
	"testing"

	"microservice/internal/adapters/dtos"
	"microservice/internal/adapters/gateways"
	"microservice/internal/domain/entities"
	"microservice/internal/domain/value_objects"
)

func TestNewProcessPaymentConfirmationUseCase(t *testing.T) {
	orderGateway := &gateways.OrderGateway{}
	statusGateway := &gateways.OrderStatusGateway{}

	uc := NewProcessPaymentConfirmationUseCase(orderGateway, statusGateway)

	if uc == nil {
		t.Error("Expected use case to be created")
		return
	}

	if uc.orderGateway != orderGateway {
		t.Error("Expected order gateway to be set")
	}

	if uc.orderStatusGateway != statusGateway {
		t.Error("Expected order status gateway to be set")
	}
}

func TestProcessPaymentConfirmationUseCase_ValidateInput(t *testing.T) {
	uc := &ProcessPaymentConfirmationUseCase{}

	testCases := []struct {
		name        string
		dto         PaymentConfirmationDTO
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid input",
			dto: PaymentConfirmationDTO{
				OrderID:   "order-1",
				PaymentID: "payment-1",
				Status:    "confirmed",
				Amount:    25.0,
			},
			expectError: false,
		},
		{
			name: "missing order ID",
			dto: PaymentConfirmationDTO{
				PaymentID: "payment-1",
				Status:    "confirmed",
				Amount:    25.0,
			},
			expectError: true,
			errorMsg:    "order ID is required",
		},
		{
			name: "missing payment ID",
			dto: PaymentConfirmationDTO{
				OrderID: "order-1",
				Status:  "confirmed",
				Amount:  25.0,
			},
			expectError: true,
			errorMsg:    "payment ID is required",
		},
		{
			name: "missing status",
			dto: PaymentConfirmationDTO{
				OrderID:   "order-1",
				PaymentID: "payment-1",
				Amount:    25.0,
			},
			expectError: true,
			errorMsg:    "payment status is required",
		},
		{
			name: "invalid amount",
			dto: PaymentConfirmationDTO{
				OrderID:   "order-1",
				PaymentID: "payment-1",
				Status:    "confirmed",
				Amount:    0,
			},
			expectError: true,
			errorMsg:    "payment amount must be positive",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := uc.validateInput(tc.dto)

			if tc.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if err.Error() != tc.errorMsg {
					t.Errorf("Expected error '%s', got '%s'", tc.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}

func TestProcessPaymentConfirmationUseCase_CanUpdateOrderStatus_Simple(t *testing.T) {
	uc := &ProcessPaymentConfirmationUseCase{}

	amount, _ := value_objects.NewAmount(25.0)
	status, _ := entities.NewOrderStatus("status-id", "pending")

	order := &entities.Order{
		ID:     "order-1",
		Amount: amount,
		Status: *status,
	}

	result := uc.canUpdateOrderStatus(order, "confirmed")
	if !result {
		t.Error("Expected pending order to be updatable to confirmed")
	}

	result = uc.canUpdateOrderStatus(order, "failed")
	if !result {
		t.Error("Expected pending order to be updatable to failed")
	}
}

func TestPaymentConfirmationDTO_Structure(t *testing.T) {
	dto := PaymentConfirmationDTO{
		OrderID:       "order-1",
		PaymentID:     "payment-1",
		Status:        "confirmed",
		Amount:        25.0,
		PaymentMethod: "credit_card",
	}

	if dto.OrderID != "order-1" {
		t.Errorf("Expected OrderID 'order-1', got '%s'", dto.OrderID)
	}

	if dto.PaymentID != "payment-1" {
		t.Errorf("Expected PaymentID 'payment-1', got '%s'", dto.PaymentID)
	}

	if dto.Status != "confirmed" {
		t.Errorf("Expected Status 'confirmed', got '%s'", dto.Status)
	}

	if dto.Amount != 25.0 {
		t.Errorf("Expected Amount 25.0, got %f", dto.Amount)
	}

	if dto.PaymentMethod != "credit_card" {
		t.Errorf("Expected PaymentMethod 'credit_card', got '%s'", dto.PaymentMethod)
	}
}

func TestPaymentConfirmationResult_Structure(t *testing.T) {
	amount, _ := value_objects.NewAmount(25.0)
	status, _ := entities.NewOrderStatus("status-1", "pending")
	
	order := entities.Order{
		ID:     "order-1",
		Amount: amount,
		Status: *status,
	}

	result := PaymentConfirmationResult{
		Order:               order,
		StatusChanged:       true,
		ShouldNotifyKitchen: true,
		Message:             "Payment processed",
	}

	if result.Order.ID != "order-1" {
		t.Errorf("Expected Order ID 'order-1', got '%s'", result.Order.ID)
	}

	if !result.StatusChanged {
		t.Error("Expected StatusChanged to be true")
	}

	if !result.ShouldNotifyKitchen {
		t.Error("Expected ShouldNotifyKitchen to be true")
	}

	if result.Message != "Payment processed" {
		t.Errorf("Expected Message 'Payment processed', got '%s'", result.Message)
	}
}
func TestProcessPaymentConfirmationUseCase_Execute_ValidInput_ValidationOnly(t *testing.T) {
	uc := &ProcessPaymentConfirmationUseCase{}

	dto := PaymentConfirmationDTO{
		OrderID:   "550e8400-e29b-41d4-a716-446655440000",
		PaymentID: "payment-1",
		Status:    "confirmed",
		Amount:    25.0,
	}

	err := uc.validateInput(dto)
	if err != nil {
		t.Errorf("Expected no validation error, got %v", err)
	}
}

func TestProcessPaymentConfirmationUseCase_canUpdateOrderStatus_AllCases(t *testing.T) {
	uc := &ProcessPaymentConfirmationUseCase{}

	testCases := []struct {
		currentStatus string
		newStatus     string
		expected      bool
	}{
		{"pending", "confirmed", true},
		{"pending", "failed", true},
		{"pending", "cancelled", true},
		{"created", "confirmed", true},
		{"created", "failed", true},
		{"created", "cancelled", true},
		{"paid", "confirmed", false},
		{"paid", "failed", false},
		{"failed", "confirmed", false},
		{"cancelled", "confirmed", false},
		{"unknown", "confirmed", false},
	}

	for _, tc := range testCases {
		t.Run(tc.currentStatus+"_to_"+tc.newStatus, func(t *testing.T) {
			amount, _ := value_objects.NewAmount(25.0)
			status, _ := entities.NewOrderStatus("status-id", tc.currentStatus)

			order := &entities.Order{
				ID:     "order-1",
				Amount: amount,
				Status: *status,
			}

			result := uc.canUpdateOrderStatus(order, tc.newStatus)
			if result != tc.expected {
				t.Errorf("Expected %v for %s to %s, got %v", tc.expected, tc.currentStatus, tc.newStatus, result)
			}
		})
	}
}

func TestProcessPaymentConfirmationUseCase_findOrCreateStatus_MethodExists(t *testing.T) {
	// Test that the findOrCreateStatus method exists
	uc := &ProcessPaymentConfirmationUseCase{}

	// Just verify the method signature is correct
	_ = uc.findOrCreateStatus
}

func TestProcessPaymentConfirmationUseCase_processConfirmedPayment_MethodExists(t *testing.T) {
	uc := &ProcessPaymentConfirmationUseCase{}

	amount, _ := value_objects.NewAmount(99.99)
	status, _ := entities.NewOrderStatus("pending", "Pending")
	order := &entities.Order{
		ID:     "order-123",
		Amount: amount,
		Status: *status,
	}

	dto := PaymentConfirmationDTO{
		OrderID:   "order-123",
		PaymentID: "payment-456",
		Status:    "confirmed",
		Amount:    99.99,
	}

	_ = uc.processConfirmedPayment
	_ = order
	_ = dto
}

func TestProcessPaymentConfirmationUseCase_processFailedPayment_MethodExists(t *testing.T) {
	uc := &ProcessPaymentConfirmationUseCase{}

	amount, _ := value_objects.NewAmount(99.99)
	status, _ := entities.NewOrderStatus("pending", "Pending")
	order := &entities.Order{
		ID:     "order-123",
		Amount: amount,
		Status: *status,
	}

	dto := PaymentConfirmationDTO{
		OrderID:   "order-123",
		PaymentID: "payment-456",
		Status:    "failed",
		Amount:    99.99,
	}

	_ = uc.processFailedPayment
	_ = order
	_ = dto
}

func TestProcessPaymentConfirmationUseCase_updateOrder_MethodExists(t *testing.T) {
	uc := &ProcessPaymentConfirmationUseCase{}

	dto := dtos.UpdateOrderDTO{
		ID:       "550e8400-e29b-41d4-a716-446655440000",
		StatusID: "status-1",
	}

	_ = uc.updateOrder
	_ = dto
}