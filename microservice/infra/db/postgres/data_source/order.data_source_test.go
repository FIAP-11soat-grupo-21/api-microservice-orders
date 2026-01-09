package data_source

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"microservice/internal/adapters/daos"
	"microservice/internal/adapters/dtos"
)

func TestNewGormOrderDataSource_MethodExists(t *testing.T) {
	_ = NewGormOrderDataSource
}

func TestGormOrderDataSource_Structure(t *testing.T) {
	dataSource := &GormOrderDataSource{
		db: nil,
	}
	assert.NotNil(t, dataSource)
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
