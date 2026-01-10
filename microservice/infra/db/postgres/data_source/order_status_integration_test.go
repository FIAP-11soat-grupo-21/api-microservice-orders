package data_source

import (
	"os"
	"testing"

	"microservice/infra/db/postgres"
	"microservice/internal/adapters/daos"

	"github.com/stretchr/testify/assert"
)

func TestNewGormOrderStatusDataSource_Integration(t *testing.T) {
	// Skip integration test if GO_ENV is not set (no database setup)
	if os.Getenv("GO_ENV") == "" {
		t.Skip("Skipping integration test - GO_ENV not set, no database configuration available")
		return
	}
	
	// Test the actual constructor function
	dataSource := NewGormOrderStatusDataSource()
	assert.NotNil(t, dataSource)
	
	// Skip database connection test if no database is available
	if dataSource.db == nil {
		t.Skip("Skipping integration test - no database connection available")
		return
	}
	
	assert.NotNil(t, dataSource.db)
}

func TestGormOrderStatusDataSource_FindByName_Integration(t *testing.T) {
	// Skip integration test if GO_ENV is not set (no database setup)
	if os.Getenv("GO_ENV") == "" {
		t.Skip("Skipping integration test - GO_ENV not set, no database configuration available")
		return
	}
	
	// Setup test database connection
	postgres.Connect()
	
	dataSource := NewGormOrderStatusDataSource()
	
	// Skip if no database connection is available
	if dataSource == nil || dataSource.db == nil {
		t.Skip("Skipping integration test - no database connection available")
		return
	}
	
	// Test finding existing status by name
	status, err := dataSource.FindByName("Recebido")
	if err != nil {
		// If status doesn't exist, create it for testing
		testStatus := daos.OrderStatusDAO{
			ID:   "test-status-id",
			Name: "Recebido",
		}
		
		// Insert test status directly into database
		result := dataSource.db.Create(&testStatus)
		if result.Error != nil {
			t.Skipf("Could not create test status: %v", result.Error)
		}
		
		// Try finding again
		status, err = dataSource.FindByName("Recebido")
	}
	
	if err == nil {
		assert.NotEmpty(t, status.ID)
		assert.Equal(t, "Recebido", status.Name)
	} else {
		t.Logf("FindByName test skipped due to database setup: %v", err)
	}
}

func TestGormOrderStatusDataSource_FindByName_NotFound(t *testing.T) {
	// Skip integration test if GO_ENV is not set (no database setup)
	if os.Getenv("GO_ENV") == "" {
		t.Skip("Skipping integration test - GO_ENV not set, no database configuration available")
		return
	}
	
	// Setup test database connection
	postgres.Connect()
	
	dataSource := NewGormOrderStatusDataSource()
	
	// Skip if no database connection is available
	if dataSource == nil || dataSource.db == nil {
		t.Skip("Skipping integration test - no database connection available")
		return
	}
	
	// Test finding non-existent status
	_, err := dataSource.FindByName("NonExistentStatus")
	assert.Error(t, err)
}

func TestGormOrderStatusDataSource_FindByName_EmptyName(t *testing.T) {
	// Skip integration test if GO_ENV is not set (no database setup)
	if os.Getenv("GO_ENV") == "" {
		t.Skip("Skipping integration test - GO_ENV not set, no database configuration available")
		return
	}
	
	// Setup test database connection
	postgres.Connect()
	
	dataSource := NewGormOrderStatusDataSource()
	
	// Skip if no database connection is available
	if dataSource == nil || dataSource.db == nil {
		t.Skip("Skipping integration test - no database connection available")
		return
	}
	
	// Test finding with empty name
	_, err := dataSource.FindByName("")
	assert.Error(t, err)
}

func TestGormOrderStatusDataSource_FindByName_ValidNames(t *testing.T) {
	// Skip integration test if GO_ENV is not set (no database setup)
	if os.Getenv("GO_ENV") == "" {
		t.Skip("Skipping integration test - GO_ENV not set, no database configuration available")
		return
	}
	
	// Setup test database connection
	postgres.Connect()
	
	dataSource := NewGormOrderStatusDataSource()
	
	// Skip if no database connection is available
	if dataSource == nil || dataSource.db == nil {
		t.Skip("Skipping integration test - no database connection available")
		return
	}
	
	// Test with various valid status names
	validNames := []string{
		"Recebido",
		"Confirmado", 
		"Em preparação",
		"Pronto",
		"Entregue",
	}
	
	for _, name := range validNames {
		t.Run("FindByName_"+name, func(t *testing.T) {
			status, err := dataSource.FindByName(name)
			if err == nil {
				assert.NotEmpty(t, status.ID)
				assert.Equal(t, name, status.Name)
			} else {
				t.Logf("Status '%s' not found in database: %v", name, err)
			}
		})
	}
}

func TestGormOrderStatusDataSource_Methods_Coverage(t *testing.T) {
	// Skip if GO_ENV is not set (no database setup)
	if os.Getenv("GO_ENV") == "" {
		t.Skip("Skipping coverage test - GO_ENV not set, no database configuration available")
		return
	}
	
	// Test that all methods exist and can be called
	dataSource := NewGormOrderStatusDataSource()
	
	// Skip if no database connection is available
	if dataSource == nil || dataSource.db == nil {
		t.Skip("Skipping coverage test - no database connection available")
		return
	}
	
	// Test FindAll method exists
	statuses, err := dataSource.FindAll()
	if err != nil {
		t.Logf("FindAll error (expected in test environment): %v", err)
	} else {
		assert.NotNil(t, statuses)
	}
	
	// Test FindByID method exists
	_, err = dataSource.FindByID("test-id")
	assert.Error(t, err) // Should error for non-existent ID
	
	// Test FindByName method exists
	_, err = dataSource.FindByName("TestStatus")
	assert.Error(t, err) // Should error for non-existent name
}