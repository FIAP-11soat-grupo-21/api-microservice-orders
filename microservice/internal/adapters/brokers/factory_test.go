package brokers

import (
	"testing"
)

func TestNewFactory(t *testing.T) {
	factory := NewFactory()
	if factory == nil {
		t.Error("Expected factory to be created")
	}
}

func TestFactory_CreateBroker_SQS(t *testing.T) {
	factory := NewFactory()
	config := BrokerConfig{
		Type:               "sqs",
		SQSPaymentQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/payment-queue",
		SQSKitchenQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/kitchen-queue",
		AWSRegion:          "us-east-1",
	}

	broker, err := factory.CreateBroker(config)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if broker == nil {
		t.Error("Expected broker to be created")
	}

	// Check if it's the correct type
	if _, ok := broker.(*SQSBroker); !ok {
		t.Error("Expected SQSBroker type")
	}
}

func TestFactory_CreateBroker_RabbitMQ(t *testing.T) {
	factory := NewFactory()
	config := BrokerConfig{
		Type:                 "rabbitmq",
		RabbitMQURL:          "amqp://guest:guest@localhost:5672/",
		RabbitMQPaymentQueue: "payment-confirmations",
		RabbitMQKitchenQueue: "kitchen-orders",
	}

	broker, err := factory.CreateBroker(config)
	if err != nil {
		t.Logf("RabbitMQ connection failed as expected in test environment: %v", err)
		return
	}

	if broker == nil {
		t.Error("Expected broker to be created")
	}

	if _, ok := broker.(*RabbitMQBroker); !ok {
		t.Error("Expected RabbitMQBroker type")
	}
}

func TestFactory_CreateBroker_CaseInsensitive(t *testing.T) {
	factory := NewFactory()
	
	testCases := []string{"SQS", "sqs", "Sqs", "SqS"}
	
	for _, brokerType := range testCases {
		config := BrokerConfig{
			Type:               brokerType,
			SQSPaymentQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/payment-queue",
			SQSKitchenQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/kitchen-queue",
			AWSRegion:          "us-east-1",
		}

		broker, err := factory.CreateBroker(config)
		if err != nil {
			t.Errorf("Expected no error for type '%s', got %v", brokerType, err)
		}

		if broker == nil {
			t.Errorf("Expected broker to be created for type '%s'", brokerType)
		}
	}
}

func TestFactory_CreateBroker_UnsupportedType(t *testing.T) {
	factory := NewFactory()
	config := BrokerConfig{
		Type: "unsupported",
	}

	broker, err := factory.CreateBroker(config)
	if err == nil {
		t.Error("Expected error for unsupported broker type")
	}

	if broker != nil {
		t.Error("Expected broker to be nil for unsupported type")
	}

	expectedError := "unsupported broker type: unsupported"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}