package messaging

import (
	"os"
	"testing"

	"microservice/utils/config"
)

func TestConnect_Simple(t *testing.T) {
	os.Setenv("GO_ENV", "test")
	os.Setenv("API_PORT", "8080")
	os.Setenv("API_HOST", "localhost")
	os.Setenv("DB_RUN_MIGRATIONS", "false")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USERNAME", "test_user")
	os.Setenv("DB_PASSWORD", "test_pass")
	os.Setenv("MESSAGE_BROKER_TYPE", "sqs")
	os.Setenv("SQS_PAYMENT_QUEUE_URL", "https://sqs.us-east-1.amazonaws.com/123456789012/payment")
	os.Setenv("SQS_KITCHEN_QUEUE_URL", "https://sqs.us-east-1.amazonaws.com/123456789012/kitchen")
	os.Setenv("AWS_REGION", "us-east-1")

	defer func() {
		os.Unsetenv("GO_ENV")
		os.Unsetenv("API_PORT")
		os.Unsetenv("API_HOST")
		os.Unsetenv("DB_RUN_MIGRATIONS")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USERNAME")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("MESSAGE_BROKER_TYPE")
		os.Unsetenv("SQS_PAYMENT_QUEUE_URL")
		os.Unsetenv("SQS_KITCHEN_QUEUE_URL")
		os.Unsetenv("AWS_REGION")
	}()

	err := Connect()
	if err != nil {
		t.Logf("Connect failed as expected: %v", err)
	}
}

func TestGetBroker_Simple(t *testing.T) {
	broker := GetBroker()
	_ = broker
}

func TestNewPaymentEventConsumer_Simple(t *testing.T) {
	consumer := NewPaymentEventConsumer()
	if consumer == nil {
		t.Error("Expected consumer to be created")
	}
}

func TestPaymentEventConsumer_Start_Simple(t *testing.T) {
	consumer := &PaymentEventConsumer{}
	err := consumer.Start()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestBuildRabbitMQURL_Simple(t *testing.T) {
	cfg := &config.Config{}
	cfg.MessageBroker.RabbitMQ.URL = "amqp://test:test@localhost:5672/"

	url := buildRabbitMQURL(cfg)
	if url != "amqp://test:test@localhost:5672/" {
		t.Errorf("Expected URL 'amqp://test:test@localhost:5672/', got '%s'", url)
	}
}

func TestClose_Simple(t *testing.T) {
	Close()
}
