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
		"GO_ENV":                "development",
		"API_PORT":              "8080",
		"API_HOST":              "localhost",
		"DB_RUN_MIGRATIONS":     "true",
		"DB_HOST":               "localhost",
		"DB_NAME":               "testdb",
		"DB_PORT":               "5432",
		"DB_USERNAME":           "user",
		"DB_PASSWORD":           "pass",
		"RABBITMQ_HOST":         "localhost",
		"RABBITMQ_PORT":         "5672",
		"RABBITMQ_USER":         "guest",
		"RABBITMQ_PASSWORD":     "guest",
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
		"GO_ENV":                "development",
		"API_PORT":              "8080",
		"API_HOST":              "localhost",
		"DB_RUN_MIGRATIONS":     "false",
		"DB_HOST":               "localhost",
		"DB_NAME":               "testdb",
		"DB_PORT":               "5432",
		"DB_USERNAME":           "user",
		"DB_PASSWORD":           "pass",
		"RABBITMQ_URL":          "amqp://guest:guest@localhost:5672/",
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
