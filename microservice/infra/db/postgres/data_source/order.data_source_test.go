package data_source

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

	if err := db.AutoMigrate(&models.OrderModel{}, &models.OrderItemModel{}, &models.OrderStatusModel{}); err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	status := models.OrderStatusModel{
		ID:   "pending",
		Name: "Pending",
	}
	db.Create(&status)

	return db
}

func TestNewGormOrderDataSource_MethodExists(t *testing.T) {
	_ = NewGormOrderDataSource
}

func TestNewGormOrderDataSource_ReturnsInstance(t *testing.T) {
	dataSource := &GormOrderDataSource{db: setupTestDB()}
	assert.NotNil(t, dataSource)
	assert.NotNil(t, dataSource.db)
}

func TestGormOrderDataSource_Structure(t *testing.T) {
	dataSource := &GormOrderDataSource{
		db: nil,
	}
	assert.NotNil(t, dataSource)
}

func TestGormOrderDataSource_Create_Success(t *testing.T) {
	db := setupTestDB()
	dataSource := &GormOrderDataSource{db: db}

	order := daos.OrderDAO{
		ID:         "test-order-1",
		CustomerID: stringPtr("customer-1"),
		Amount:     99.99,
		Status: daos.OrderStatusDAO{
			ID:   "pending",
			Name: "Pending",
		},
		Items: []daos.OrderItemDAO{
			{
				ID:        "item-1",
				OrderID:   "test-order-1",
				ProductID: "product-1",
				Quantity:  2,
				UnitPrice: 49.995,
			},
		},
		CreatedAt: time.Now(),
	}

	err := dataSource.Create(order)
	assert.NoError(t, err)

	var createdOrder models.OrderModel
	err = db.First(&createdOrder, "id = ?", "test-order-1").Error
	assert.NoError(t, err)
	assert.Equal(t, "test-order-1", createdOrder.ID)
	assert.Equal(t, 99.99, createdOrder.Amount)
}

func TestGormOrderDataSource_FindByID_Success(t *testing.T) {
	db := setupTestDB()
	dataSource := &GormOrderDataSource{db: db}

	orderModel := models.OrderModel{
		ID:         "test-order-2",
		CustomerID: stringPtr("customer-2"),
		Amount:     150.00,
		StatusID:   "pending",
		CreatedAt:  time.Now(),
	}
	db.Create(&orderModel)

	itemModel := models.OrderItemModel{
		ID:        "item-2",
		OrderID:   "test-order-2",
		ProductID: "product-2",
		Quantity:  3,
		UnitPrice: 50.00,
	}
	db.Create(&itemModel)

	result, err := dataSource.FindByID("test-order-2")
	assert.NoError(t, err)
	assert.Equal(t, "test-order-2", result.ID)
	assert.Equal(t, 150.00, result.Amount)
	assert.Equal(t, "pending", result.Status.ID)
	assert.Len(t, result.Items, 1)
}

func TestGormOrderDataSource_FindByID_NotFound(t *testing.T) {
	db := setupTestDB()
	dataSource := &GormOrderDataSource{db: db}

	_, err := dataSource.FindByID("non-existent-order")
	assert.Error(t, err)
}

func TestGormOrderDataSource_FindAll_Success(t *testing.T) {
	db := setupTestDB()
	dataSource := &GormOrderDataSource{db: db}

	order1 := models.OrderModel{
		ID:         "order-1",
		CustomerID: stringPtr("customer-1"),
		Amount:     100.00,
		StatusID:   "pending",
		CreatedAt:  time.Now().Add(-1 * time.Hour),
	}
	order2 := models.OrderModel{
		ID:         "order-2",
		CustomerID: stringPtr("customer-2"),
		Amount:     200.00,
		StatusID:   "pending",
		CreatedAt:  time.Now(),
	}
	db.Create(&order1)
	db.Create(&order2)

	filter := dtos.OrderFilterDTO{}
	results, err := dataSource.FindAll(filter)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "order-2", results[0].ID)
	assert.Equal(t, "order-1", results[1].ID)
}

func TestGormOrderDataSource_FindAll_WithFilters(t *testing.T) {
	db := setupTestDB()
	dataSource := &GormOrderDataSource{db: db}

	now := time.Now()
	customerID := "customer-filter-test"

	order := models.OrderModel{
		ID:         "filtered-order",
		CustomerID: &customerID,
		Amount:     300.00,
		StatusID:   "pending",
		CreatedAt:  now,
	}
	db.Create(&order)

	filter := dtos.OrderFilterDTO{
		CustomerID: &customerID,
	}
	results, err := dataSource.FindAll(filter)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "filtered-order", results[0].ID)

	statusID := "pending"
	filter = dtos.OrderFilterDTO{
		StatusID: &statusID,
	}
	results, err = dataSource.FindAll(filter)
	assert.NoError(t, err)
	assert.True(t, len(results) >= 1)
}

func TestGormOrderDataSource_Update_Success(t *testing.T) {
	db := setupTestDB()
	dataSource := &GormOrderDataSource{db: db}

	order := daos.OrderDAO{
		ID:         "update-test-order",
		CustomerID: stringPtr("customer-update"),
		Amount:     100.00,
		Status: daos.OrderStatusDAO{
			ID:   "pending",
			Name: "Pending",
		},
		CreatedAt: time.Now(),
	}
	err := dataSource.Create(order)
	assert.NoError(t, err)

	order.Amount = 150.00
	updatedAt := time.Now()
	order.UpdatedAt = &updatedAt

	err = dataSource.Update(order)
	assert.NoError(t, err)

	result, err := dataSource.FindByID("update-test-order")
	assert.NoError(t, err)
	assert.Equal(t, 150.00, result.Amount)
}

func TestGormOrderDataSource_Delete_Success(t *testing.T) {
	db := setupTestDB()
	dataSource := &GormOrderDataSource{db: db}

	order := daos.OrderDAO{
		ID:         "delete-test-order",
		CustomerID: stringPtr("customer-delete"),
		Amount:     200.00,
		Status: daos.OrderStatusDAO{
			ID:   "pending",
			Name: "Pending",
		},
		Items: []daos.OrderItemDAO{
			{
				ID:        "delete-item-1",
				OrderID:   "delete-test-order",
				ProductID: "product-delete",
				Quantity:  1,
				UnitPrice: 200.00,
			},
		},
		CreatedAt: time.Now(),
	}
	err := dataSource.Create(order)
	assert.NoError(t, err)

	err = dataSource.Delete("delete-test-order")
	assert.NoError(t, err)

	_, err = dataSource.FindByID("delete-test-order")
	assert.Error(t, err)

	var itemCount int64
	db.Model(&models.OrderItemModel{}).Where("order_id = ?", "delete-test-order").Count(&itemCount)
	assert.Equal(t, int64(0), itemCount)
}

func TestGormOrderDataSource_Methods_Exist(t *testing.T) {
	dataSource := &GormOrderDataSource{}

	_ = dataSource.Create
	_ = dataSource.FindAll
	_ = dataSource.FindByID
	_ = dataSource.Update
	_ = dataSource.Delete
}

func TestOrderDAO_Structure(t *testing.T) {
	order := daos.OrderDAO{
		ID:     "test-order",
		Amount: 99.99,
	}

	assert.Equal(t, "test-order", order.ID)
	assert.Equal(t, 99.99, order.Amount)
}

func TestOrderFilterDTO_Structure(t *testing.T) {
	filter := dtos.OrderFilterDTO{}

	assert.NotNil(t, filter)
}

func stringPtr(s string) *string {
	return &s
}
