package data_source

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"microservice/infra/db/postgres/models"
	"microservice/internal/adapters/daos"
	"microservice/internal/adapters/dtos"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&models.OrderModel{}, &models.OrderItemModel{}, &models.OrderStatusModel{})
	if err != nil {
		panic("failed to migrate database")
	}

	status := models.OrderStatusModel{
		ID:   "status-1",
		Name: "Pending",
	}
	db.Create(&status)

	return db
}

func TestNewGormOrderDataSource(t *testing.T) {
	ds := &GormOrderDataSource{db: setupTestDB()}
	if ds.db == nil {
		t.Error("Expected db to be set")
	}
}

func TestGormOrderDataSource_Create(t *testing.T) {
	db := setupTestDB()
	ds := &GormOrderDataSource{db: db}

	order := daos.OrderDAO{
		ID:         "order-1",
		CustomerID: stringPtr("customer-1"),
		Amount:     25.0,
		Status: daos.OrderStatusDAO{
			ID:   "status-1",
			Name: "Pending",
		},
		Items: []daos.OrderItemDAO{
			{
				ID:        "item-1",
				OrderID:   "order-1",
				ProductID: "product-1",
				Quantity:  1,
				UnitPrice: 25.0,
			},
		},
		CreatedAt: time.Now(),
	}

	err := ds.Create(order)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	var count int64
	db.Model(&models.OrderModel{}).Where("id = ?", "order-1").Count(&count)
	if count != 1 {
		t.Errorf("Expected 1 order, got %d", count)
	}
}

func TestGormOrderDataSource_FindAll(t *testing.T) {
	db := setupTestDB()
	ds := &GormOrderDataSource{db: db}

	order := models.OrderModel{
		ID:         "order-1",
		CustomerID: stringPtr("customer-1"),
		Amount:     25.0,
		StatusID:   "status-1",
		CreatedAt:  time.Now(),
	}
	db.Create(&order)

	item := models.OrderItemModel{
		ID:        "item-1",
		OrderID:   "order-1",
		ProductID: "product-1",
		Quantity:  1,
		UnitPrice: 25.0,
	}
	db.Create(&item)

	orders, err := ds.FindAll(dtos.OrderFilterDTO{})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(orders) != 1 {
		t.Errorf("Expected 1 order, got %d", len(orders))
	}

	if orders[0].ID != "order-1" {
		t.Errorf("Expected order ID 'order-1', got '%s'", orders[0].ID)
	}
}

func TestGormOrderDataSource_FindAll_WithFilters(t *testing.T) {
	db := setupTestDB()
	ds := &GormOrderDataSource{db: db}

	now := time.Now()
	order := models.OrderModel{
		ID:         "order-1",
		CustomerID: stringPtr("customer-1"),
		Amount:     25.0,
		StatusID:   "status-1",
		CreatedAt:  now,
	}
	db.Create(&order)

	customerID := "customer-1"
	filter := dtos.OrderFilterDTO{
		CustomerID: &customerID,
	}

	orders, err := ds.FindAll(filter)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(orders) != 1 {
		t.Errorf("Expected 1 order, got %d", len(orders))
	}

	statusID := "status-1"
	filter = dtos.OrderFilterDTO{
		StatusID: &statusID,
	}

	orders, err = ds.FindAll(filter)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(orders) != 1 {
		t.Errorf("Expected 1 order, got %d", len(orders))
	}
}

func TestGormOrderDataSource_FindByID(t *testing.T) {
	db := setupTestDB()
	ds := &GormOrderDataSource{db: db}

	order := models.OrderModel{
		ID:         "order-1",
		CustomerID: stringPtr("customer-1"),
		Amount:     25.0,
		StatusID:   "status-1",
		CreatedAt:  time.Now(),
	}
	db.Create(&order)

	item := models.OrderItemModel{
		ID:        "item-1",
		OrderID:   "order-1",
		ProductID: "product-1",
		Quantity:  1,
		UnitPrice: 25.0,
	}
	db.Create(&item)

	result, err := ds.FindByID("order-1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != "order-1" {
		t.Errorf("Expected order ID 'order-1', got '%s'", result.ID)
	}

	if len(result.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(result.Items))
	}
}

func TestGormOrderDataSource_FindByID_NotFound(t *testing.T) {
	db := setupTestDB()
	ds := &GormOrderDataSource{db: db}

	_, err := ds.FindByID("nonexistent")
	if err == nil {
		t.Error("Expected error when order not found")
	}
}

func TestGormOrderDataSource_Update(t *testing.T) {
	db := setupTestDB()
	ds := &GormOrderDataSource{db: db}

	order := models.OrderModel{
		ID:         "order-1",
		CustomerID: stringPtr("customer-1"),
		Amount:     25.0,
		StatusID:   "status-1",
		CreatedAt:  time.Now(),
	}
	db.Create(&order)

	updatedOrder := daos.OrderDAO{
		ID:         "order-1",
		CustomerID: stringPtr("customer-1"),
		Amount:     30.0,
		Status: daos.OrderStatusDAO{
			ID:   "status-1",
			Name: "Pending",
		},
		CreatedAt: order.CreatedAt,
		UpdatedAt: &time.Time{},
	}

	err := ds.Update(updatedOrder)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	var updatedModel models.OrderModel
	db.First(&updatedModel, "id = ?", "order-1")
	if updatedModel.Amount != 30.0 {
		t.Errorf("Expected amount 30.0, got %f", updatedModel.Amount)
	}
}

func TestGormOrderDataSource_Delete(t *testing.T) {
	db := setupTestDB()
	ds := &GormOrderDataSource{db: db}

	order := models.OrderModel{
		ID:         "order-1",
		CustomerID: stringPtr("customer-1"),
		Amount:     25.0,
		StatusID:   "status-1",
		CreatedAt:  time.Now(),
	}
	db.Create(&order)

	item := models.OrderItemModel{
		ID:        "item-1",
		OrderID:   "order-1",
		ProductID: "product-1",
		Quantity:  1,
		UnitPrice: 25.0,
	}
	db.Create(&item)

	err := ds.Delete("order-1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	var count int64
	db.Model(&models.OrderModel{}).Where("id = ?", "order-1").Count(&count)
	if count != 0 {
		t.Errorf("Expected 0 orders, got %d", count)
	}

	db.Model(&models.OrderItemModel{}).Where("order_id = ?", "order-1").Count(&count)
	if count != 0 {
		t.Errorf("Expected 0 items, got %d", count)
	}
}

func stringPtr(s string) *string {
	return &s
}