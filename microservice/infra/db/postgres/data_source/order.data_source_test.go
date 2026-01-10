package data_source

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"microservice/infra/db/postgres/models"
	"microservice/internal/adapters/daos"
)

// ============================================================================
// Tests for NewGormOrderDataSource
// ============================================================================

func TestNewGormOrderDataSource(t *testing.T) {
	ds := NewGormOrderDataSource()
	assert.NotNil(t, ds)
	// Note: db might be nil if postgres.GetDB() returns nil in test environment
	// This is expected behavior in unit tests without database setup
}

// ============================================================================
// Tests for FromDAOToModel conversion
// ============================================================================

func TestFromDAOToModel_BasicConversion(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	dao := daos.OrderDAO{
		ID:         "order-1",
		CustomerID: &customerID,
		Amount:     100.0,
		Status: daos.OrderStatusDAO{
			ID:   "status-1",
			Name: "PENDING",
		},
		Items: []daos.OrderItemDAO{
			{
				ID:        "item-1",
				OrderID:   "order-1",
				ProductID: "product-1",
				Quantity:  2,
				UnitPrice: 50.0,
			},
		},
		CreatedAt: now,
	}

	model := FromDAOToModel(dao)

	assert.Equal(t, "order-1", model.ID)
	assert.Equal(t, customerID, *model.CustomerID)
	assert.Equal(t, 100.0, model.Amount)
	assert.Equal(t, "status-1", model.StatusID)
	assert.Equal(t, "PENDING", model.Status.Name)
	assert.Len(t, model.Items, 1)
	assert.Equal(t, "product-1", model.Items[0].ProductID)
	assert.Equal(t, 2, model.Items[0].Quantity)
	assert.Equal(t, 50.0, model.Items[0].UnitPrice)
}

func TestFromDAOToModel_NilCustomerIDUnit(t *testing.T) {
	now := time.Now()

	dao := daos.OrderDAO{
		ID:         "order-1",
		CustomerID: nil,
		Amount:     100.0,
		Status: daos.OrderStatusDAO{
			ID:   "status-1",
			Name: "PENDING",
		},
		Items:     []daos.OrderItemDAO{},
		CreatedAt: now,
	}

	model := FromDAOToModel(dao)

	assert.Nil(t, model.CustomerID)
	assert.Equal(t, "order-1", model.ID)
}

func TestFromDAOToModel_EmptyItems(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	dao := daos.OrderDAO{
		ID:         "order-1",
		CustomerID: &customerID,
		Amount:     100.0,
		Status: daos.OrderStatusDAO{
			ID:   "status-1",
			Name: "PENDING",
		},
		Items:     []daos.OrderItemDAO{},
		CreatedAt: now,
	}

	model := FromDAOToModel(dao)

	assert.Len(t, model.Items, 0)
}

func TestFromDAOToModel_MultipleItems(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	dao := daos.OrderDAO{
		ID:         "order-1",
		CustomerID: &customerID,
		Amount:     300.0,
		Status: daos.OrderStatusDAO{
			ID:   "status-1",
			Name: "PENDING",
		},
		Items: []daos.OrderItemDAO{
			{ID: "item-1", OrderID: "order-1", ProductID: "product-1", Quantity: 2, UnitPrice: 50.0},
			{ID: "item-2", OrderID: "order-1", ProductID: "product-2", Quantity: 1, UnitPrice: 100.0},
			{ID: "item-3", OrderID: "order-1", ProductID: "product-3", Quantity: 5, UnitPrice: 20.0},
		},
		CreatedAt: now,
	}

	model := FromDAOToModel(dao)

	assert.Len(t, model.Items, 3)
	assert.Equal(t, "product-1", model.Items[0].ProductID)
	assert.Equal(t, "product-2", model.Items[1].ProductID)
	assert.Equal(t, "product-3", model.Items[2].ProductID)
}

func TestFromDAOToModel_WithUpdatedAt(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()
	updatedAt := now.Add(time.Hour)

	dao := daos.OrderDAO{
		ID:         "order-1",
		CustomerID: &customerID,
		Amount:     100.0,
		Status: daos.OrderStatusDAO{
			ID:   "status-1",
			Name: "PENDING",
		},
		Items:     []daos.OrderItemDAO{},
		CreatedAt: now,
		UpdatedAt: &updatedAt,
	}

	model := FromDAOToModel(dao)

	assert.NotNil(t, model.UpdatedAt)
	assert.Equal(t, updatedAt, *model.UpdatedAt)
}

func TestFromDAOToModel_NilUpdatedAt(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	dao := daos.OrderDAO{
		ID:         "order-1",
		CustomerID: &customerID,
		Amount:     100.0,
		Status: daos.OrderStatusDAO{
			ID:   "status-1",
			Name: "PENDING",
		},
		Items:     []daos.OrderItemDAO{},
		CreatedAt: now,
		UpdatedAt: nil,
	}

	model := FromDAOToModel(dao)

	assert.Nil(t, model.UpdatedAt)
}

func TestFromDAOToModel_DifferentAmounts(t *testing.T) {
	tests := []struct {
		name   string
		amount float64
	}{
		{"Zero amount", 0.0},
		{"Small amount", 0.01},
		{"Large amount", 999999.99},
		{"Decimal amount", 123.45},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customerID := "customer-123"
			now := time.Now()

			dao := daos.OrderDAO{
				ID:         "order-1",
				CustomerID: &customerID,
				Amount:     tt.amount,
				Status: daos.OrderStatusDAO{
					ID:   "status-1",
					Name: "PENDING",
				},
				Items:     []daos.OrderItemDAO{},
				CreatedAt: now,
			}

			model := FromDAOToModel(dao)

			assert.Equal(t, tt.amount, model.Amount)
		})
	}
}

// ============================================================================
// Tests for FromModelToDAO conversion
// ============================================================================

func TestFromModelToDAO_BasicConversion(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	model := models.OrderModel{
		ID:         "order-1",
		CustomerID: &customerID,
		Amount:     100.0,
		StatusID:   "status-1",
		Status: models.OrderStatusModel{
			ID:   "status-1",
			Name: "PENDING",
		},
		Items: []models.OrderItemModel{
			{
				ID:        "item-1",
				OrderID:   "order-1",
				ProductID: "product-1",
				Quantity:  2,
				UnitPrice: 50.0,
			},
		},
		CreatedAt: now,
	}

	dao := FromModelToDAO(model)

	assert.Equal(t, "order-1", dao.ID)
	assert.Equal(t, customerID, *dao.CustomerID)
	assert.Equal(t, 100.0, dao.Amount)
	assert.Equal(t, "PENDING", dao.Status.Name)
	assert.Len(t, dao.Items, 1)
}

func TestFromModelToDAO_NilCustomerIDUnit(t *testing.T) {
	now := time.Now()

	model := models.OrderModel{
		ID:         "order-1",
		CustomerID: nil,
		Amount:     100.0,
		StatusID:   "status-1",
		Status: models.OrderStatusModel{
			ID:   "status-1",
			Name: "PENDING",
		},
		Items:     []models.OrderItemModel{},
		CreatedAt: now,
	}

	dao := FromModelToDAO(model)

	assert.Nil(t, dao.CustomerID)
}

func TestFromModelToDAO_EmptyItems(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	model := models.OrderModel{
		ID:         "order-1",
		CustomerID: &customerID,
		Amount:     100.0,
		StatusID:   "status-1",
		Status: models.OrderStatusModel{
			ID:   "status-1",
			Name: "PENDING",
		},
		Items:     []models.OrderItemModel{},
		CreatedAt: now,
	}

	dao := FromModelToDAO(model)

	assert.Len(t, dao.Items, 0)
}

func TestFromModelToDAO_MultipleItems(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	model := models.OrderModel{
		ID:         "order-1",
		CustomerID: &customerID,
		Amount:     300.0,
		StatusID:   "status-1",
		Status: models.OrderStatusModel{
			ID:   "status-1",
			Name: "PENDING",
		},
		Items: []models.OrderItemModel{
			{ID: "item-1", OrderID: "order-1", ProductID: "product-1", Quantity: 2, UnitPrice: 50.0},
			{ID: "item-2", OrderID: "order-1", ProductID: "product-2", Quantity: 1, UnitPrice: 100.0},
			{ID: "item-3", OrderID: "order-1", ProductID: "product-3", Quantity: 5, UnitPrice: 20.0},
		},
		CreatedAt: now,
	}

	dao := FromModelToDAO(model)

	assert.Len(t, dao.Items, 3)
	assert.Equal(t, "product-1", dao.Items[0].ProductID)
	assert.Equal(t, "product-2", dao.Items[1].ProductID)
	assert.Equal(t, "product-3", dao.Items[2].ProductID)
}

func TestFromModelToDAO_WithUpdatedAt(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()
	updatedAt := now.Add(time.Hour)

	model := models.OrderModel{
		ID:         "order-1",
		CustomerID: &customerID,
		Amount:     100.0,
		StatusID:   "status-1",
		Status: models.OrderStatusModel{
			ID:   "status-1",
			Name: "PENDING",
		},
		Items:     []models.OrderItemModel{},
		CreatedAt: now,
		UpdatedAt: &updatedAt,
	}

	dao := FromModelToDAO(model)

	assert.NotNil(t, dao.UpdatedAt)
	assert.Equal(t, updatedAt, *dao.UpdatedAt)
}

func TestFromModelToDAO_NilUpdatedAt(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	model := models.OrderModel{
		ID:         "order-1",
		CustomerID: &customerID,
		Amount:     100.0,
		StatusID:   "status-1",
		Status: models.OrderStatusModel{
			ID:   "status-1",
			Name: "PENDING",
		},
		Items:     []models.OrderItemModel{},
		CreatedAt: now,
		UpdatedAt: nil,
	}

	dao := FromModelToDAO(model)

	assert.Nil(t, dao.UpdatedAt)
}

// ============================================================================
// Tests for FromModelArrayToDAOArray conversion
// ============================================================================

func TestFromModelArrayToDAOArray_SingleOrder(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	models := []models.OrderModel{
		{
			ID:         "order-1",
			CustomerID: &customerID,
			Amount:     100.0,
			StatusID:   "status-1",
			Status: models.OrderStatusModel{
				ID:   "status-1",
				Name: "PENDING",
			},
			Items:     []models.OrderItemModel{},
			CreatedAt: now,
		},
	}

	daos := FromModelArrayToDAOArray(models)

	assert.Len(t, daos, 1)
	assert.Equal(t, "order-1", daos[0].ID)
}

func TestFromModelArrayToDAOArray_MultipleOrders(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	models := []models.OrderModel{
		{
			ID:         "order-1",
			CustomerID: &customerID,
			Amount:     100.0,
			StatusID:   "status-1",
			Status: models.OrderStatusModel{
				ID:   "status-1",
				Name: "PENDING",
			},
			Items:     []models.OrderItemModel{},
			CreatedAt: now,
		},
		{
			ID:         "order-2",
			CustomerID: &customerID,
			Amount:     200.0,
			StatusID:   "status-2",
			Status: models.OrderStatusModel{
				ID:   "status-2",
				Name: "CONFIRMED",
			},
			Items:     []models.OrderItemModel{},
			CreatedAt: now,
		},
		{
			ID:         "order-3",
			CustomerID: &customerID,
			Amount:     300.0,
			StatusID:   "status-3",
			Status: models.OrderStatusModel{
				ID:   "status-3",
				Name: "CANCELLED",
			},
			Items:     []models.OrderItemModel{},
			CreatedAt: now,
		},
	}

	daos := FromModelArrayToDAOArray(models)

	assert.Len(t, daos, 3)
	assert.Equal(t, "order-1", daos[0].ID)
	assert.Equal(t, "order-2", daos[1].ID)
	assert.Equal(t, "order-3", daos[2].ID)
	assert.Equal(t, "PENDING", daos[0].Status.Name)
	assert.Equal(t, "CONFIRMED", daos[1].Status.Name)
	assert.Equal(t, "CANCELLED", daos[2].Status.Name)
}

func TestFromModelArrayToDAOArray_EmptyUnit(t *testing.T) {
	models := []models.OrderModel{}
	daos := FromModelArrayToDAOArray(models)

	assert.Len(t, daos, 0)
	assert.NotNil(t, daos)
}

func TestFromModelArrayToDAOArray_WithItems(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	models := []models.OrderModel{
		{
			ID:         "order-1",
			CustomerID: &customerID,
			Amount:     100.0,
			StatusID:   "status-1",
			Status: models.OrderStatusModel{
				ID:   "status-1",
				Name: "PENDING",
			},
			Items: []models.OrderItemModel{
				{ID: "item-1", OrderID: "order-1", ProductID: "product-1", Quantity: 2, UnitPrice: 50.0},
				{ID: "item-2", OrderID: "order-1", ProductID: "product-2", Quantity: 1, UnitPrice: 100.0},
			},
			CreatedAt: now,
		},
	}

	daos := FromModelArrayToDAOArray(models)

	assert.Len(t, daos, 1)
	assert.Len(t, daos[0].Items, 2)
	assert.Equal(t, "product-1", daos[0].Items[0].ProductID)
	assert.Equal(t, "product-2", daos[0].Items[1].ProductID)
}

// ============================================================================
// Tests for Round-trip conversions
// ============================================================================

func TestRoundTripConversion_DAOToModelToDAO(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	originalDAO := daos.OrderDAO{
		ID:         "order-1",
		CustomerID: &customerID,
		Amount:     100.0,
		Status: daos.OrderStatusDAO{
			ID:   "status-1",
			Name: "PENDING",
		},
		Items: []daos.OrderItemDAO{
			{
				ID:        "item-1",
				OrderID:   "order-1",
				ProductID: "product-1",
				Quantity:  2,
				UnitPrice: 50.0,
			},
		},
		CreatedAt: now,
	}

	model := FromDAOToModel(originalDAO)
	resultDAO := FromModelToDAO(model)

	assert.Equal(t, originalDAO.ID, resultDAO.ID)
	assert.Equal(t, *originalDAO.CustomerID, *resultDAO.CustomerID)
	assert.Equal(t, originalDAO.Amount, resultDAO.Amount)
	assert.Equal(t, originalDAO.Status.ID, resultDAO.Status.ID)
	assert.Equal(t, originalDAO.Status.Name, resultDAO.Status.Name)
	assert.Len(t, resultDAO.Items, len(originalDAO.Items))
}

func TestRoundTripConversion_ModelToDAOToModel(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	originalModel := models.OrderModel{
		ID:         "order-1",
		CustomerID: &customerID,
		Amount:     100.0,
		StatusID:   "status-1",
		Status: models.OrderStatusModel{
			ID:   "status-1",
			Name: "PENDING",
		},
		Items: []models.OrderItemModel{
			{
				ID:        "item-1",
				OrderID:   "order-1",
				ProductID: "product-1",
				Quantity:  2,
				UnitPrice: 50.0,
			},
		},
		CreatedAt: now,
	}

	dao := FromModelToDAO(originalModel)
	resultModel := FromDAOToModel(dao)

	assert.Equal(t, originalModel.ID, resultModel.ID)
	assert.Equal(t, *originalModel.CustomerID, *resultModel.CustomerID)
	assert.Equal(t, originalModel.Amount, resultModel.Amount)
	assert.Equal(t, originalModel.StatusID, resultModel.StatusID)
	assert.Len(t, resultModel.Items, len(originalModel.Items))
}

// ============================================================================
// Tests for edge cases
// ============================================================================

func TestFromDAOToModel_LargeAmount(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	dao := daos.OrderDAO{
		ID:         "order-1",
		CustomerID: &customerID,
		Amount:     999999999.99,
		Status: daos.OrderStatusDAO{
			ID:   "status-1",
			Name: "PENDING",
		},
		Items:     []daos.OrderItemDAO{},
		CreatedAt: now,
	}

	model := FromDAOToModel(dao)

	assert.Equal(t, 999999999.99, model.Amount)
}

func TestFromDAOToModel_SpecialCharactersInIDs(t *testing.T) {
	customerID := "customer-123-special-!@#"
	now := time.Now()

	dao := daos.OrderDAO{
		ID:         "order-special-!@#",
		CustomerID: &customerID,
		Amount:     100.0,
		Status: daos.OrderStatusDAO{
			ID:   "status-special-!@#",
			Name: "PENDING",
		},
		Items:     []daos.OrderItemDAO{},
		CreatedAt: now,
	}

	model := FromDAOToModel(dao)

	assert.Equal(t, "order-special-!@#", model.ID)
	assert.Equal(t, "customer-123-special-!@#", *model.CustomerID)
	assert.Equal(t, "status-special-!@#", model.Status.ID)
}

func TestFromModelArrayToDAOArray_LargeArray(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	modelArray := make([]models.OrderModel, 100)
	for i := 0; i < 100; i++ {
		modelArray[i] = models.OrderModel{
			ID:         "order-" + string(rune(i)),
			CustomerID: &customerID,
			Amount:     float64(i) * 10.0,
			StatusID:   "status-1",
			Status: models.OrderStatusModel{
				ID:   "status-1",
				Name: "PENDING",
			},
			Items:     []models.OrderItemModel{},
			CreatedAt: now,
		}
	}

	daos := FromModelArrayToDAOArray(modelArray)

	assert.Len(t, daos, 100)
}

func TestFromDAOToModel_ItemsPreservation(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	dao := daos.OrderDAO{
		ID:         "order-1",
		CustomerID: &customerID,
		Amount:     100.0,
		Status: daos.OrderStatusDAO{
			ID:   "status-1",
			Name: "PENDING",
		},
		Items: []daos.OrderItemDAO{
			{ID: "item-1", OrderID: "order-1", ProductID: "product-1", Quantity: 2, UnitPrice: 50.0},
			{ID: "item-2", OrderID: "order-1", ProductID: "product-2", Quantity: 1, UnitPrice: 100.0},
			{ID: "item-3", OrderID: "order-1", ProductID: "product-3", Quantity: 5, UnitPrice: 20.0},
		},
		CreatedAt: now,
	}

	model := FromDAOToModel(dao)

	for i, item := range model.Items {
		assert.Equal(t, dao.Items[i].ID, item.ID)
		assert.Equal(t, dao.Items[i].OrderID, item.OrderID)
		assert.Equal(t, dao.Items[i].ProductID, item.ProductID)
		assert.Equal(t, dao.Items[i].Quantity, item.Quantity)
		assert.Equal(t, dao.Items[i].UnitPrice, item.UnitPrice)
	}
}

func TestFromModelToDAO_ItemsPreservation(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	model := models.OrderModel{
		ID:         "order-1",
		CustomerID: &customerID,
		Amount:     100.0,
		StatusID:   "status-1",
		Status: models.OrderStatusModel{
			ID:   "status-1",
			Name: "PENDING",
		},
		Items: []models.OrderItemModel{
			{ID: "item-1", OrderID: "order-1", ProductID: "product-1", Quantity: 2, UnitPrice: 50.0},
			{ID: "item-2", OrderID: "order-1", ProductID: "product-2", Quantity: 1, UnitPrice: 100.0},
			{ID: "item-3", OrderID: "order-1", ProductID: "product-3", Quantity: 5, UnitPrice: 20.0},
		},
		CreatedAt: now,
	}

	dao := FromModelToDAO(model)

	for i, item := range dao.Items {
		assert.Equal(t, model.Items[i].ID, item.ID)
		assert.Equal(t, model.Items[i].OrderID, item.OrderID)
		assert.Equal(t, model.Items[i].ProductID, item.ProductID)
		assert.Equal(t, model.Items[i].Quantity, item.Quantity)
		assert.Equal(t, model.Items[i].UnitPrice, item.UnitPrice)
	}
}

func TestFromDAOToModel_StatusPreservation(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	statuses := []daos.OrderStatusDAO{
		{ID: "status-1", Name: "PENDING"},
		{ID: "status-2", Name: "CONFIRMED"},
		{ID: "status-3", Name: "CANCELLED"},
	}

	for _, status := range statuses {
		dao := daos.OrderDAO{
			ID:         "order-1",
			CustomerID: &customerID,
			Amount:     100.0,
			Status:     status,
			Items:      []daos.OrderItemDAO{},
			CreatedAt:  now,
		}

		model := FromDAOToModel(dao)

		assert.Equal(t, status.ID, model.Status.ID)
		assert.Equal(t, status.Name, model.Status.Name)
		assert.Equal(t, status.ID, model.StatusID)
	}
}

func TestFromModelToDAO_StatusPreservation(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	statuses := []models.OrderStatusModel{
		{ID: "status-1", Name: "PENDING"},
		{ID: "status-2", Name: "CONFIRMED"},
		{ID: "status-3", Name: "CANCELLED"},
	}

	for _, status := range statuses {
		model := models.OrderModel{
			ID:         "order-1",
			CustomerID: &customerID,
			Amount:     100.0,
			StatusID:   status.ID,
			Status:     status,
			Items:      []models.OrderItemModel{},
			CreatedAt:  now,
		}

		dao := FromModelToDAO(model)

		assert.Equal(t, status.ID, dao.Status.ID)
		assert.Equal(t, status.Name, dao.Status.Name)
	}
}

func TestFromDAOToModel_TimestampPreservation(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()
	updatedAt := now.Add(24 * time.Hour)

	dao := daos.OrderDAO{
		ID:         "order-1",
		CustomerID: &customerID,
		Amount:     100.0,
		Status: daos.OrderStatusDAO{
			ID:   "status-1",
			Name: "PENDING",
		},
		Items:     []daos.OrderItemDAO{},
		CreatedAt: now,
		UpdatedAt: &updatedAt,
	}

	model := FromDAOToModel(dao)

	assert.Equal(t, now, model.CreatedAt)
	assert.Equal(t, updatedAt, *model.UpdatedAt)
}

func TestFromModelToDAO_TimestampPreservation(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()
	updatedAt := now.Add(24 * time.Hour)

	model := models.OrderModel{
		ID:         "order-1",
		CustomerID: &customerID,
		Amount:     100.0,
		StatusID:   "status-1",
		Status: models.OrderStatusModel{
			ID:   "status-1",
			Name: "PENDING",
		},
		Items:     []models.OrderItemModel{},
		CreatedAt: now,
		UpdatedAt: &updatedAt,
	}

	dao := FromModelToDAO(model)

	assert.Equal(t, now, dao.CreatedAt)
	assert.Equal(t, updatedAt, *dao.UpdatedAt)
}
