package data_source

import (
	"gorm.io/gorm"

	"microservice/infra/db/postgres"
	"microservice/infra/db/postgres/models"
	"microservice/internal/adapters/daos"
	"microservice/internal/adapters/dtos"
)

type GormOrderDataSource struct {
	db *gorm.DB
}

func NewGormOrderDataSource() *GormOrderDataSource {
	return &GormOrderDataSource{
		db: postgres.GetDB(),
	}
}

func (r *GormOrderDataSource) Create(order daos.OrderDAO) error {
	orderModel := FromDAOToModel(order)
	return r.db.Create(&orderModel).Error
}

func (r *GormOrderDataSource) FindAll(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
	var orders []models.OrderModel

	query := r.db.
		Preload("Status").
		Preload("Items").
		Order("orders.created_at DESC")

	if filter.CreatedAtFrom != nil {
		query = query.Where("orders.created_at >= ?", *filter.CreatedAtFrom)
	}
	if filter.CreatedAtTo != nil {
		query = query.Where("orders.created_at <= ?", *filter.CreatedAtTo)
	}
	if filter.StatusID != nil {
		query = query.Where("orders.status_id = ?", *filter.StatusID)
	}
	if filter.CustomerID != nil {
		query = query.Where("orders.customer_id = ?", *filter.CustomerID)
	}

	if err := query.Find(&orders).Error; err != nil {
		return nil, err
	}

	return FromModelArrayToDAOArray(orders), nil
}

func (r *GormOrderDataSource) FindByID(id string) (daos.OrderDAO, error) {
	var order models.OrderModel

	if err := r.db.Preload("Status").Preload("Items").First(&order, "id = ?", id).Error; err != nil {
		return daos.OrderDAO{}, err
	}

	return FromModelToDAO(order), nil
}

func (r *GormOrderDataSource) Update(order daos.OrderDAO) error {
	orderModel := FromDAOToModel(order)
	return r.db.Save(&orderModel).Error
}

func (r *GormOrderDataSource) Delete(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.OrderItemModel{}, "order_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&models.OrderModel{}, "id = ?", id).Error
	})
}
