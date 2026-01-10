package use_cases

import (
	"microservice/internal/adapters/dtos"
	"microservice/internal/domain/entities"
	"microservice/internal/interfaces"
)

type FindAllOrdersUseCase struct {
	orderGateway interfaces.IOrderGateway
}

func NewFindAllOrdersUseCase(orderGateway interfaces.IOrderGateway) *FindAllOrdersUseCase {
	return &FindAllOrdersUseCase{
		orderGateway: orderGateway,
	}
}

func (uc *FindAllOrdersUseCase) Execute(filter dtos.OrderFilterDTO) ([]entities.Order, error) {
	return uc.orderGateway.FindAll(filter)
}
