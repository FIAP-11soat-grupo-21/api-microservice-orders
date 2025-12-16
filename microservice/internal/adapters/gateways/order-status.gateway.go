package gateways

import (
	"microservice/internal/domain/entities"
	"microservice/internal/interfaces"
)

type OrderStatusGateway struct {
	datasource interfaces.IOrderStatusDataSource
}

func NewOrderStatusGateway(datasource interfaces.IOrderStatusDataSource) *OrderStatusGateway {
	return &OrderStatusGateway{datasource: datasource}
}

func (g *OrderStatusGateway) FindAll() ([]entities.OrderStatus, error) {
	statusDAOs, err := g.datasource.FindAll()
	if err != nil {
		return nil, err
	}

	statuses := make([]entities.OrderStatus, 0, len(statusDAOs))
	for _, statusDAO := range statusDAOs {
		status, err := entities.NewOrderStatus(statusDAO.ID, statusDAO.Name)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, *status)
	}

	return statuses, nil
}

func (g *OrderStatusGateway) FindByID(id string) (*entities.OrderStatus, error) {
	statusDAO, err := g.datasource.FindByID(id)
	if err != nil {
		return nil, err
	}

	status, err := entities.NewOrderStatus(statusDAO.ID, statusDAO.Name)
	if err != nil {
		return nil, err
	}

	return status, nil
}
