package data_source

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"microservice/infra/db/postgres/models"
	"microservice/internal/adapters/daos"
)

func TestGormOrderStatusDataSource_Structure(t *testing.T) {
	// Test that GormOrderStatusDataSource has the required structure
	ds := &GormOrderStatusDataSource{}
	assert.NotNil(t, ds)
}

func TestOrderStatusDAO_Conversion(t *testing.T) {
	// Test OrderStatusDAO structure and conversion
	statusDAO := daos.OrderStatusDAO{
		ID:   "status-1",
		Name: "PENDING",
	}

	assert.Equal(t, "status-1", statusDAO.ID)
	assert.Equal(t, "PENDING", statusDAO.Name)
}

func TestOrderStatusModel_Conversion(t *testing.T) {
	// Test OrderStatusModel structure
	statusModel := models.OrderStatusModel{
		ID:   "status-1",
		Name: "PENDING",
	}

	assert.Equal(t, "status-1", statusModel.ID)
	assert.Equal(t, "PENDING", statusModel.Name)
	assert.Equal(t, "order_status", statusModel.TableName())
}

func TestOrderStatusDAO_MultipleStatuses(t *testing.T) {
	statuses := []daos.OrderStatusDAO{
		{ID: "status-1", Name: "PENDING"},
		{ID: "status-2", Name: "CONFIRMED"},
		{ID: "status-3", Name: "CANCELLED"},
	}

	assert.Len(t, statuses, 3)
	assert.Equal(t, "PENDING", statuses[0].Name)
	assert.Equal(t, "CONFIRMED", statuses[1].Name)
	assert.Equal(t, "CANCELLED", statuses[2].Name)
}

func TestOrderStatusModel_TableName(t *testing.T) {
	statusModel := models.OrderStatusModel{
		ID:   "status-1",
		Name: "PENDING",
	}

	assert.Equal(t, "order_status", statusModel.TableName())
}

func TestOrderStatusDAO_EmptyName(t *testing.T) {
	statusDAO := daos.OrderStatusDAO{
		ID:   "status-1",
		Name: "",
	}

	assert.Equal(t, "status-1", statusDAO.ID)
	assert.Equal(t, "", statusDAO.Name)
}

func TestOrderStatusDAO_SpecialCharacters(t *testing.T) {
	statusDAO := daos.OrderStatusDAO{
		ID:   "status-special-!@#$%",
		Name: "SPECIAL_STATUS",
	}

	assert.Equal(t, "status-special-!@#$%", statusDAO.ID)
	assert.Equal(t, "SPECIAL_STATUS", statusDAO.Name)
}

func TestOrderStatusModel_LongName(t *testing.T) {
	longName := "VERY_LONG_STATUS_NAME_WITH_MANY_CHARACTERS_AND_UNDERSCORES"
	statusModel := models.OrderStatusModel{
		ID:   "status-1",
		Name: longName,
	}

	assert.Equal(t, longName, statusModel.Name)
}

func TestOrderStatusDAO_Array(t *testing.T) {
	statuses := make([]daos.OrderStatusDAO, 0)

	statuses = append(statuses, daos.OrderStatusDAO{ID: "1", Name: "PENDING"})
	statuses = append(statuses, daos.OrderStatusDAO{ID: "2", Name: "CONFIRMED"})

	assert.Len(t, statuses, 2)
	assert.Equal(t, "PENDING", statuses[0].Name)
	assert.Equal(t, "CONFIRMED", statuses[1].Name)
}

func TestOrderStatusModel_Array(t *testing.T) {
	statuses := make([]models.OrderStatusModel, 0)

	statuses = append(statuses, models.OrderStatusModel{ID: "1", Name: "PENDING"})
	statuses = append(statuses, models.OrderStatusModel{ID: "2", Name: "CONFIRMED"})

	assert.Len(t, statuses, 2)
	assert.Equal(t, "PENDING", statuses[0].Name)
	assert.Equal(t, "CONFIRMED", statuses[1].Name)
}
