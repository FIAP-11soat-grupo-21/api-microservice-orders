package mocks

import "os"

// Helper functions
func SetupEnv() {
	defaultEnvVars := map[string]string{
		"GO_ENV":                            "test",
		"API_PORT":                          "8080",
		"API_HOST":                          "localhost",
		"API_GATEWAY_URL":                   "http://localhost:8081",
		"DB_RUN_MIGRATIONS":                 "false",
		"DB_HOST":                           "localhost",
		"DB_NAME":                           "test_db",
		"DB_PORT":                           "5432",
		"DB_USERNAME":                       "test_user",
		"DB_PASSWORD":                       "test_pass",
		"MESSAGE_BROKER_TYPE":               "sqs",
		"AWS_REGION":                        "us-west-2",
		"AWS_ACCESS_KEY_ID":                 "test",
		"AWS_SECRET_ACCESS_KEY":             "test",
		"AWS_ENDPOINT":                      "http://localhost:4566",
		"SQS_UPDATE_ORDER_STATUS_QUEUE_URL": "http://localhost:4566/000000000000/update-order-status-queue",
		"SQS_ORDER_ERROR_QUEUE_URL":         "http://localhost:4566/000000000000/order-error-queue",
		"SNS_ORDER_ERROR_TOPIC_ARN":         "arn:aws:sns:us-west-2:000000000000:order-error-topic",
		"SNS_ORDER_CREATED_TOPIC_ARN":       "arn:aws:sns:us-west-2:000000000000:order-created-topic",
	}

	for key, value := range defaultEnvVars {
		os.Setenv(key, value)
	}
}

func CleanupEnv() {
	envVars := []string{
		"GO_ENV", "API_PORT", "API_HOST", "API_GATEWAY_URL", "DB_RUN_MIGRATIONS",
		"DB_HOST", "DB_NAME", "DB_PORT", "DB_USERNAME", "DB_PASSWORD",
		"MESSAGE_BROKER_TYPE", "SQS_UPDATE_ORDER_STATUS_QUEUE_URL", "SQS_ORDER_ERROR_QUEUE_URL", "AWS_REGION",
		"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_ENDPOINT",
		"SNS_ORDER_ERROR_TOPIC_ARN", "SNS_ORDER_CREATED_TOPIC_ARN",
	}

	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}
