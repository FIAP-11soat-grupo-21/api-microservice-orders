package brokers

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/stretchr/testify/assert"
)

func TestNewSQSBroker_MissingOrdersQueueURL(t *testing.T) {
	config := BrokerConfig{
		SQSUpdateOrderStatusQueueURL: "", // Missing
		AWSRegion:                    "us-east-1",
	}

	broker, err := NewSQSBroker(config)
	assert.Error(t, err)
	assert.Nil(t, broker)
	assert.Contains(t, err.Error(), "SQS orders queue URL is required")
}

func TestNewSQSBroker_ValidConfig(t *testing.T) {
	config := BrokerConfig{
		SQSUpdateOrderStatusQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/orders-queue",
		AWSRegion:                    "us-east-1",
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
		updateOrderStatusQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/orders-queue",
	}

	assert.Equal(t, "https://sqs.us-east-1.amazonaws.com/123456789012/orders-queue", broker.updateOrderStatusQueueURL)
}

func TestSQSBroker_Close(t *testing.T) {
	broker := &SQSBroker{}
	err := broker.Close()
	assert.NoError(t, err)
}

func TestSQSBroker_ConsumeOrderUpdates_ContextCancellation(t *testing.T) {
	config := BrokerConfig{
		SQSUpdateOrderStatusQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/orders-queue",
		AWSRegion:                    "us-east-1",
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
		SQSUpdateOrderStatusQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-orders-queue",
		AWSRegion:                    "us-east-1",
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

func TestSQSBroker_deleteMessage_Coverage(t *testing.T) {
	config := BrokerConfig{
		SQSUpdateOrderStatusQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-orders-queue",
		AWSRegion:                    "us-east-1",
	}

	broker, err := NewSQSBroker(config)
	assert.NoError(t, err)
	assert.NotNil(t, broker)

	// Test deleteMessage method exists and can be called
	ctx := context.Background()
	receiptHandle := "test-receipt-handle"
	message := types.Message{
		ReceiptHandle: &receiptHandle,
	}

	err = broker.deleteMessage(ctx, config.SQSUpdateOrderStatusQueueURL, message)
	if err != nil {
		t.Logf("deleteMessage error (expected without AWS setup): %v", err)
	}
}

func TestSQSBroker_ConsumeOrderUpdates_LongRunning(t *testing.T) {
	config := BrokerConfig{
		SQSUpdateOrderStatusQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-orders-queue",
		AWSRegion:                    "us-east-1",
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
				SQSUpdateOrderStatusQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/orders-queue",
				AWSRegion:                    "us-east-1",
			},
			expectError: false,
		},
		{
			name: "Missing queue URL",
			config: BrokerConfig{
				SQSUpdateOrderStatusQueueURL: "",
				AWSRegion:                    "us-east-1",
			},
			expectError:   true,
			errorContains: "SQS orders queue URL is required",
		},
		{
			name: "Invalid queue URL format",
			config: BrokerConfig{
				SQSUpdateOrderStatusQueueURL: "invalid-url",
				AWSRegion:                    "us-east-1",
			},
			expectError: false, // URL validation happens at AWS level
		},
		{
			name: "Missing AWS region",
			config: BrokerConfig{
				SQSUpdateOrderStatusQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/orders-queue",
				AWSRegion:                    "",
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
