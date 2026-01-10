package brokers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestNewRabbitMQBroker_ValidConfig(t *testing.T) {
	config := BrokerConfig{
		RabbitMQURL:         "amqp://localhost:5672",
		RabbitMQOrdersQueue: "orders-queue",
	}

	// This may succeed if RabbitMQ is running locally
	broker, err := NewRabbitMQBroker(config)
	if err != nil {
		// Expected if no RabbitMQ server is running
		assert.Error(t, err)
		assert.Nil(t, broker)
	} else {
		// If RabbitMQ is running, broker should be valid
		assert.NotNil(t, broker)
		broker.Close() // Clean up
	}
}

func TestNewRabbitMQBroker_InvalidURL(t *testing.T) {
	config := BrokerConfig{
		RabbitMQURL:         "invalid-url",
		RabbitMQOrdersQueue: "orders-queue",
	}

	// This may succeed if RabbitMQ is running locally
	broker, err := NewRabbitMQBroker(config)
	assert.Error(t, err)
	assert.Nil(t, broker)
}

func TestNewRabbitMQBroker_EmptyURL(t *testing.T) {
	config := BrokerConfig{
		RabbitMQURL:         "",
		RabbitMQOrdersQueue: "orders-queue",
	}

	broker, err := NewRabbitMQBroker(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "RabbitMQ URL is required")
	assert.Nil(t, broker)
}

func TestNewRabbitMQBroker_EmptyOrdersQueue(t *testing.T) {
	config := BrokerConfig{
		RabbitMQURL:         "amqp://localhost:5672",
		RabbitMQOrdersQueue: "",
	}

	broker, err := NewRabbitMQBroker(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "RabbitMQ orders queue is required")
	assert.Nil(t, broker)
}

func TestBrokerConfig_Structure(t *testing.T) {
	config := BrokerConfig{
		Type:                "rabbitmq",
		RabbitMQURL:         "amqp://localhost:5672",
		RabbitMQOrdersQueue: "orders-queue",
		SQSOrdersQueueURL:   "https://sqs.amazonaws.com/orders",
		AWSRegion:           "us-east-1",
	}

	assert.Equal(t, "rabbitmq", config.Type)
	assert.Equal(t, "amqp://localhost:5672", config.RabbitMQURL)
	assert.Equal(t, "orders-queue", config.RabbitMQOrdersQueue)
}

func TestRabbitMQBroker_shouldDiscardMessage(t *testing.T) {
	broker := &RabbitMQBroker{}

	// Test recoverable errors (should not discard)
	recoverableErrors := []string{
		"network error",
		"temporary failure",
		"connection lost",
	}

	for _, errMsg := range recoverableErrors {
		err := &mockError{message: errMsg}
		assert.False(t, broker.shouldDiscardMessage(err))
	}

	// Test non-recoverable errors (should discard)
	nonRecoverableErrors := []string{
		"Order not found",
		"Invalid order ID",
		"Invalid payment confirmation",
		"order ID is required",
		"payment ID is required",
		"payment status is required",
	}

	for _, errMsg := range nonRecoverableErrors {
		err := &mockError{message: errMsg}
		assert.True(t, broker.shouldDiscardMessage(err))
	}
}

func TestRabbitMQBroker_Close_NilChannel(t *testing.T) {
	broker := &RabbitMQBroker{
		channel: nil,
	}

	err := broker.Close()
	assert.NoError(t, err) // Should handle nil channel gracefully
}

type mockError struct {
	message string
}

func (e *mockError) Error() string {
	return e.message
}
func TestRabbitMQBroker_Structure(t *testing.T) {
	broker := &RabbitMQBroker{
		ordersQueue: "orders.updates",
	}

	assert.Equal(t, "orders.updates", broker.ordersQueue)
}

func TestRabbitMQBroker_shouldDiscardMessage_RecoverableErrors(t *testing.T) {
	broker := &RabbitMQBroker{}

	recoverableErrors := []string{
		"connection refused",
		"timeout",
		"network unreachable",
		"temporary failure",
	}

	for _, errMsg := range recoverableErrors {
		t.Run(errMsg, func(t *testing.T) {
			err := &mockError{message: errMsg}
			assert.False(t, broker.shouldDiscardMessage(err))
		})
	}
}

func TestRabbitMQBroker_shouldDiscardMessage_NonRecoverableErrors(t *testing.T) {
	broker := &RabbitMQBroker{}

	nonRecoverableErrors := []string{
		"Order not found",
		"Invalid order ID",
		"Invalid payment confirmation",
		"order ID is required",
		"payment ID is required",
		"payment status is required",
	}

	for _, errMsg := range nonRecoverableErrors {
		t.Run(errMsg, func(t *testing.T) {
			err := &mockError{message: errMsg}
			assert.True(t, broker.shouldDiscardMessage(err))
		})
	}
}
func TestRabbitMQBroker_ConsumeOrderUpdates_Integration(t *testing.T) {
	config := BrokerConfig{
		RabbitMQURL:         "amqp://localhost:5672",
		RabbitMQOrdersQueue: "test-orders-queue",
	}

	broker, err := NewRabbitMQBroker(config)
	if err != nil {
		t.Skipf("Skipping RabbitMQ integration test - RabbitMQ not available: %v", err)
		return
	}
	defer broker.Close()

	// Test ConsumeOrderUpdates method
	ctx := context.Background()
	handler := func(message OrderUpdateMessage) error {
		return nil
	}

	// This should not block in test environment
	go func() {
		err := broker.ConsumeOrderUpdates(ctx, handler)
		if err != nil {
			t.Logf("ConsumeOrderUpdates error (expected in test): %v", err)
		}
	}()

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)
}

func TestRabbitMQBroker_processOrderUpdateMessage_ValidJSON(t *testing.T) {
	broker := &RabbitMQBroker{}

	// Create a mock amqp.Delivery
	validJSON := `{"order_id": "order-123", "status": "Em preparação"}`
	mockDelivery := amqp.Delivery{
		Body: []byte(validJSON),
	}

	handlerCalled := false
	handler := func(message OrderUpdateMessage) error {
		handlerCalled = true
		assert.Equal(t, "order-123", message.OrderID)
		assert.Equal(t, "Em preparação", message.Status)
		return nil
	}

	err := broker.processOrderUpdateMessage(mockDelivery, handler)
	assert.NoError(t, err)
	assert.True(t, handlerCalled)
}

func TestRabbitMQBroker_processOrderUpdateMessage_InvalidJSON(t *testing.T) {
	broker := &RabbitMQBroker{}

	// Create a mock amqp.Delivery with invalid JSON
	invalidJSON := `{"invalid": json}`
	mockDelivery := amqp.Delivery{
		Body: []byte(invalidJSON),
	}

	handler := func(message OrderUpdateMessage) error {
		t.Error("Handler should not be called for invalid JSON")
		return nil
	}

	err := broker.processOrderUpdateMessage(mockDelivery, handler)
	assert.Error(t, err)
}

func TestRabbitMQBroker_processOrderUpdateMessage_HandlerError(t *testing.T) {
	broker := &RabbitMQBroker{}

	// Create a mock amqp.Delivery
	validJSON := `{"order_id": "order-123", "status": "Em preparação"}`
	mockDelivery := amqp.Delivery{
		Body: []byte(validJSON),
	}

	expectedError := errors.New("handler error")
	handler := func(message OrderUpdateMessage) error {
		return expectedError
	}

	err := broker.processOrderUpdateMessage(mockDelivery, handler)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestRabbitMQBroker_declareQueues_Coverage(t *testing.T) {
	config := BrokerConfig{
		RabbitMQURL:         "amqp://localhost:5672",
		RabbitMQOrdersQueue: "test-orders-queue",
	}

	broker, err := NewRabbitMQBroker(config)
	if err != nil {
		t.Skipf("Skipping RabbitMQ test - RabbitMQ not available: %v", err)
		return
	}
	defer broker.Close()

	// The declareQueues method is called during NewRabbitMQBroker
	// This test ensures it was called successfully
	assert.NotNil(t, broker.channel)
	assert.Equal(t, "test-orders-queue", broker.ordersQueue)
}

func TestRabbitMQBroker_ConsumeOrderUpdates_ContextCancellation(t *testing.T) {
	config := BrokerConfig{
		RabbitMQURL:         "amqp://localhost:5672",
		RabbitMQOrdersQueue: "test-orders-queue",
	}

	broker, err := NewRabbitMQBroker(config)
	if err != nil {
		t.Skipf("Skipping RabbitMQ test - RabbitMQ not available: %v", err)
		return
	}
	defer broker.Close()

	// Test context cancellation
	ctx, cancel := context.WithCancel(context.Background())

	handler := func(message OrderUpdateMessage) error {
		return nil
	}

	// Start consuming in goroutine
	done := make(chan error, 1)
	go func() {
		done <- broker.ConsumeOrderUpdates(ctx, handler)
	}()

	// Cancel context after a short delay
	time.Sleep(50 * time.Millisecond)
	cancel()

	// Wait for consumption to stop
	select {
	case err := <-done:
		// Should return due to context cancellation
		if err != nil && err != context.Canceled {
			t.Logf("ConsumeOrderUpdates returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Error("ConsumeOrderUpdates did not stop after context cancellation")
	}
}

func TestRabbitMQBroker_Close_WithNilChannel(t *testing.T) {
	// Create a broker with nil channel to test Close method
	broker := &RabbitMQBroker{
		conn:        nil,
		channel:     nil,
		ordersQueue: "test-queue",
	}

	// Should not panic with nil channel
	err := broker.Close()
	assert.NoError(t, err)
}

func TestRabbitMQBroker_Close_WithValidChannel(t *testing.T) {
	config := BrokerConfig{
		RabbitMQURL:         "amqp://localhost:5672",
		RabbitMQOrdersQueue: "test-orders-queue",
	}

	broker, err := NewRabbitMQBroker(config)
	if err != nil {
		t.Skipf("Skipping RabbitMQ test - RabbitMQ not available: %v", err)
		return
	}

	// Close should work without error
	err = broker.Close()
	assert.NoError(t, err)
}
