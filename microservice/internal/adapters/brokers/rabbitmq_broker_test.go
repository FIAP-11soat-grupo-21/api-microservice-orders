package brokers

import (
	"testing"

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
		Type:              "rabbitmq",
		RabbitMQURL:       "amqp://localhost:5672",
		RabbitMQOrdersQueue: "orders-queue",
		SQSOrdersQueueURL: "https://sqs.amazonaws.com/orders",
		AWSRegion:         "us-east-1",
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

type mockError struct {
	message string
}

func (e *mockError) Error() string {
	return e.message
}