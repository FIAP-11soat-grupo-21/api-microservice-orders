package seed

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"microservice/infra/db/postgres/models"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&models.OrderStatusModel{})
	if err != nil {
		panic("failed to migrate database")
	}

	return db
}

func TestSeedOrderStatus_NewDatabase(t *testing.T) {
	db := setupTestDB()

	SeedOrderStatus(db)

	var count int64
	db.Model(&models.OrderStatusModel{}).Count(&count)

	if count != 5 {
		t.Errorf("Expected 5 order statuses, got %d", count)
	}

	expectedStatuses := map[string]string{
		ORDER_STATUS_RECEIVED_ID:  "Recebido",
		ORDER_STATUS_CONFIRMED_ID: "Confirmado",
		ORDER_STATUS_PREPARING_ID: "Em preparação",
		ORDER_STATUS_READY_ID:     "Pronto",
		ORDER_STATUS_DELIVERED_ID: "Entregue",
	}

	for id, name := range expectedStatuses {
		var status models.OrderStatusModel
		err := db.Where("id = ?", id).First(&status).Error
		if err != nil {
			t.Errorf("Expected status with ID %s to exist", id)
		}

		if status.Name != name {
			t.Errorf("Expected status name '%s', got '%s'", name, status.Name)
		}
	}
}

func TestSeedOrderStatus_ExistingData(t *testing.T) {
	db := setupTestDB()

	existingStatus := models.OrderStatusModel{
		ID:   ORDER_STATUS_RECEIVED_ID,
		Name: "Recebido",
	}
	db.Create(&existingStatus)

	SeedOrderStatus(db)

	var count int64
	db.Model(&models.OrderStatusModel{}).Count(&count)

	if count != 5 {
		t.Errorf("Expected 5 order statuses, got %d", count)
	}

	var receivedStatuses []models.OrderStatusModel
	db.Where("id = ?", ORDER_STATUS_RECEIVED_ID).Find(&receivedStatuses)

	if len(receivedStatuses) != 1 {
		t.Errorf("Expected 1 'Recebido' status, got %d", len(receivedStatuses))
	}
}

func TestSeedOrderStatus_Constants(t *testing.T) {
	constants := []string{
		ORDER_STATUS_RECEIVED_ID,
		ORDER_STATUS_CONFIRMED_ID,
		ORDER_STATUS_PREPARING_ID,
		ORDER_STATUS_READY_ID,
		ORDER_STATUS_DELIVERED_ID,
	}

	for _, constant := range constants {
		if len(constant) != 36 {
			t.Errorf("Expected UUID format (36 chars), got %d chars for %s", len(constant), constant)
		}

		if constant == "" {
			t.Error("Expected non-empty constant")
		}
	}
}
