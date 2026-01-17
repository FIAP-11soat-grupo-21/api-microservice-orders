package messaging

import (
	"os"
	"testing"
)

func TestConnect_Simple(t *testing.T) {
	os.Setenv("GO_ENV", "development")
	os.Setenv("API_HOST", "localhost")
	os.Setenv("API_PORT", "3000")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_USERNAME", "test_user")
	os.Setenv("DB_PASSWORD", "test_pass")
	os.Setenv("DB_RUN_MIGRATIONS", "false")
	os.Setenv("MESSAGE_BROKER_TYPE", "sqs")
	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	os.Setenv("AWS_ENDPOINT", "http://localhost:4566")
	os.Setenv("SQS_UPDATE_ORDER_STATUS_QUEUE_URL", "http://localhost:4566/000000000000/update-order-status-queue")
	os.Setenv("SQS_ORDER_ERROR_QUEUE_URL", "http://localhost:4566/000000000000/order-error-queue")
	os.Setenv("SNS_ORDER_ERROR_TOPIC_ARN", "arn:aws:sns:us-west-2:000000000000:order-error-topic")
	os.Setenv("SNS_ORDER_CREATED_TOPIC_ARN", "arn:aws:sns:us-west-2:000000000000:order-created-topic")

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
		os.Unsetenv("AWS_REGION")
		os.Unsetenv("AWS_ACCESS_KEY_ID")
		os.Unsetenv("AWS_SECRET_ACCESS_KEY")
		os.Unsetenv("AWS_ENDPOINT")
		os.Unsetenv("SQS_UPDATE_ORDER_STATUS_QUEUE_URL")
		os.Unsetenv("SQS_ORDER_ERROR_QUEUE_URL")
		os.Unsetenv("SNS_ORDER_ERROR_TOPIC_ARN")
		os.Unsetenv("SNS_ORDER_CREATED_TOPIC_ARN")
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

func TestClose_Simple(t *testing.T) {
	Close()
}
