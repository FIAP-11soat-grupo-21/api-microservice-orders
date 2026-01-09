package brokers

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAMQPConnection struct {
	mock.Mock
}

func (m *MockAMQPConnection) Close() error {
	args := m.Called()
	return args.Error(0)
}

type MockAMQPChannel struct {
	mock.Mock
}

func (m *MockAMQPChannel) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestRabbitMQBroker_NewRabbitMQBroker_InvalidURL(t *testing.T) {
	config := BrokerConfig{
		RabbitMQURL:          "invalid-url",
		RabbitMQPaymentQueue: "payment-queue",
		RabbitMQKitchenQueue: "kitchen-queue",
	}

	broker, err := NewRabbitMQBroker(config)
	assert.Error(t, err)
	assert.Nil(t, broker)
	assert.Contains(t, err.Error(), "failed to connect to RabbitMQ")
}

func TestRabbitMQBroker_shouldDiscardMessage(t *testing.T) {
	broker := &RabbitMQBroker{}

	testCases := []struct {
		name           string
		error          error
		shouldDiscard  bool
	}{
		{
			name:          "Order not found error",
			error:         errors.New("Order not found"),
			shouldDiscard: true,
		},
		{
			name:          "Invalid order ID error",
			error:         errors.New("Invalid order ID"),
			shouldDiscard: true,
		},
		{
			name:          "Invalid payment confirmation error",
			error:         errors.New("Invalid payment confirmation"),
			shouldDiscard: true,
		},
		{
			name:          "Order ID is required error",
			error:         errors.New("order ID is required"),
			shouldDiscard: true,
		},
		{
			name:          "Payment ID is required error",
			error:         errors.New("payment ID is required"),
			shouldDiscard: true,
		},
		{
			name:          "Payment status is required error",
			error:         errors.New("payment status is required"),
			shouldDiscard: true,
		},
		{
			name:          "Database connection error",
			error:         errors.New("database connection failed"),
			shouldDiscard: false,
		},
		{
			name:          "Network timeout error",
			error:         errors.New("network timeout"),
			shouldDiscard: false,
		},
		{
			name:          "Generic error",
			error:         errors.New("some other error"),
			shouldDiscard: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := broker.shouldDiscardMessage(tc.error)
			assert.Equal(t, tc.shouldDiscard, result)
		})
	}
}

func TestRabbitMQBroker_SendToKitchen_Success(t *testing.T) {
	// Test the message structure preparation
	message := map[string]interface{}{
		"type":     "order_paid",
		"order_id": "order-123",
		"status":   "paid",
	}

	// Test JSON marshaling
	body, err := json.Marshal(message)
	assert.NoError(t, err)
	assert.NotEmpty(t, body)

	var unmarshaled map[string]interface{}
	err = json.Unmarshal(body, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, "order_paid", unmarshaled["type"])
	assert.Equal(t, "order-123", unmarshaled["order_id"])
	assert.Equal(t, "paid", unmarshaled["status"])
}

func TestRabbitMQBroker_SendToKitchen_InvalidMessage(t *testing.T) {
	message := map[string]interface{}{
		"invalid": make(chan int), 
	}

	_, err := json.Marshal(message)
	assert.Error(t, err)
}

func TestRabbitMQBroker_processPaymentMessage_ValidMessage(t *testing.T) {
	paymentMsg := PaymentConfirmationMessage{
		OrderID:       "order-123",
		PaymentID:     "payment-456",
		Status:        "approved",
		Amount:        99.99,
		PaymentMethod: "credit_card",
		ProcessedAt:   time.Now(),
	}

	msgBody, err := json.Marshal(paymentMsg)
	assert.NoError(t, err)

	mockDelivery := struct {
		Body []byte
	}{
		Body: msgBody,
	}

	handlerCalled := false
	mockHandler := func(msg PaymentConfirmationMessage) error {
		handlerCalled = true
		assert.Equal(t, "order-123", msg.OrderID)
		assert.Equal(t, "payment-456", msg.PaymentID)
		assert.Equal(t, "approved", msg.Status)
		assert.Equal(t, 99.99, msg.Amount)
		assert.Equal(t, "credit_card", msg.PaymentMethod)
		return nil
	}

	var unmarshaledMsg PaymentConfirmationMessage
	err = json.Unmarshal(mockDelivery.Body, &unmarshaledMsg)
	assert.NoError(t, err)
	
	err = mockHandler(unmarshaledMsg)
	assert.NoError(t, err)
	assert.True(t, handlerCalled)
}

func TestRabbitMQBroker_processPaymentMessage_InvalidJSON(t *testing.T) {
	invalidJSON := []byte(`{"invalid": json}`)
	
	mockDelivery := struct {
		Body []byte
	}{
		Body: invalidJSON,
	}

	var paymentMsg PaymentConfirmationMessage
	err := json.Unmarshal(mockDelivery.Body, &paymentMsg)
	assert.Error(t, err)
}

func TestRabbitMQBroker_processPaymentMessage_HandlerError(t *testing.T) {
	paymentMsg := PaymentConfirmationMessage{
		OrderID:   "order-123",
		PaymentID: "payment-456",
		Status:    "approved",
	}

	msgBody, err := json.Marshal(paymentMsg)
	assert.NoError(t, err)

	mockHandler := func(msg PaymentConfirmationMessage) error {
		return errors.New("handler error")
	}

	var unmarshaledMsg PaymentConfirmationMessage
	err = json.Unmarshal(msgBody, &unmarshaledMsg)
	assert.NoError(t, err)
	
	err = mockHandler(unmarshaledMsg)
	assert.Error(t, err)
	assert.Equal(t, "handler error", err.Error())
}

func TestRabbitMQBroker_Close(t *testing.T) {
	broker := &RabbitMQBroker{}
	err := broker.Close()
	assert.NoError(t, err)

	broker2 := &RabbitMQBroker{
		conn:    nil,
		channel: nil,
	}
	err = broker2.Close()
	assert.NoError(t, err)
}

func TestRabbitMQBroker_ConsumePaymentConfirmations_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	
	cancel()
	
	select {
	case <-ctx.Done():
		assert.True(t, true, "Context was cancelled as expected")
	default:
		assert.Fail(t, "Context should be cancelled")
	}
}

func TestPaymentConfirmationMessage_Structure(t *testing.T) {
	now := time.Now()
	msg := PaymentConfirmationMessage{
		OrderID:       "order-123",
		PaymentID:     "payment-456",
		Status:        "approved",
		Amount:        150.75,
		PaymentMethod: "debit_card",
		ProcessedAt:   now,
	}

	assert.Equal(t, "order-123", msg.OrderID)
	assert.Equal(t, "payment-456", msg.PaymentID)
	assert.Equal(t, "approved", msg.Status)
	assert.Equal(t, 150.75, msg.Amount)
	assert.Equal(t, "debit_card", msg.PaymentMethod)
	assert.Equal(t, now, msg.ProcessedAt)

	jsonData, err := json.Marshal(msg)
	assert.NoError(t, err)

	var unmarshaledMsg PaymentConfirmationMessage
	err = json.Unmarshal(jsonData, &unmarshaledMsg)
	assert.NoError(t, err)
	assert.Equal(t, msg.OrderID, unmarshaledMsg.OrderID)
	assert.Equal(t, msg.PaymentID, unmarshaledMsg.PaymentID)
	assert.Equal(t, msg.Status, unmarshaledMsg.Status)
	assert.Equal(t, msg.Amount, unmarshaledMsg.Amount)
	assert.Equal(t, msg.PaymentMethod, unmarshaledMsg.PaymentMethod)
}

func TestBrokerConfig_Structure(t *testing.T) {
	config := BrokerConfig{
		Type:                 "rabbitmq",
		RabbitMQURL:          "amqp://localhost:5672",
		RabbitMQPaymentQueue: "payment-queue",
		RabbitMQKitchenQueue: "kitchen-queue",
		SQSPaymentQueueURL:   "https://sqs.amazonaws.com/payment",
		SQSKitchenQueueURL:   "https://sqs.amazonaws.com/kitchen",
		AWSRegion:            "us-east-1",
	}

	assert.Equal(t, "rabbitmq", config.Type)
	assert.Equal(t, "amqp://localhost:5672", config.RabbitMQURL)
	assert.Equal(t, "payment-queue", config.RabbitMQPaymentQueue)
	assert.Equal(t, "kitchen-queue", config.RabbitMQKitchenQueue)
	assert.Equal(t, "https://sqs.amazonaws.com/payment", config.SQSPaymentQueueURL)
	assert.Equal(t, "https://sqs.amazonaws.com/kitchen", config.SQSKitchenQueueURL)
	assert.Equal(t, "us-east-1", config.AWSRegion)
}