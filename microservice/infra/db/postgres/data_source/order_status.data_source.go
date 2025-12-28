package data_source

import (
	"gorm.io/gorm"

	"microservice/infra/db/postgres"
	"microservice/infra/db/postgres/models"
	"microservice/internal/adapters/daos"
)

type GormOrderStatusDataSource struct {
	db *gorm.DB
}

func NewGormOrderStatusDataSource() *GormOrderStatusDataSource {
	return &GormOrderStatusDataSource{
		db: postgres.GetDB(),
	}
}

func (r *GormOrderStatusDataSource) FindAll() ([]daos.OrderStatusDAO, error) {
	var statuses []models.OrderStatusModel

	if err := r.db.Find(&statuses).Error; err != nil {
		return nil, err
	}

	result := make([]daos.OrderStatusDAO, len(statuses))
	for i, status := range statuses {
		result[i] = daos.OrderStatusDAO{
			ID:   status.ID,
			Name: status.Name,
		}
	}

	return result, nil
}

func (r *GormOrderStatusDataSource) FindByID(id string) (daos.OrderStatusDAO, error) {
	var status models.OrderStatusModel

	if err := r.db.First(&status, "id = ?", id).Error; err != nil {
		return daos.OrderStatusDAO{}, err
	}

	return daos.OrderStatusDAO{
		ID:   status.ID,
		Name: status.Name,
	}, nil
}
