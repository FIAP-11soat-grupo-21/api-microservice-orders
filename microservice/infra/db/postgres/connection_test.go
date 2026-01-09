package postgres

import (
	"os"
	"sync"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetDB_Singleton(t *testing.T) {
	originalInstance := instance
	originalOnce := once
	defer func() {
		instance = originalInstance
		once = originalOnce
	}()

	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	dbConnection = testDB
	instance = nil
	once = sync.Once{}

	db1 := GetDB()
	db2 := GetDB()

	if db1 != db2 {
		t.Error("GetDB should return the same instance (singleton pattern)")
	}

	if db1 != testDB {
		t.Error("GetDB should return the configured database connection")
	}
}

func TestConnect_AlreadyConnected(t *testing.T) {
	originalConnection := dbConnection
	defer func() {
		dbConnection = originalConnection
	}()

	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	dbConnection = testDB

	Connect()

	if dbConnection != testDB {
		t.Error("Connect should not change existing database connection")
	}
}

func TestConnect_WithValidConfig(t *testing.T) {
	originalConnection := dbConnection
	defer func() {
		dbConnection = originalConnection
	}()

	os.Setenv("GO_ENV", "test")
	os.Setenv("API_PORT", "8080")
	os.Setenv("API_HOST", "localhost")
	os.Setenv("DB_RUN_MIGRATIONS", "false")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USERNAME", "test_user")
	os.Setenv("DB_PASSWORD", "test_pass")

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

	dbConnection = nil

	defer func() {
		if r := recover(); r != nil {
		}
	}()

	t.Skip("Skipping test that requires actual PostgreSQL connection")
}

func TestClose_NoConnection(t *testing.T) {
	originalConnection := dbConnection
	defer func() {
		dbConnection = originalConnection
	}()

	dbConnection = nil

	Close()
}

func TestClose_WithConnection(t *testing.T) {
	originalConnection := dbConnection
	defer func() {
		dbConnection = originalConnection
	}()

	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	dbConnection = testDB

	Close()
}

func TestRunMigrations_WithConnection(t *testing.T) {
	originalConnection := dbConnection
	defer func() {
		dbConnection = originalConnection
	}()

	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	dbConnection = testDB

	RunMigrations()

	if !testDB.Migrator().HasTable("orders") {
		t.Error("Expected orders table to be created")
	}

	if !testDB.Migrator().HasTable("order_items") {
		t.Error("Expected order_items table to be created")
	}

	if !testDB.Migrator().HasTable("order_status") {
		t.Error("Expected order_status table to be created")
	}
}