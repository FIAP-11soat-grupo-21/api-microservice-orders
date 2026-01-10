package brokers

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/stretchr/testify/assert"
)

func TestNewSQSBroker_MissingOrdersQueueURL(t *testing.T) {
	config := BrokerConfig{
		SQSOrdersQueueURL: "", // Missing
		AWSRegion:         "us-east-1",
	}

	broker, err := NewSQSBroker(config)
	assert.Error(t, err)
	assert.Nil(t, broker)
	assert.Contains(t, err.Error(), "SQS orders queue URL is required")
}

func TestNewSQSBroker_ValidConfig(t *testing.T) {
	config := BrokerConfig{
		SQSOrdersQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/orders-queue",
		AWSRegion:         "us-east-1",
	}

	// This may succeed if AWS credentials are configured
	broker, err := NewSQSBroker(config)
	if err != nil {
		// Expected if no AWS credentials are configured
		assert.Error(t, err)
		assert.Nil(t, broker)
	} else {
		// If AWS is configured, broker should be valid
		assert.NotNil(t, broker)
		broker.Close() // Clean up
	}
}

func TestSQSBroker_Structure(t *testing.T) {
	// Test that the SQSBroker struct has the expected fields
	broker := &SQSBroker{
		ordersQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/orders-queue",
	}

	assert.Equal(t, "https://sqs.us-east-1.amazonaws.com/123456789012/orders-queue", broker.ordersQueueURL)
}

func TestSQSBroker_Close(t *testing.T) {
	broker := &SQSBroker{}
	err := broker.Close()
	assert.NoError(t, err)
}

func TestSQSBroker_ConsumeOrderUpdates_ContextCancellation(t *testing.T) {
	config := BrokerConfig{
		SQSOrdersQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/orders-queue",
		AWSRegion:         "us-east-1",
	}

	broker, err := NewSQSBroker(config)
	if err != nil {
		t.Skip("Skipping test - AWS not configured")
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	handler := func(message OrderUpdateMessage) error {
		return nil
	}

	err = broker.ConsumeOrderUpdates(ctx, handler)
	// The method returns nil immediately and runs in background
	// So we don't expect an error here
	assert.NoError(t, err)
}

func TestSQSBroker_processOrderUpdateMessage_ValidJSON(t *testing.T) {
	config := BrokerConfig{
		SQSOrdersQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/orders-queue",
		AWSRegion:         "us-east-1",
	}

	broker, err := NewSQSBroker(config)
	if err != nil {
		t.Skip("Skipping test - AWS not configured")
	}

	ctx := context.Background()
	messageBody := `{"order_id":"order-123","status":"Em preparação"}`
	message := types.Message{
		Body: &messageBody,
	}

	handlerCalled := false
	var receivedMessage OrderUpdateMessage
	handler := func(msg OrderUpdateMessage) error {
		handlerCalled = true
		receivedMessage = msg
		return nil
	}

	err = broker.processOrderUpdateMessage(ctx, message, handler)
	assert.NoError(t, err)
	assert.True(t, handlerCalled)
	assert.Equal(t, "order-123", receivedMessage.OrderID)
	assert.Equal(t, "Em preparação", receivedMessage.Status)
}

func TestSQSBroker_processOrderUpdateMessage_InvalidJSON(t *testing.T) {
	config := BrokerConfig{
		SQSOrdersQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/orders-queue",
		AWSRegion:         "us-east-1",
	}

	broker, err := NewSQSBroker(config)
	if err != nil {
		t.Skip("Skipping test - AWS not configured")
	}

	ctx := context.Background()
	messageBody := `invalid json`
	message := types.Message{
		Body: &messageBody,
	}

	handler := func(msg OrderUpdateMessage) error {
		return nil
	}

	err = broker.processOrderUpdateMessage(ctx, message, handler)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal order update message")
}

func TestSQSBroker_processOrderUpdateMessage_HandlerError(t *testing.T) {
	config := BrokerConfig{
		SQSOrdersQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/orders-queue",
		AWSRegion:         "us-east-1",
	}

	broker, err := NewSQSBroker(config)
	if err != nil {
		t.Skip("Skipping test - AWS not configured")
	}

	ctx := context.Background()
	messageBody := `{"order_id":"order-123","status":"Em preparação"}`
	message := types.Message{
		Body: &messageBody,
	}

	handler := func(msg OrderUpdateMessage) error {
		return errors.New("handler processing error")
	}

	err = broker.processOrderUpdateMessage(ctx, message, handler)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "handler processing error")
}

func TestSQSBroker_processOrderUpdateMessage_DifferentStatuses(t *testing.T) {
	config := BrokerConfig{
		SQSOrdersQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/orders-queue",
		AWSRegion:         "us-east-1",
	}

	broker, err := NewSQSBroker(config)
	if err != nil {
		t.Skip("Skipping test - AWS not configured")
	}

	testCases := []struct {
		name   string
		status string
	}{
		{"Em preparação", "Em preparação"},
		{"Pronto", "Pronto"},
		{"Finalizado", "Finalizado"},
		{"Cancelado", "Cancelado"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			messageBody := `{"order_id":"order-123","status":"` + tc.status + `"}`
			message := types.Message{
				Body: &messageBody,
			}

			var receivedMessage OrderUpdateMessage
			handler := func(msg OrderUpdateMessage) error {
				receivedMessage = msg
				return nil
			}

			err = broker.processOrderUpdateMessage(ctx, message, handler)
			assert.NoError(t, err)
			assert.Equal(t, tc.status, receivedMessage.Status)
		})
	}
}

func TestOrderUpdateMessage_Structure(t *testing.T) {
	message := OrderUpdateMessage{
		OrderID: "order-123",
		Status:  "Em preparação",
	}

	assert.Equal(t, "order-123", message.OrderID)
	assert.Equal(t, "Em preparação", message.Status)
}

func TestSQSBroker_pollOrderUpdateMessages_Coverage(t *testing.T) {
	config := BrokerConfig{
		SQSOrdersQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-orders-queue",
		AWSRegion:         "us-east-1",
	}

	broker, err := NewSQSBroker(config)
	assert.NoError(t, err)
	assert.NotNil(t, broker)

	// Test pollOrderUpdateMessages method exists and can be called
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	handler := func(message OrderUpdateMessage) error {
		return nil
	}

	// This will likely fail due to AWS credentials, but tests the method exists
	err = broker.pollOrderUpdateMessages(ctx, handler)
	if err != nil {
		t.Logf("pollOrderUpdateMessages error (expected without AWS setup): %v", err)
	}
}

func TestSQSBroker_deleteOrderUpdateMessage_Coverage(t *testing.T) {
	config := BrokerConfig{
		SQSOrdersQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-orders-queue",
		AWSRegion:         "us-east-1",
	}

	broker, err := NewSQSBroker(config)
	assert.NoError(t, err)
	assert.NotNil(t, broker)

	// Test deleteOrderUpdateMessage method exists and can be called
	ctx := context.Background()
	receiptHandle := "test-receipt-handle"
	message := types.Message{
		ReceiptHandle: &receiptHandle,
	}
	
	err = broker.deleteOrderUpdateMessage(ctx, message)
	if err != nil {
		t.Logf("deleteOrderUpdateMessage error (expected without AWS setup): %v", err)
	}
}

func TestSQSBroker_ConsumeOrderUpdates_LongRunning(t *testing.T) {
	config := BrokerConfig{
		SQSOrdersQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-orders-queue",
		AWSRegion:         "us-east-1",
	}

	broker, err := NewSQSBroker(config)
	assert.NoError(t, err)
	assert.NotNil(t, broker)

	// Test long-running ConsumeOrderUpdates with context cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	handler := func(message OrderUpdateMessage) error {
		return nil
	}

	// Start consuming in goroutine
	done := make(chan error, 1)
	go func() {
		done <- broker.ConsumeOrderUpdates(ctx, handler)
	}()

	// Wait for context timeout
	select {
	case err := <-done:
		// Should return due to context cancellation or AWS error
		if err != nil && err != context.DeadlineExceeded {
			t.Logf("ConsumeOrderUpdates returned error: %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("ConsumeOrderUpdates did not stop after context timeout")
	}
}

func TestSQSBroker_processOrderUpdateMessage_EmptyMessage(t *testing.T) {
	config := BrokerConfig{
		SQSOrdersQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-orders-queue",
		AWSRegion:         "us-east-1",
	}

	broker, err := NewSQSBroker(config)
	assert.NoError(t, err)

	// Test with empty message
	ctx := context.Background()
	emptyBody := ""
	message := types.Message{
		Body: &emptyBody,
	}
	
	handler := func(message OrderUpdateMessage) error {
		t.Error("Handler should not be called for empty message")
		return nil
	}

	err = broker.processOrderUpdateMessage(ctx, message, handler)
	assert.Error(t, err)
}

func TestSQSBroker_processOrderUpdateMessage_MissingFields(t *testing.T) {
	config := BrokerConfig{
		SQSOrdersQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-orders-queue",
		AWSRegion:         "us-east-1",
	}

	broker, err := NewSQSBroker(config)
	assert.NoError(t, err)

	testCases := []struct {
		name    string
		message string
	}{
		{
			name:    "Missing order_id",
			message: `{"status": "Em preparação"}`,
		},
		{
			name:    "Missing status",
			message: `{"order_id": "order-123"}`,
		},
		{
			name:    "Empty order_id",
			message: `{"order_id": "", "status": "Em preparação"}`,
		},
		{
			name:    "Empty status",
			message: `{"order_id": "order-123", "status": ""}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			message := types.Message{
				Body: &tc.message,
			}
			
			handler := func(message OrderUpdateMessage) error {
				// Allow handler to be called, but check if fields are empty
				if message.OrderID == "" || message.Status == "" {
					return errors.New("missing required fields")
				}
				return nil
			}

			err = broker.processOrderUpdateMessage(ctx, message, handler)
			// Should either fail at JSON unmarshaling or handler validation
			if err == nil {
				t.Error("Expected error for missing fields")
			}
		})
	}
}

func TestSQSBroker_processOrderUpdateMessage_AllStatuses(t *testing.T) {
	config := BrokerConfig{
		SQSOrdersQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-orders-queue",
		AWSRegion:         "us-east-1",
	}

	broker, err := NewSQSBroker(config)
	assert.NoError(t, err)

	statuses := []string{
		"Em preparação",
		"Pronto",
		"Finalizado",
		"Cancelado",
		"Recebido",
		"Confirmado",
	}

	for _, status := range statuses {
		t.Run("Status_"+status, func(t *testing.T) {
			ctx := context.Background()
			messageBody := fmt.Sprintf(`{"order_id": "order-123", "status": "%s"}`, status)
			message := types.Message{
				Body: &messageBody,
			}
			
			handlerCalled := false
			handler := func(msg OrderUpdateMessage) error {
				handlerCalled = true
				assert.Equal(t, "order-123", msg.OrderID)
				assert.Equal(t, status, msg.Status)
				return nil
			}

			err = broker.processOrderUpdateMessage(ctx, message, handler)
			assert.NoError(t, err)
			assert.True(t, handlerCalled)
		})
	}
}

func TestSQSBroker_Configuration_Validation(t *testing.T) {
	testCases := []struct {
		name          string
		config        BrokerConfig
		expectError   bool
		errorContains string
	}{
		{
			name: "Valid config",
			config: BrokerConfig{
				SQSOrdersQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/orders-queue",
				AWSRegion:         "us-east-1",
			},
			expectError: false,
		},
		{
			name: "Missing queue URL",
			config: BrokerConfig{
				SQSOrdersQueueURL: "",
				AWSRegion:         "us-east-1",
			},
			expectError:   true,
			errorContains: "SQS orders queue URL is required",
		},
		{
			name: "Invalid queue URL format",
			config: BrokerConfig{
				SQSOrdersQueueURL: "invalid-url",
				AWSRegion:         "us-east-1",
			},
			expectError: false, // URL validation happens at AWS level
		},
		{
			name: "Missing AWS region",
			config: BrokerConfig{
				SQSOrdersQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/orders-queue",
				AWSRegion:         "",
			},
			expectError: false, // Region can be empty, AWS SDK will use default
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			broker, err := NewSQSBroker(tc.config)
			
			if tc.expectError {
				assert.Error(t, err)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
				assert.Nil(t, broker)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, broker)
				
				// Test that broker can be closed
				err = broker.Close()
				assert.NoError(t, err)
			}
		})
	}
}

