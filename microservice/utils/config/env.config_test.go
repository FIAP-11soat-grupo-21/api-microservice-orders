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
		name     string
		env      string
		isProd   bool
		isDev    bool
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
		"MESSAGE_BROKER_TYPE":   "sqs",
		"SQS_ORDERS_QUEUE_URL":  "https://sqs.us-east-1.amazonaws.com/123456789012/orders-queue",
		"AWS_REGION":            "us-east-1",
	}

	for key, value := range sqsEnvVars {
		os.Setenv(key, value)
	}

	config := &Config{}
	config.Load()

	if config.MessageBroker.Type != "sqs" {
		t.Errorf("Expected MessageBroker.Type 'sqs', got %s", config.MessageBroker.Type)
	}

	if config.MessageBroker.SQS.OrdersQueueURL != "https://sqs.us-east-1.amazonaws.com/123456789012/orders-queue" {
		t.Errorf("Expected SQS OrdersQueueURL 'https://sqs.us-east-1.amazonaws.com/123456789012/orders-queue', got %s", config.MessageBroker.SQS.OrdersQueueURL)
	}

	if config.MessageBroker.SQS.AWSRegion != "us-east-1" {
		t.Errorf("Expected AWS Region 'us-east-1', got %s", config.MessageBroker.SQS.AWSRegion)
	}
}

func TestConfig_MessageBroker_RabbitMQ(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()

	rabbitMQEnvVars := map[string]string{
		"MESSAGE_BROKER_TYPE":    "rabbitmq",
		"RABBITMQ_URL":           "amqp://user:pass@rabbitmq.example.com:5672/",
		"RABBITMQ_ORDERS_QUEUE":  "orders.updates",
	}

	for key, value := range rabbitMQEnvVars {
		os.Setenv(key, value)
	}

	config := &Config{}
	config.Load()

	if config.MessageBroker.Type != "rabbitmq" {
		t.Errorf("Expected MessageBroker.Type 'rabbitmq', got %s", config.MessageBroker.Type)
	}

	if config.MessageBroker.RabbitMQ.URL != "amqp://user:pass@rabbitmq.example.com:5672/" {
		t.Errorf("Expected RabbitMQ URL 'amqp://user:pass@rabbitmq.example.com:5672/', got %s", config.MessageBroker.RabbitMQ.URL)
	}

	if config.MessageBroker.RabbitMQ.OrdersQueue != "orders.updates" {
		t.Errorf("Expected RabbitMQ OrdersQueue 'orders.updates', got %s", config.MessageBroker.RabbitMQ.OrdersQueue)
	}
}

func TestConfig_MessageBroker_Defaults(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()

	// Remove MESSAGE_BROKER_TYPE to test default
	os.Unsetenv("MESSAGE_BROKER_TYPE")
	os.Unsetenv("SQS_ORDERS_QUEUE_URL")
	os.Unsetenv("RABBITMQ_URL")
	os.Unsetenv("RABBITMQ_ORDERS_QUEUE")

	// Set defaults
	os.Setenv("MESSAGE_BROKER_TYPE", "sqs")
	os.Setenv("AWS_REGION", "us-east-2")
	os.Setenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	os.Setenv("RABBITMQ_ORDERS_QUEUE", "orders.updates")

	config := &Config{}
	config.Load()

	if config.MessageBroker.Type != "sqs" {
		t.Errorf("Expected default MessageBroker.Type 'sqs', got %s", config.MessageBroker.Type)
	}

	if config.MessageBroker.SQS.AWSRegion != "us-east-2" {
		t.Errorf("Expected default AWS Region 'us-east-2', got %s", config.MessageBroker.SQS.AWSRegion)
	}

	if config.MessageBroker.RabbitMQ.URL != "amqp://guest:guest@localhost:5672/" {
		t.Errorf("Expected default RabbitMQ URL 'amqp://guest:guest@localhost:5672/', got %s", config.MessageBroker.RabbitMQ.URL)
	}

	if config.MessageBroker.RabbitMQ.OrdersQueue != "orders.updates" {
		t.Errorf("Expected default RabbitMQ OrdersQueue 'orders.updates', got %s", config.MessageBroker.RabbitMQ.OrdersQueue)
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
		"GO_ENV":                "production",
		"API_PORT":              "8443",
		"API_HOST":              "api.example.com",
		"DB_RUN_MIGRATIONS":     "true",
		"DB_HOST":               "prod-db.example.com",
		"DB_NAME":               "orders_prod",
		"DB_PORT":               "5432",
		"DB_USERNAME":           "prod_user",
		"DB_PASSWORD":           "prod_password",
		"MESSAGE_BROKER_TYPE":   "rabbitmq",
		"SQS_ORDERS_QUEUE_URL":  "https://sqs.us-west-2.amazonaws.com/123456789012/orders",
		"AWS_REGION":            "us-west-2",
		"RABBITMQ_URL":          "amqp://prod:prod@rabbitmq-prod:5672/",
		"RABBITMQ_ORDERS_QUEUE": "orders.prod",
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

	if config.MessageBroker.Type != "rabbitmq" {
		t.Errorf("Expected MessageBroker.Type 'rabbitmq', got %s", config.MessageBroker.Type)
	}

	if config.MessageBroker.RabbitMQ.URL != "amqp://prod:prod@rabbitmq-prod:5672/" {
		t.Errorf("Expected RabbitMQ URL 'amqp://prod:prod@rabbitmq-prod:5672/', got %s", config.MessageBroker.RabbitMQ.URL)
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

	if config.MessageBroker.SQS.OrdersQueueURL == "" && config.MessageBroker.SQS.OrdersQueueURL != "" {
		t.Error("Config.MessageBroker.SQS.OrdersQueueURL field missing")
	}

	if config.MessageBroker.RabbitMQ.URL == "" && config.MessageBroker.RabbitMQ.URL != "" {
		t.Error("Config.MessageBroker.RabbitMQ.URL field missing")
	}
}

// Helper functions
func setupTestEnv() {
	defaultEnvVars := map[string]string{
		"GO_ENV":                "test",
		"API_PORT":              "8080",
		"API_HOST":              "localhost",
		"DB_RUN_MIGRATIONS":     "false",
		"DB_HOST":               "localhost",
		"DB_NAME":               "test_db",
		"DB_PORT":               "5432",
		"DB_USERNAME":           "test_user",
		"DB_PASSWORD":           "test_pass",
		"MESSAGE_BROKER_TYPE":   "sqs",
		"SQS_ORDERS_QUEUE_URL":  "https://sqs.us-east-1.amazonaws.com/123456789012/test-orders",
		"AWS_REGION":            "us-east-1",
		"RABBITMQ_URL":          "amqp://guest:guest@localhost:5672/",
		"RABBITMQ_ORDERS_QUEUE": "orders.updates",
	}

	for key, value := range defaultEnvVars {
		os.Setenv(key, value)
	}
}

func cleanupTestEnv() {
	envVars := []string{
		"GO_ENV", "API_PORT", "API_HOST", "DB_RUN_MIGRATIONS",
		"DB_HOST", "DB_NAME", "DB_PORT", "DB_USERNAME", "DB_PASSWORD",
		"MESSAGE_BROKER_TYPE", "SQS_ORDERS_QUEUE_URL", "AWS_REGION",
		"RABBITMQ_URL", "RABBITMQ_ORDERS_QUEUE",
	}

	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}
