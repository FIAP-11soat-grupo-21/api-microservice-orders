package interfaces

import (
	"microservice/internal/adapters/dtos"
	"microservice/internal/domain/entities"
)

type IOrderGateway interface {
	Create(order entities.Order) error
	FindByID(id string) (*entities.Order, error)
	FindAll(filter dtos.OrderFilterDTO) ([]entities.Order, error)
	Update(order entities.Order) error
	Delete(id string) error
}

type IOrderStatusGateway interface {
	FindAll() ([]entities.OrderStatus, error)
	FindByID(id string) (*entities.OrderStatus, error)
}
