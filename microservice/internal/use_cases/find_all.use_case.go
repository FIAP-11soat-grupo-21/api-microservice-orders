package use_cases

import (
	"microservice/internal/adapters/dtos"
	"microservice/internal/adapters/gateways"
	"microservice/internal/domain/entities"
)

type FindAllOrdersUseCase struct {
	orderGateway gateways.OrderGateway
}

func NewFindAllOrdersUseCase(orderGateway gateways.OrderGateway) *FindAllOrdersUseCase {
	return &FindAllOrdersUseCase{
		orderGateway: orderGateway,
	}
}

func (uc *FindAllOrdersUseCase) Execute(filter dtos.OrderFilterDTO) ([]entities.Order, error) {
	return uc.orderGateway.FindAll(filter)
}
