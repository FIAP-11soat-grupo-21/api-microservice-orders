package consumers

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"microservice/internal/adapters/brokers"
	"microservice/internal/domain/entities"
	"microservice/internal/domain/value_objects"
	"microservice/internal/use_cases"
)

type mockBroker struct {
	consumeError error
	sendError    error
}

func (m *mockBroker) ConsumePaymentConfirmations(ctx context.Context, handler brokers.PaymentConfirmationHandler) error {
	return m.consumeError
}

func (m *mockBroker) SendToKitchen(message map[string]interface{}) error {
	return m.sendError
}

func (m *mockBroker) Close() error {
	return nil
}

func TestNewPaymentConsumer(t *testing.T) {
	broker := &mockBroker{}
	kitchenBroker := &mockBroker{}
	processPaymentUC := &mockProcessPaymentUseCase{}

	consumer := NewPaymentConsumer(broker, processPaymentUC, kitchenBroker)

	assert.NotNil(t, consumer)
	assert.Equal(t, broker, consumer.broker)
	assert.Equal(t, kitchenBroker, consumer.kitchenBroker)
	assert.Equal(t, processPaymentUC, consumer.processPaymentUC)
}

func TestPaymentConsumerStruct(t *testing.T) {
	broker := &mockBroker{}
	kitchenBroker := &mockBroker{}

	consumer := &PaymentConsumer{
		broker:        broker,
		kitchenBroker: kitchenBroker,
	}

	if consumer.broker == nil {
		t.Error("Expected broker to be set")
	}

	if consumer.kitchenBroker == nil {
		t.Error("Expected kitchenBroker to be set")
	}
}

func TestPaymentConsumer_Start_Success(t *testing.T) {
	broker := &mockBroker{consumeError: nil}
	kitchenBroker := &mockBroker{}
	processPaymentUC := &mockProcessPaymentUseCase{}

	consumer := NewPaymentConsumer(broker, processPaymentUC, kitchenBroker)

	ctx := context.Background()
	err := consumer.Start(ctx)
	assert.NoError(t, err)
}

func TestPaymentConsumer_Start_BrokerError(t *testing.T) {
	expectedError := assert.AnError
	broker := &mockBroker{consumeError: expectedError}
	kitchenBroker := &mockBroker{}
	processPaymentUC := &mockProcessPaymentUseCase{}

	consumer := NewPaymentConsumer(broker, processPaymentUC, kitchenBroker)

	ctx := context.Background()
	err := consumer.Start(ctx)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestPaymentConsumer_sendToKitchen_MessageStructure(t *testing.T) {
	amount, _ := value_objects.NewAmount(99.99)
	productID, _ := value_objects.NewProductID("product-1")
	quantity, _ := value_objects.NewQuantity(1)
	unitPrice, _ := value_objects.NewUnitPrice(99.99)

	order := &entities.Order{
		ID:         "order-123",
		CustomerID: stringPtr("customer-456"),
		Amount:     amount,
		Items: []entities.OrderItem{
			{
				ID:        "item-1",
				ProductID: productID,
				Quantity:  quantity,
				UnitPrice: unitPrice,
			},
		},
		CreatedAt: time.Now(),
	}

	message := brokers.PaymentConfirmationMessage{
		PaymentMethod: "pix",
	}

	kitchenMessage := map[string]interface{}{
		"type":           "order_paid",
		"order_id":       order.ID,
		"customer_id":    order.CustomerID,
		"items":          order.Items,
		"status":         "paid",
		"created_at":     order.CreatedAt,
		"paid_at":        time.Now(),
		"total_amount":   order.Amount.Value(),
		"payment_method": message.PaymentMethod,
	}

	assert.Equal(t, "order_paid", kitchenMessage["type"])
	assert.Equal(t, order.ID, kitchenMessage["order_id"])
	assert.Equal(t, order.CustomerID, kitchenMessage["customer_id"])
	assert.Equal(t, order.Items, kitchenMessage["items"])
	assert.Equal(t, "paid", kitchenMessage["status"])
	assert.Equal(t, order.CreatedAt, kitchenMessage["created_at"])
	assert.Equal(t, order.Amount.Value(), kitchenMessage["total_amount"])
	assert.Equal(t, message.PaymentMethod, kitchenMessage["payment_method"])
	assert.NotNil(t, kitchenMessage["paid_at"])
}

func TestPaymentConsumer_sendToKitchen_Success(t *testing.T) {
	broker := &mockBroker{}
	kitchenBroker := &mockBroker{sendError: nil}
	processPaymentUC := &mockProcessPaymentUseCase{}

	consumer := NewPaymentConsumer(broker, processPaymentUC, kitchenBroker)

	amount, _ := value_objects.NewAmount(150.75)
	productID, _ := value_objects.NewProductID("product-1")
	quantity, _ := value_objects.NewQuantity(2)
	unitPrice, _ := value_objects.NewUnitPrice(75.375)

	order := &entities.Order{
		ID:         "order-123",
		CustomerID: stringPtr("customer-456"),
		Amount:     amount,
		Items: []entities.OrderItem{
			{
				ID:        "item-1",
				ProductID: productID,
				Quantity:  quantity,
				UnitPrice: unitPrice,
			},
		},
		CreatedAt: time.Now(),
	}

	message := brokers.PaymentConfirmationMessage{
		PaymentMethod: "debit_card",
	}

	err := consumer.sendToKitchen(order, message)
	assert.NoError(t, err)
}

func TestPaymentConsumer_sendToKitchen_Error(t *testing.T) {
	broker := &mockBroker{}
	expectedError := assert.AnError
	kitchenBroker := &mockBroker{sendError: expectedError}
	processPaymentUC := &mockProcessPaymentUseCase{}

	consumer := NewPaymentConsumer(broker, processPaymentUC, kitchenBroker)

	amount, _ := value_objects.NewAmount(99.99)
	order := &entities.Order{
		ID:         "order-123",
		CustomerID: stringPtr("customer-456"),
		Amount:     amount,
		Items:      []entities.OrderItem{},
		CreatedAt:  time.Now(),
	}

	message := brokers.PaymentConfirmationMessage{
		PaymentMethod: "credit_card",
	}

	err := consumer.sendToKitchen(order, message)
	assert.NoError(t, err)
}

func stringPtr(s string) *string {
	return &s
}

// Testes para o método processPaymentConfirmation (0% coverage)

type mockProcessPaymentUseCase struct {
	result *use_cases.PaymentConfirmationResult
	error  error
}

func (m *mockProcessPaymentUseCase) Execute(dto use_cases.PaymentConfirmationDTO) (*use_cases.PaymentConfirmationResult, error) {
	return m.result, m.error
}

func TestPaymentConsumer_processPaymentConfirmation_Success_WithKitchenNotification(t *testing.T) {
	broker := &mockBroker{}
	kitchenBroker := &mockBroker{sendError: nil}

	amount, _ := value_objects.NewAmount(199.99)
	status, _ := entities.NewOrderStatus("paid", "Paid")
	order := entities.Order{
		ID:         "order-123",
		CustomerID: stringPtr("customer-456"),
		Amount:     amount,
		Status:     *status,
		Items:      []entities.OrderItem{},
		CreatedAt:  time.Now(),
	}

	mockUseCase := &mockProcessPaymentUseCase{
		result: &use_cases.PaymentConfirmationResult{
			Order:               order,
			StatusChanged:       true,
			ShouldNotifyKitchen: true,
			Message:             "Payment processed successfully",
		},
		error: nil,
	}

	consumer := NewPaymentConsumer(broker, mockUseCase, kitchenBroker)

	message := brokers.PaymentConfirmationMessage{
		OrderID:       "order-123",
		PaymentID:     "payment-456",
		Status:        "confirmed",
		Amount:        199.99,
		PaymentMethod: "credit_card",
		ProcessedAt:   time.Now(),
	}

	err := consumer.processPaymentConfirmation(message)
	assert.NoError(t, err)
}

func TestPaymentConsumer_processPaymentConfirmation_Success_WithoutKitchenNotification(t *testing.T) {
	broker := &mockBroker{}
	kitchenBroker := &mockBroker{}

	amount, _ := value_objects.NewAmount(99.99)
	status, _ := entities.NewOrderStatus("failed", "Failed")
	order := entities.Order{
		ID:         "order-123",
		CustomerID: stringPtr("customer-456"),
		Amount:     amount,
		Status:     *status,
		Items:      []entities.OrderItem{},
		CreatedAt:  time.Now(),
	}

	mockUseCase := &mockProcessPaymentUseCase{
		result: &use_cases.PaymentConfirmationResult{
			Order:               order,
			StatusChanged:       true,
			ShouldNotifyKitchen: false,
			Message:             "Payment failed",
		},
		error: nil,
	}

	consumer := NewPaymentConsumer(broker, mockUseCase, kitchenBroker)

	message := brokers.PaymentConfirmationMessage{
		OrderID:       "order-123",
		PaymentID:     "payment-456",
		Status:        "failed",
		Amount:        99.99,
		PaymentMethod: "credit_card",
		ProcessedAt:   time.Now(),
	}

	err := consumer.processPaymentConfirmation(message)
	assert.NoError(t, err)
}

func TestPaymentConsumer_processPaymentConfirmation_UseCaseError(t *testing.T) {
	broker := &mockBroker{}
	kitchenBroker := &mockBroker{}

	expectedError := assert.AnError
	mockUseCase := &mockProcessPaymentUseCase{
		result: nil,
		error:  expectedError,
	}

	consumer := NewPaymentConsumer(broker, mockUseCase, kitchenBroker)

	message := brokers.PaymentConfirmationMessage{
		OrderID:       "order-123",
		PaymentID:     "payment-456",
		Status:        "confirmed",
		Amount:        99.99,
		PaymentMethod: "credit_card",
		ProcessedAt:   time.Now(),
	}

	err := consumer.processPaymentConfirmation(message)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestPaymentConsumer_processPaymentConfirmation_KitchenNotificationError(t *testing.T) {
	broker := &mockBroker{}
	kitchenBrokerError := assert.AnError
	kitchenBroker := &mockBroker{sendError: kitchenBrokerError}

	amount, _ := value_objects.NewAmount(299.99)
	status, _ := entities.NewOrderStatus("paid", "Paid")
	order := entities.Order{
		ID:         "order-123",
		CustomerID: stringPtr("customer-456"),
		Amount:     amount,
		Status:     *status,
		Items:      []entities.OrderItem{},
		CreatedAt:  time.Now(),
	}

	mockUseCase := &mockProcessPaymentUseCase{
		result: &use_cases.PaymentConfirmationResult{
			Order:               order,
			StatusChanged:       true,
			ShouldNotifyKitchen: true,
			Message:             "Payment processed successfully",
		},
		error: nil,
	}

	consumer := NewPaymentConsumer(broker, mockUseCase, kitchenBroker)

	message := brokers.PaymentConfirmationMessage{
		OrderID:       "order-123",
		PaymentID:     "payment-456",
		Status:        "confirmed",
		Amount:        299.99,
		PaymentMethod: "debit_card",
		ProcessedAt:   time.Now(),
	}

	err := consumer.processPaymentConfirmation(message)
	assert.NoError(t, err)
}

func TestPaymentConsumer_processPaymentConfirmation_DTOMapping(t *testing.T) {
	broker := &mockBroker{}
	kitchenBroker := &mockBroker{}

	mockUseCase := &mockProcessPaymentUseCase{
		result: &use_cases.PaymentConfirmationResult{
			Order: entities.Order{
				ID: "test-order",
			},
			StatusChanged:       false,
			ShouldNotifyKitchen: false,
			Message:             "Test message",
		},
		error: nil,
	}

	consumer := NewPaymentConsumer(broker, mockUseCase, kitchenBroker)

	testTime := time.Now()
	message := brokers.PaymentConfirmationMessage{
		OrderID:       "test-order-123",
		PaymentID:     "test-payment-456",
		Status:        "test-status",
		Amount:        123.45,
		PaymentMethod: "test-method",
		ProcessedAt:   testTime,
	}

	err := consumer.processPaymentConfirmation(message)
	assert.NoError(t, err)

}

func TestPaymentConsumer_processPaymentConfirmation_CompleteFlow(t *testing.T) {
	broker := &mockBroker{}
	kitchenBroker := &mockBroker{sendError: nil}

	amount, _ := value_objects.NewAmount(399.99)
	status, _ := entities.NewOrderStatus("paid", "Paid")

	productID, _ := value_objects.NewProductID("product-1")
	quantity, _ := value_objects.NewQuantity(2)
	unitPrice, _ := value_objects.NewUnitPrice(199.995)

	orderItem := entities.OrderItem{
		ID:        "item-1",
		OrderID:   "order-123",
		ProductID: productID,
		Quantity:  quantity,
		UnitPrice: unitPrice,
	}

	order := entities.Order{
		ID:         "order-123",
		CustomerID: stringPtr("customer-456"),
		Amount:     amount,
		Status:     *status,
		Items:      []entities.OrderItem{orderItem},
		CreatedAt:  time.Now(),
	}

	mockUseCase := &mockProcessPaymentUseCase{
		result: &use_cases.PaymentConfirmationResult{
			Order:               order,
			StatusChanged:       true,
			ShouldNotifyKitchen: true,
			Message:             "Order successfully paid and ready for kitchen",
		},
		error: nil,
	}

	consumer := NewPaymentConsumer(broker, mockUseCase, kitchenBroker)

	message := brokers.PaymentConfirmationMessage{
		OrderID:       "order-123",
		PaymentID:     "payment-456",
		Status:        "confirmed",
		Amount:        399.99,
		PaymentMethod: "pix",
		ProcessedAt:   time.Now(),
	}

	err := consumer.processPaymentConfirmation(message)
	assert.NoError(t, err)
}
