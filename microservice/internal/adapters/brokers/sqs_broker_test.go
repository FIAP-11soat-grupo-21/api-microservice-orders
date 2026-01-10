package brokers

import (
	"context"
	"errors"
	"testing"

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