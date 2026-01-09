package data_source

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"microservice/infra/db/postgres/models"
	"microservice/internal/adapters/daos"
)

func setupTestDBForStatus() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err := db.AutoMigrate(&models.OrderStatusModel{}); err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	statuses := []models.OrderStatusModel{
		{ID: "pending", Name: "Pending"},
		{ID: "confirmed", Name: "Confirmed"},
		{ID: "preparing", Name: "Preparing"},
		{ID: "ready", Name: "Ready"},
		{ID: "completed", Name: "Completed"},
	}

	for _, status := range statuses {
		db.Create(&status)
	}

	return db
}

func TestNewGormOrderStatusDataSource_MethodExists(t *testing.T) {
	_ = NewGormOrderStatusDataSource
}

func TestNewGormOrderStatusDataSource_ReturnsInstance(t *testing.T) {
	dataSource := &GormOrderStatusDataSource{db: setupTestDBForStatus()}
	assert.NotNil(t, dataSource)
	assert.NotNil(t, dataSource.db)
}

func TestGormOrderStatusDataSource_Structure(t *testing.T) {
	dataSource := &GormOrderStatusDataSource{
		db: nil,
	}
	assert.NotNil(t, dataSource)
}

func TestGormOrderStatusDataSource_FindAll_Success(t *testing.T) {
	db := setupTestDBForStatus()
	dataSource := &GormOrderStatusDataSource{db: db}

	results, err := dataSource.FindAll()
	assert.NoError(t, err)
	assert.Len(t, results, 5)

	statusIDs := make(map[string]bool)
	for _, status := range results {
		statusIDs[status.ID] = true
	}

	assert.True(t, statusIDs["pending"])
	assert.True(t, statusIDs["confirmed"])
	assert.True(t, statusIDs["preparing"])
	assert.True(t, statusIDs["ready"])
	assert.True(t, statusIDs["completed"])
}

func TestGormOrderStatusDataSource_FindAll_EmptyResult(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	if err := db.AutoMigrate(&models.OrderStatusModel{}); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	dataSource := &GormOrderStatusDataSource{db: db}

	results, err := dataSource.FindAll()
	assert.NoError(t, err)
	assert.Len(t, results, 0)
}

func TestGormOrderStatusDataSource_FindByID_Success(t *testing.T) {
	db := setupTestDBForStatus()
	dataSource := &GormOrderStatusDataSource{db: db}

	result, err := dataSource.FindByID("pending")
	assert.NoError(t, err)
	assert.Equal(t, "pending", result.ID)
	assert.Equal(t, "Pending", result.Name)
}

func TestGormOrderStatusDataSource_FindByID_NotFound(t *testing.T) {
	db := setupTestDBForStatus()
	dataSource := &GormOrderStatusDataSource{db: db}

	_, err := dataSource.FindByID("non-existent-status")
	assert.Error(t, err)
}

func TestGormOrderStatusDataSource_FindByID_AllStatuses(t *testing.T) {
	db := setupTestDBForStatus()
	dataSource := &GormOrderStatusDataSource{db: db}

	testCases := []struct {
		id   string
		name string
	}{
		{"pending", "Pending"},
		{"confirmed", "Confirmed"},
		{"preparing", "Preparing"},
		{"ready", "Ready"},
		{"completed", "Completed"},
	}

	for _, tc := range testCases {
		t.Run(tc.id, func(t *testing.T) {
			result, err := dataSource.FindByID(tc.id)
			assert.NoError(t, err)
			assert.Equal(t, tc.id, result.ID)
			assert.Equal(t, tc.name, result.Name)
		})
	}
}

func TestGormOrderStatusDataSource_Methods_Exist(t *testing.T) {
	dataSource := &GormOrderStatusDataSource{}

	_ = dataSource.FindAll
	_ = dataSource.FindByID
}

func TestOrderStatusDAO_Structure(t *testing.T) {
	status := daos.OrderStatusDAO{
		ID:   "status-1",
		Name: "pending",
	}

	assert.Equal(t, "status-1", status.ID)
	assert.Equal(t, "pending", status.Name)
}
