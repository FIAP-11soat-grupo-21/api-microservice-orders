package config

import (
	"os"
	"testing"
)

func TestConfig_IsProduction(t *testing.T) {
	tests := []struct {
		name     string
		goEnv    string
		expected bool
	}{
		{"production env", "production", true},
		{"development env", "development", false},
		{"staging env", "staging", false},
		{"empty env", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{GoEnv: tt.goEnv}
			result := cfg.IsProduction()
			if result != tt.expected {
				t.Errorf("IsProduction() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConfig_IsDevelopment(t *testing.T) {
	tests := []struct {
		name     string
		goEnv    string
		expected bool
	}{
		{"development env", "development", true},
		{"production env", "production", false},
		{"staging env", "staging", false},
		{"empty env", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{GoEnv: tt.goEnv}
			result := cfg.IsDevelopment()
			if result != tt.expected {
				t.Errorf("IsDevelopment() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetEnv_ExistingVariable(t *testing.T) {
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	value := os.Getenv("TEST_VAR")
	if value != "test_value" {
		t.Errorf("getEnv() = %v, want test_value", value)
	}
}

func TestConfig_Load_WithEnvVars(t *testing.T) {
	envVars := map[string]string{
		"GO_ENV":                 "development",
		"API_PORT":               "8080",
		"API_HOST":               "localhost",
		"DB_RUN_MIGRATIONS":      "true",
		"DB_HOST":                "localhost",
		"DB_NAME":                "testdb",
		"DB_PORT":                "5432",
		"DB_USERNAME":            "user",
		"DB_PASSWORD":            "pass",
		"RABBITMQ_HOST":          "localhost",
		"RABBITMQ_PORT":          "5672",
		"RABBITMQ_USER":          "guest",
		"RABBITMQ_PASSWORD":      "guest",
		"RABBITMQ_PAYMENT_QUEUE": "payments",
	}

	for key, value := range envVars {
		os.Setenv(key, value)
	}
	defer func() {
		for key := range envVars {
			os.Unsetenv(key)
		}
	}()

	cfg := &Config{}
	cfg.Load()

	if cfg.GoEnv != "development" {
		t.Errorf("Load() GoEnv = %v, want development", cfg.GoEnv)
	}
	if cfg.APIPort != "8080" {
		t.Errorf("Load() APIPort = %v, want 8080", cfg.APIPort)
	}
	if cfg.APIHost != "localhost" {
		t.Errorf("Load() APIHost = %v, want localhost", cfg.APIHost)
	}
	if !cfg.Database.RunMigrations {
		t.Error("Load() Database.RunMigrations = false, want true")
	}
	if cfg.Database.Host != "localhost" {
		t.Errorf("Load() Database.Host = %v, want localhost", cfg.Database.Host)
	}
	if cfg.Database.Name != "testdb" {
		t.Errorf("Load() Database.Name = %v, want testdb", cfg.Database.Name)
	}
	if cfg.MessageBroker.RabbitMQ.URL != "amqp://guest:guest@localhost:5672/" {
		t.Errorf("Load() MessageBroker.RabbitMQ.URL = %v, want amqp://guest:guest@localhost:5672/", cfg.MessageBroker.RabbitMQ.URL)
	}
}

func TestConfig_Load_MigrationsDisabled(t *testing.T) {
	envVars := map[string]string{
		"GO_ENV":                 "development",
		"API_PORT":               "8080",
		"API_HOST":               "localhost",
		"DB_RUN_MIGRATIONS":      "false",
		"DB_HOST":                "localhost",
		"DB_NAME":                "testdb",
		"DB_PORT":                "5432",
		"DB_USERNAME":            "user",
		"DB_PASSWORD":            "pass",
		"RABBITMQ_URL":           "amqp://guest:guest@localhost:5672/",
		"RABBITMQ_PAYMENT_QUEUE": "payments",
	}

	for key, value := range envVars {
		os.Setenv(key, value)
	}
	defer func() {
		for key := range envVars {
			os.Unsetenv(key)
		}
	}()

	cfg := &Config{}
	cfg.Load()

	if cfg.Database.RunMigrations {
		t.Error("Load() Database.RunMigrations = true, want false")
	}
}
func TestConfig_Load_WithDefaults(t *testing.T) {
	os.Setenv("GO_ENV", "test")
	os.Setenv("API_PORT", "8080")
	os.Setenv("API_HOST", "localhost")
	os.Setenv("DB_RUN_MIGRATIONS", "false")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USERNAME", "test_user")
	os.Setenv("DB_PASSWORD", "test_pass")

	os.Unsetenv("MESSAGE_BROKER_TYPE")
	os.Unsetenv("SQS_PAYMENT_QUEUE_URL")
	os.Unsetenv("SQS_KITCHEN_QUEUE_URL")
	os.Unsetenv("AWS_REGION")
	os.Unsetenv("RABBITMQ_URL")
	os.Unsetenv("RABBITMQ_PAYMENT_QUEUE")
	os.Unsetenv("RABBITMQ_KITCHEN_QUEUE")

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
	}()

	config := &Config{}
	config.Load()

	if config.MessageBroker.Type != "sqs" {
		t.Errorf("Expected default MessageBroker.Type 'sqs', got '%s'", config.MessageBroker.Type)
	}

	if config.MessageBroker.SQS.PaymentQueueURL != "" {
		t.Errorf("Expected empty default SQS PaymentQueueURL, got '%s'", config.MessageBroker.SQS.PaymentQueueURL)
	}

	if config.MessageBroker.SQS.KitchenQueueURL != "" {
		t.Errorf("Expected empty default SQS KitchenQueueURL, got '%s'", config.MessageBroker.SQS.KitchenQueueURL)
	}

	if config.MessageBroker.SQS.AWSRegion != "us-east-2" {
		t.Errorf("Expected default AWS Region 'us-east-2', got '%s'", config.MessageBroker.SQS.AWSRegion)
	}

	if config.MessageBroker.RabbitMQ.URL != "amqp://guest:guest@localhost:5672/" {
		t.Errorf("Expected default RabbitMQ URL, got '%s'", config.MessageBroker.RabbitMQ.URL)
	}

	if config.MessageBroker.RabbitMQ.PaymentQueue != "payment.confirmation" {
		t.Errorf("Expected default RabbitMQ PaymentQueue, got '%s'", config.MessageBroker.RabbitMQ.PaymentQueue)
	}

	if config.MessageBroker.RabbitMQ.KitchenQueue != "kitchen.orders" {
		t.Errorf("Expected default RabbitMQ KitchenQueue, got '%s'", config.MessageBroker.RabbitMQ.KitchenQueue)
	}
}

func TestGetEnv_WithDefault(t *testing.T) {
	value := getEnv("NON_EXISTENT_VAR", "default_value")
	if value != "default_value" {
		t.Errorf("Expected 'default_value', got '%s'", value)
	}
}

func TestGetEnv_ExistingVar(t *testing.T) {
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	value := getEnv("TEST_VAR", "default_value")
	if value != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", value)
	}
}
