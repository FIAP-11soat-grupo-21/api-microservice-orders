package data_source

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"microservice/infra/db/postgres/models"
)

func setupTestDBForStatus() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.OrderStatusModel{})

	statuses := []models.OrderStatusModel{
		{ID: "status-1", Name: "Pending"},
		{ID: "status-2", Name: "Confirmed"},
		{ID: "status-3", Name: "Preparing"},
	}

	for _, status := range statuses {
		db.Create(&status)
	}

	return db
}

func TestNewGormOrderStatusDataSource(t *testing.T) {
	ds := &GormOrderStatusDataSource{db: setupTestDBForStatus()}
	if ds.db == nil {
		t.Error("Expected db to be set")
	}
}

func TestGormOrderStatusDataSource_FindAll(t *testing.T) {
	db := setupTestDBForStatus()
	ds := &GormOrderStatusDataSource{db: db}

	statuses, err := ds.FindAll()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(statuses) != 3 {
		t.Errorf("Expected 3 statuses, got %d", len(statuses))
	}

	statusNames := make(map[string]bool)
	for _, status := range statuses {
		statusNames[status.Name] = true
	}

	expectedNames := []string{"Pending", "Confirmed", "Preparing"}
	for _, name := range expectedNames {
		if !statusNames[name] {
			t.Errorf("Expected status '%s' not found", name)
		}
	}
}

func TestGormOrderStatusDataSource_FindByID(t *testing.T) {
	db := setupTestDBForStatus()
	ds := &GormOrderStatusDataSource{db: db}

	status, err := ds.FindByID("status-1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if status.ID != "status-1" {
		t.Errorf("Expected status ID 'status-1', got '%s'", status.ID)
	}

	if status.Name != "Pending" {
		t.Errorf("Expected status name 'Pending', got '%s'", status.Name)
	}
}

func TestGormOrderStatusDataSource_FindByID_NotFound(t *testing.T) {
	db := setupTestDBForStatus()
	ds := &GormOrderStatusDataSource{db: db}

	_, err := ds.FindByID("nonexistent")
	if err == nil {
		t.Error("Expected error when status not found")
	}
}