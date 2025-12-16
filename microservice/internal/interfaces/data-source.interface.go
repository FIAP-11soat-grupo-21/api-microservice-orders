package interfaces

import (
	"microservice/internal/adapters/daos"
	"microservice/internal/adapters/dtos"
)

type IOrderDataSource interface {
	Create(order daos.OrderDAO) error
	FindAll(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error)
	FindByID(id string) (daos.OrderDAO, error)
	Update(order daos.OrderDAO) error
	Delete(id string) error
}

type IOrderStatusDataSource interface {
	FindByID(id string) (daos.OrderStatusDAO, error)
	FindAll() ([]daos.OrderStatusDAO, error)
}
