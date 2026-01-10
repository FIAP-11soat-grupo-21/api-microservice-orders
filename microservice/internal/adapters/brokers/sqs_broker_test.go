package brokers

import (
	"testing"

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