package data_source

import (
	"testing"
	"time"

	"microservice/infra/db/postgres/models"
	"microservice/internal/adapters/daos"
)

func TestFromDAOToModel(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()
	updatedAt := now.Add(time.Hour)

	dao := daos.OrderDAO{
		ID:         "order-1",
		CustomerID: &customerID,
		Amount:     100.0,
		Status: daos.OrderStatusDAO{
			ID:   "status-1",
			Name: "Pending",
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
		UpdatedAt: &updatedAt,
	}

	model := FromDAOToModel(dao)

	if model.ID != "order-1" {
		t.Errorf("FromDAOToModel() ID = %v, want order-1", model.ID)
	}
	if *model.CustomerID != customerID {
		t.Errorf("FromDAOToModel() CustomerID = %v, want %v", *model.CustomerID, customerID)
	}
	if model.Amount != 100.0 {
		t.Errorf("FromDAOToModel() Amount = %v, want 100.0", model.Amount)
	}
	if model.StatusID != "status-1" {
		t.Errorf("FromDAOToModel() StatusID = %v, want status-1", model.StatusID)
	}
	if model.Status.ID != "status-1" {
		t.Errorf("FromDAOToModel() Status.ID = %v, want status-1", model.Status.ID)
	}
	if model.Status.Name != "Pending" {
		t.Errorf("FromDAOToModel() Status.Name = %v, want Pending", model.Status.Name)
	}
	if len(model.Items) != 1 {
		t.Errorf("FromDAOToModel() Items length = %v, want 1", len(model.Items))
	}
	if model.Items[0].ID != "item-1" {
		t.Errorf("FromDAOToModel() Items[0].ID = %v, want item-1", model.Items[0].ID)
	}
	if model.Items[0].ProductID != "product-1" {
		t.Errorf("FromDAOToModel() Items[0].ProductID = %v, want product-1", model.Items[0].ProductID)
	}
	if model.Items[0].Quantity != 2 {
		t.Errorf("FromDAOToModel() Items[0].Quantity = %v, want 2", model.Items[0].Quantity)
	}
	if model.Items[0].UnitPrice != 50.0 {
		t.Errorf("FromDAOToModel() Items[0].UnitPrice = %v, want 50.0", model.Items[0].UnitPrice)
	}
}

func TestFromDAOToModel_NilCustomerID(t *testing.T) {
	now := time.Now()

	dao := daos.OrderDAO{
		ID:         "order-1",
		CustomerID: nil,
		Amount:     100.0,
		Status: daos.OrderStatusDAO{
			ID:   "status-1",
			Name: "Pending",
		},
		Items:     []daos.OrderItemDAO{},
		CreatedAt: now,
		UpdatedAt: nil,
	}

	model := FromDAOToModel(dao)

	if model.CustomerID != nil {
		t.Errorf("FromDAOToModel() CustomerID = %v, want nil", model.CustomerID)
	}
	if model.UpdatedAt != nil {
		t.Errorf("FromDAOToModel() UpdatedAt = %v, want nil", model.UpdatedAt)
	}
}

func TestFromModelToDAO(t *testing.T) {
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
			Name: "Pending",
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
		UpdatedAt: &updatedAt,
	}

	dao := FromModelToDAO(model)

	if dao.ID != "order-1" {
		t.Errorf("FromModelToDAO() ID = %v, want order-1", dao.ID)
	}
	if *dao.CustomerID != customerID {
		t.Errorf("FromModelToDAO() CustomerID = %v, want %v", *dao.CustomerID, customerID)
	}
	if dao.Amount != 100.0 {
		t.Errorf("FromModelToDAO() Amount = %v, want 100.0", dao.Amount)
	}
	if dao.Status.ID != "status-1" {
		t.Errorf("FromModelToDAO() Status.ID = %v, want status-1", dao.Status.ID)
	}
	if dao.Status.Name != "Pending" {
		t.Errorf("FromModelToDAO() Status.Name = %v, want Pending", dao.Status.Name)
	}
	if len(dao.Items) != 1 {
		t.Errorf("FromModelToDAO() Items length = %v, want 1", len(dao.Items))
	}
	if dao.Items[0].ID != "item-1" {
		t.Errorf("FromModelToDAO() Items[0].ID = %v, want item-1", dao.Items[0].ID)
	}
}

func TestFromModelToDAO_NilCustomerID(t *testing.T) {
	now := time.Now()

	model := models.OrderModel{
		ID:         "order-1",
		CustomerID: nil,
		Amount:     100.0,
		StatusID:   "status-1",
		Status: models.OrderStatusModel{
			ID:   "status-1",
			Name: "Pending",
		},
		Items:     []models.OrderItemModel{},
		CreatedAt: now,
		UpdatedAt: nil,
	}

	dao := FromModelToDAO(model)

	if dao.CustomerID != nil {
		t.Errorf("FromModelToDAO() CustomerID = %v, want nil", dao.CustomerID)
	}
	if dao.UpdatedAt != nil {
		t.Errorf("FromModelToDAO() UpdatedAt = %v, want nil", dao.UpdatedAt)
	}
}

func TestFromModelArrayToDAOArray(t *testing.T) {
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
				Name: "Pending",
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
				Name: "Confirmed",
			},
			Items:     []models.OrderItemModel{},
			CreatedAt: now,
		},
	}

	daos := FromModelArrayToDAOArray(models)

	if len(daos) != 2 {
		t.Errorf("FromModelArrayToDAOArray() length = %v, want 2", len(daos))
	}
	if daos[0].ID != "order-1" {
		t.Errorf("FromModelArrayToDAOArray()[0].ID = %v, want order-1", daos[0].ID)
	}
	if daos[1].ID != "order-2" {
		t.Errorf("FromModelArrayToDAOArray()[1].ID = %v, want order-2", daos[1].ID)
	}
}

func TestFromModelArrayToDAOArray_Empty(t *testing.T) {
	models := []models.OrderModel{}
	daos := FromModelArrayToDAOArray(models)

	if len(daos) != 0 {
		t.Errorf("FromModelArrayToDAOArray() length = %v, want 0", len(daos))
	}
}
