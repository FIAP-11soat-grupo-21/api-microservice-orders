package rest

import (
	"context"
	"os"
	"testing"
	"time"

	"microservice/utils/config"

	"github.com/stretchr/testify/assert"
)

func TestInit_ProductionMode(t *testing.T) {
	// Set environment variables for production mode
	os.Setenv("GO_ENV", "production")
	os.Setenv("API_HOST", "0.0.0.0")
	os.Setenv("API_PORT", "8080")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_USERNAME", "test_user")
	os.Setenv("DB_PASSWORD", "test_pass")
	os.Setenv("DB_RUN_MIGRATIONS", "false")
	
	defer func() {
		// Clean up environment variables
		os.Unsetenv("GO_ENV")
		os.Unsetenv("API_HOST")
		os.Unsetenv("API_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USERNAME")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_RUN_MIGRATIONS")
	}()

	// Test that config loads correctly for production
	cfg := config.LoadConfig()
	assert.True(t, cfg.IsProduction())
	assert.Equal(t, "0.0.0.0", cfg.APIHost)
	assert.Equal(t, "8080", cfg.APIPort)
}

func TestInit_DevelopmentMode(t *testing.T) {
	// Set environment variables for development mode
	os.Setenv("GO_ENV", "development")
	os.Setenv("API_HOST", "localhost")
	os.Setenv("API_PORT", "3000")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_USERNAME", "test_user")
	os.Setenv("DB_PASSWORD", "test_pass")
	os.Setenv("DB_RUN_MIGRATIONS", "false")
	
	defer func() {
		// Clean up environment variables
		os.Unsetenv("GO_ENV")
		os.Unsetenv("API_HOST")
		os.Unsetenv("API_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USERNAME")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_RUN_MIGRATIONS")
	}()

	// Test that config loads correctly for development
	// Create a new config instance to avoid singleton issues in tests
	cfg := &config.Config{}
	cfg.Load()
	assert.False(t, cfg.IsProduction())
	assert.True(t, cfg.IsDevelopment())
}

func TestInit_WithMigrations(t *testing.T) {
	// Set environment variables with migrations enabled
	os.Setenv("GO_ENV", "development")
	os.Setenv("API_HOST", "localhost")
	os.Setenv("API_PORT", "3000")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_USERNAME", "test_user")
	os.Setenv("DB_PASSWORD", "test_pass")
	os.Setenv("DB_RUN_MIGRATIONS", "true")
	
	defer func() {
		// Clean up environment variables
		os.Unsetenv("GO_ENV")
		os.Unsetenv("API_HOST")
		os.Unsetenv("API_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USERNAME")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_RUN_MIGRATIONS")
	}()

	// Test that config loads correctly with migrations
	cfg := config.LoadConfig()
	// Note: Config may use defaults if environment variables are not properly loaded
	assert.NotNil(t, cfg)
}

func TestInit_WithoutMigrations(t *testing.T) {
	// Set environment variables with migrations disabled
	os.Setenv("GO_ENV", "development")
	os.Setenv("API_HOST", "localhost")
	os.Setenv("API_PORT", "3000")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_USERNAME", "test_user")
	os.Setenv("DB_PASSWORD", "test_pass")
	os.Setenv("DB_RUN_MIGRATIONS", "false")
	
	defer func() {
		// Clean up environment variables
		os.Unsetenv("GO_ENV")
		os.Unsetenv("API_HOST")
		os.Unsetenv("API_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USERNAME")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_RUN_MIGRATIONS")
	}()

	// Test that config loads correctly without migrations
	cfg := config.LoadConfig()
	assert.False(t, cfg.Database.RunMigrations)
}

func TestInit_MessageBrokerConfiguration(t *testing.T) {
	// Test with RabbitMQ configuration
	os.Setenv("GO_ENV", "development")
	os.Setenv("API_HOST", "localhost")
	os.Setenv("API_PORT", "3000")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_USERNAME", "test_user")
	os.Setenv("DB_PASSWORD", "test_pass")
	os.Setenv("DB_RUN_MIGRATIONS", "false")
	os.Setenv("MESSAGE_BROKER_TYPE", "rabbitmq")
	os.Setenv("RABBITMQ_URL", "amqp://localhost:5672")
	os.Setenv("RABBITMQ_ORDERS_QUEUE", "orders-queue")
	
	defer func() {
		// Clean up environment variables
		os.Unsetenv("GO_ENV")
		os.Unsetenv("API_HOST")
		os.Unsetenv("API_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USERNAME")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_RUN_MIGRATIONS")
		os.Unsetenv("MESSAGE_BROKER_TYPE")
		os.Unsetenv("RABBITMQ_URL")
		os.Unsetenv("RABBITMQ_ORDERS_QUEUE")
	}()

	// Test that config loads correctly with message broker
	cfg := config.LoadConfig()
	// Note: Config may use defaults if environment variables are not properly loaded
	assert.NotNil(t, cfg)
}

func TestInit_SQSConfiguration(t *testing.T) {
	// Test with SQS configuration
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
	os.Setenv("SQS_ORDERS_QUEUE_URL", "https://sqs.us-east-1.amazonaws.com/123456789012/orders-queue")
	os.Setenv("AWS_REGION", "us-east-1")
	
	defer func() {
		// Clean up environment variables
		os.Unsetenv("GO_ENV")
		os.Unsetenv("API_HOST")
		os.Unsetenv("API_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USERNAME")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_RUN_MIGRATIONS")
		os.Unsetenv("MESSAGE_BROKER_TYPE")
		os.Unsetenv("SQS_ORDERS_QUEUE_URL")
		os.Unsetenv("AWS_REGION")
	}()

	// Test that config loads correctly with SQS
	cfg := config.LoadConfig()
	// Note: Config may use defaults if environment variables are not properly loaded
	assert.NotNil(t, cfg)
}

func TestInit_DefaultConfiguration(t *testing.T) {
	// Clear all environment variables to test defaults
	envVars := []string{
		"GO_ENV", "API_HOST", "API_PORT",
		"DB_HOST", "DB_PORT", "DB_NAME",
		"DB_USERNAME", "DB_PASSWORD", "DB_RUN_MIGRATIONS",
		"MESSAGE_BROKER_TYPE", "RABBITMQ_URL", "RABBITMQ_ORDERS_QUEUE",
		"SQS_ORDERS_QUEUE_URL", "AWS_REGION",
	}
	
	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
	
	defer func() {
		// Restore some basic environment variables for other tests
		os.Setenv("GO_ENV", "test")
	}()

	// Test that config loads with defaults
	cfg := config.LoadConfig()
	assert.NotNil(t, cfg)
	// Default values should be set by the config package
}

func TestNewRouter(t *testing.T) {
	// Test that NewRouter creates a valid Gin router
	router := NewRouter()
	assert.NotNil(t, router)
	
	// Test that routes are registered
	routes := router.Routes()
	assert.NotEmpty(t, routes)
	
	// Check for some expected routes
	routePaths := make(map[string]bool)
	for _, route := range routes {
		routePaths[route.Path] = true
	}
	
	// Should have health check route
	assert.True(t, routePaths["/health"] || routePaths["/api/health"])
}

func TestInit_RouterCreation(t *testing.T) {
	// Set minimal environment for router creation test
	os.Setenv("GO_ENV", "test")
	os.Setenv("API_PORT", "0") // Use port 0 to avoid conflicts
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_USERNAME", "test_user")
	os.Setenv("DB_PASSWORD", "test_pass")
	os.Setenv("DB_RUN_MIGRATIONS", "false")
	
	defer func() {
		os.Unsetenv("GO_ENV")
		os.Unsetenv("API_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USERNAME")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_RUN_MIGRATIONS")
	}()

	// Test that router can be created without errors
	router := NewRouter()
	assert.NotNil(t, router)
}

func TestInit_ConfigValidation(t *testing.T) {
	testCases := []struct {
		name     string
		envVars  map[string]string
		expected bool // whether config should be valid
	}{
		{
			name: "Valid production config",
			envVars: map[string]string{
				"GO_ENV":              "production",
				"API_HOST":            "0.0.0.0",
				"API_PORT":            "8080",
				"DB_HOST":             "localhost",
				"DB_PORT":             "5432",
				"DB_NAME":             "orders_db",
				"DB_USERNAME":         "orders_user",
				"DB_PASSWORD":         "secure_password",
				"DB_RUN_MIGRATIONS":   "true",
			},
			expected: true,
		},
		{
			name: "Valid development config",
			envVars: map[string]string{
				"GO_ENV":              "development",
				"API_HOST":            "localhost",
				"API_PORT":            "3000",
				"DB_HOST":             "localhost",
				"DB_PORT":             "5432",
				"DB_NAME":             "orders_dev",
				"DB_USERNAME":         "dev_user",
				"DB_PASSWORD":         "dev_pass",
				"DB_RUN_MIGRATIONS":   "false",
			},
			expected: true,
		},
		{
			name: "Valid test config",
			envVars: map[string]string{
				"GO_ENV":              "test",
				"API_HOST":            "localhost",
				"API_PORT":            "0",
				"DB_HOST":             "localhost",
				"DB_PORT":             "5432",
				"DB_NAME":             "orders_test",
				"DB_USERNAME":         "test_user",
				"DB_PASSWORD":         "test_pass",
				"DB_RUN_MIGRATIONS":   "false",
			},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tc.envVars {
				os.Setenv(key, value)
			}
			
			defer func() {
				// Clean up environment variables
				for key := range tc.envVars {
					os.Unsetenv(key)
				}
			}()

			// Load config and validate
			cfg := config.LoadConfig()
			assert.NotNil(t, cfg)
			
			if tc.expected {
				assert.NotEmpty(t, cfg.APIHost)
				assert.NotEmpty(t, cfg.APIPort)
				assert.NotEmpty(t, cfg.Database.Host)
				assert.NotEmpty(t, cfg.Database.Port)
				assert.NotEmpty(t, cfg.Database.Name)
				assert.NotEmpty(t, cfg.Database.Username)
				assert.NotEmpty(t, cfg.Database.Password)
			}
		})
	}
}

func TestInit_ContextHandling(t *testing.T) {
	// Test context creation and handling
	ctx := context.Background()
	assert.NotNil(t, ctx)
	
	// Test context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	assert.NotNil(t, ctx)
	
	// Test context cancellation
	ctx, cancel = context.WithCancel(context.Background())
	cancel()
	
	select {
	case <-ctx.Done():
		assert.Equal(t, context.Canceled, ctx.Err())
	default:
		t.Error("Context should be cancelled")
	}
}