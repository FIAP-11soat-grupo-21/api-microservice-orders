package config

import (
	"os"
	"testing"
)

func TestLoadConfig_Singleton(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()

	// Reset singleton for testing
	instance = nil
	once.Do(func() {})

	config1 := LoadConfig()
	config2 := LoadConfig()

	if config1 != config2 {
		t.Error("Expected LoadConfig to return the same instance (singleton pattern)")
	}
}

func TestConfig_Load_Success(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()

	config := &Config{}
	result := config.Load()

	if result != config {
		t.Error("Expected Load() to return the same config instance")
	}

	if config.GoEnv != "test" {
		t.Errorf("Expected GoEnv 'test', got %s", config.GoEnv)
	}

	if config.APIPort != "8080" {
		t.Errorf("Expected APIPort '8080', got %s", config.APIPort)
	}

	if config.APIHost != "localhost" {
		t.Errorf("Expected APIHost 'localhost', got %s", config.APIHost)
	}
}

func TestConfig_Environment(t *testing.T) {
	tests := []struct {
		name   string
		env    string
		isProd bool
		isDev  bool
	}{
		{"production", "production", true, false},
		{"development", "development", false, true},
		{"test", "test", false, false},
		{"staging", "staging", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTestEnv()
			defer cleanupTestEnv()

			os.Setenv("GO_ENV", tt.env)

			config := &Config{}
			config.Load()

			if config.IsProduction() != tt.isProd {
				t.Errorf("Expected IsProduction() = %v, got %v", tt.isProd, config.IsProduction())
			}

			if config.IsDevelopment() != tt.isDev {
				t.Errorf("Expected IsDevelopment() = %v, got %v", tt.isDev, config.IsDevelopment())
			}
		})
	}
}

func TestConfig_Database(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()

	dbEnvVars := map[string]string{
		"DB_RUN_MIGRATIONS": "true",
		"DB_HOST":           "db.example.com",
		"DB_NAME":           "orders_db",
		"DB_PORT":           "5432",
		"DB_USERNAME":       "orders_user",
		"DB_PASSWORD":       "secure_password",
	}

	for key, value := range dbEnvVars {
		os.Setenv(key, value)
	}

	config := &Config{}
	config.Load()

	if !config.Database.RunMigrations {
		t.Error("Expected Database.RunMigrations to be true")
	}

	if config.Database.Host != "db.example.com" {
		t.Errorf("Expected Database.Host 'db.example.com', got %s", config.Database.Host)
	}

	if config.Database.Name != "orders_db" {
		t.Errorf("Expected Database.Name 'orders_db', got %s", config.Database.Name)
	}

	if config.Database.Port != "5432" {
		t.Errorf("Expected Database.Port '5432', got %s", config.Database.Port)
	}

	if config.Database.Username != "orders_user" {
		t.Errorf("Expected Database.Username 'orders_user', got %s", config.Database.Username)
	}

	if config.Database.Password != "secure_password" {
		t.Errorf("Expected Database.Password 'secure_password', got %s", config.Database.Password)
	}
}

func TestConfig_Database_RunMigrations_False(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()

	os.Setenv("DB_RUN_MIGRATIONS", "false")

	config := &Config{}
	config.Load()

	if config.Database.RunMigrations {
		t.Error("Expected Database.RunMigrations to be false")
	}
}

func TestConfig_MessageBroker_SQS(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()

	sqsEnvVars := map[string]string{
		"MESSAGE_BROKER_TYPE":               "sqs",
		"SQS_UPDATE_ORDER_STATUS_QUEUE_URL": "http://localhost:4566/000000000000/update-order-status-queue",
		"AWS_REGION":                        "us-east-1",
	}

	for key, value := range sqsEnvVars {
		os.Setenv(key, value)
	}

	config := &Config{}
	config.Load()

	if config.MessageBroker.Type != "sqs" {
		t.Errorf("Expected MessageBroker.Type 'sqs', got %s", config.MessageBroker.Type)
	}

	if config.MessageBroker.SQS.UpdateOrderStatusQueueURL != "http://localhost:4566/000000000000/update-order-status-queue" {
		t.Errorf("Expected SQS UpdateOrderStatusQueueURL 'http://localhost:4566/000000000000/update-order-status-queue', got %s", config.MessageBroker.SQS.UpdateOrderStatusQueueURL)
	}

	if config.AWS.Region != "us-east-1" {
		t.Errorf("Expected AWS Region 'us-east-1', got %s", config.AWS.Region)
	}
}

func TestConfig_MessageBroker_Defaults(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()

	// Remove MESSAGE_BROKER_TYPE to test default
	os.Unsetenv("MESSAGE_BROKER_TYPE")
	os.Unsetenv("AWS_REGION")

	// Set defaults
	os.Setenv("MESSAGE_BROKER_TYPE", "sqs")
	os.Setenv("AWS_REGION", "us-east-2")

	config := &Config{}
	config.Load()

	if config.MessageBroker.Type != "sqs" {
		t.Errorf("Expected default MessageBroker.Type 'sqs', got %s", config.MessageBroker.Type)
	}

	if config.AWS.Region != "us-east-2" {
		t.Errorf("Expected default AWS Region 'us-east-2', got %s", config.AWS.Region)
	}
}

func TestConfig_API_Configuration(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()

	apiEnvVars := map[string]string{
		"API_PORT": "3000",
		"API_HOST": "0.0.0.0",
	}

	for key, value := range apiEnvVars {
		os.Setenv(key, value)
	}

	config := &Config{}
	config.Load()

	if config.APIPort != "3000" {
		t.Errorf("Expected APIPort '3000', got %s", config.APIPort)
	}

	if config.APIHost != "0.0.0.0" {
		t.Errorf("Expected APIHost '0.0.0.0', got %s", config.APIHost)
	}
}

func TestConfig_AllFields(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()

	allEnvVars := map[string]string{
		"GO_ENV":                            "production",
		"API_PORT":                          "8443",
		"API_HOST":                          "api.example.com",
		"DB_RUN_MIGRATIONS":                 "true",
		"DB_HOST":                           "prod-db.example.com",
		"DB_NAME":                           "orders_prod",
		"DB_PORT":                           "5432",
		"DB_USERNAME":                       "prod_user",
		"DB_PASSWORD":                       "prod_password",
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

	for key, value := range allEnvVars {
		os.Setenv(key, value)
	}

	config := &Config{}
	config.Load()

	// Verify all fields
	if config.GoEnv != "production" {
		t.Errorf("Expected GoEnv 'production', got %s", config.GoEnv)
	}

	if config.APIPort != "8443" {
		t.Errorf("Expected APIPort '8443', got %s", config.APIPort)
	}

	if config.APIHost != "api.example.com" {
		t.Errorf("Expected APIHost 'api.example.com', got %s", config.APIHost)
	}

	if !config.Database.RunMigrations {
		t.Error("Expected Database.RunMigrations to be true")
	}

	if config.Database.Host != "prod-db.example.com" {
		t.Errorf("Expected Database.Host 'prod-db.example.com', got %s", config.Database.Host)
	}

	if config.Database.Name != "orders_prod" {
		t.Errorf("Expected Database.Name 'orders_prod', got %s", config.Database.Name)
	}

	if config.MessageBroker.Type != "sqs" {
		t.Errorf("Expected MessageBroker.Type 'sqs', got %s", config.MessageBroker.Type)
	}

	if config.AWS.Region != "us-west-2" {
		t.Errorf("Expected AWS Region 'us-west-2', got %s", config.AWS.Region)
	}

	if config.AWS.AccessKeyID != "test" {
		t.Errorf("Expected AWS AccessKeyID 'test', got %s", config.AWS.AccessKeyID)
	}

	if config.AWS.SecretAccessKey != "test" {
		t.Errorf("Expected AWS SecretAccessKey 'test', got %s", config.AWS.SecretAccessKey)
	}

	if config.AWS.Endpoint != "http://localhost:4566" {
		t.Errorf("Expected AWS Endpoint 'http://localhost:4566', got %s", config.AWS.Endpoint)
	}

	if config.MessageBroker.SQS.UpdateOrderStatusQueueURL != "http://localhost:4566/000000000000/update-order-status-queue" {
		t.Errorf("Expected SQS UpdateOrderStatusQueueURL 'http://localhost:4566/000000000000/update-order-status-queue', got %s", config.MessageBroker.SQS.UpdateOrderStatusQueueURL)
	}

	if config.MessageBroker.SQS.OrderErrorQueueURL != "http://localhost:4566/000000000000/order-error-queue" {
		t.Errorf("Expected SQS OrderErrorQueueURL 'http://localhost:4566/000000000000/order-error-queue', got %s", config.MessageBroker.SQS.OrderErrorQueueURL)
	}

	if config.MessageBroker.SNS.OrderErrorTopicARN != "arn:aws:sns:us-west-2:000000000000:order-error-topic" {
		t.Errorf("Expected SNS OrderErrorTopicARN 'arn:aws:sns:us-west-2:000000000000:order-error-topic', got %s", config.MessageBroker.SNS.OrderErrorTopicARN)
	}

	if config.MessageBroker.SNS.OrderCreatedTopicARN != "arn:aws:sns:us-west-2:000000000000:order-created-topic" {
		t.Errorf("Expected SNS OrderCreatedTopicARN 'arn:aws:sns:us-west-2:000000000000:order-created-topic', got %s", config.MessageBroker.SNS.OrderCreatedTopicARN)
	}

	if config.IsProduction() != true {
		t.Error("Expected IsProduction() to return true")
	}

	if config.IsDevelopment() != false {
		t.Error("Expected IsDevelopment() to return false")
	}
}

func TestConfig_Structure(t *testing.T) {
	config := &Config{}

	// Verify struct fields exist
	if config.GoEnv == "" && config.GoEnv != "" {
		t.Error("Config.GoEnv field missing")
	}

	if config.APIPort == "" && config.APIPort != "" {
		t.Error("Config.APIPort field missing")
	}

	if config.APIHost == "" && config.APIHost != "" {
		t.Error("Config.APIHost field missing")
	}

	// Verify nested structs
	if config.Database.Host == "" && config.Database.Host != "" {
		t.Error("Config.Database.Host field missing")
	}

	if config.MessageBroker.Type == "" && config.MessageBroker.Type != "" {
		t.Error("Config.MessageBroker.Type field missing")
	}

	if config.MessageBroker.SQS.UpdateOrderStatusQueueURL == "" && config.MessageBroker.SQS.UpdateOrderStatusQueueURL != "" {
		t.Error("Config.MessageBroker.SQS.UpdateOrderStatusQueueURL field missing")
	}
}

// Helper functions
func setupTestEnv() {
	defaultEnvVars := map[string]string{
		"GO_ENV":                            "test",
		"API_PORT":                          "8080",
		"API_HOST":                          "localhost",
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

func cleanupTestEnv() {
	envVars := []string{
		"GO_ENV", "API_PORT", "API_HOST", "DB_RUN_MIGRATIONS",
		"DB_HOST", "DB_NAME", "DB_PORT", "DB_USERNAME", "DB_PASSWORD",
		"MESSAGE_BROKER_TYPE", "SQS_UPDATE_ORDER_STATUS_QUEUE_URL", "SQS_ORDER_ERROR_QUEUE_URL", "AWS_REGION",
		"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_ENDPOINT",
		"SNS_ORDER_ERROR_TOPIC_ARN", "SNS_ORDER_CREATED_TOPIC_ARN",
	}

	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}
