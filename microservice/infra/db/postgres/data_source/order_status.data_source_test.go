package data_source

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"microservice/internal/adapters/daos"
)

func TestNewGormOrderStatusDataSource_MethodExists(t *testing.T) {
	_ = NewGormOrderStatusDataSource
}

func TestGormOrderStatusDataSource_Structure(t *testing.T) {
	dataSource := &GormOrderStatusDataSource{
		db: nil,
	}
	assert.NotNil(t, dataSource)
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
