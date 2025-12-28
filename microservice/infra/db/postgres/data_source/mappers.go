package data_source

import (
	"microservice/infra/db/postgres/models"
	"microservice/internal/adapters/daos"
)

func FromDAOToModel(order daos.OrderDAO) models.OrderModel {
	items := make([]models.OrderItemModel, len(order.Items))
	for i, item := range order.Items {
		items[i] = models.OrderItemModel{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		}
	}

	return models.OrderModel{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		Amount:     order.Amount,
		StatusID:   order.Status.ID,
		Status: models.OrderStatusModel{
			ID:   order.Status.ID,
			Name: order.Status.Name,
		},
		Items:     items,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}
}

func FromModelToDAO(order models.OrderModel) daos.OrderDAO {
	items := make([]daos.OrderItemDAO, len(order.Items))
	for i, item := range order.Items {
		items[i] = daos.OrderItemDAO{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		}
	}

	return daos.OrderDAO{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		Amount:     order.Amount,
		Status: daos.OrderStatusDAO{
			ID:   order.Status.ID,
			Name: order.Status.Name,
		},
		Items:     items,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}
}

func FromModelArrayToDAOArray(orders []models.OrderModel) []daos.OrderDAO {
	result := make([]daos.OrderDAO, len(orders))
	for i, order := range orders {
		result[i] = FromModelToDAO(order)
	}
	return result
}
