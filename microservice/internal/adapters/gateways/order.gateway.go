package gateways

import (
	"microservice/internal/adapters/daos"
	"microservice/internal/adapters/dtos"
	"microservice/internal/domain/entities"
	"microservice/internal/interfaces"
)

type OrderGateway struct {
	datasource interfaces.IOrderDataSource
}

func NewOrderGateway(datasource interfaces.IOrderDataSource) *OrderGateway {
	return &OrderGateway{datasource: datasource}
}

func (g *OrderGateway) Create(order entities.Order) error {
	items := make([]daos.OrderItemDAO, len(order.Items))
	for i, item := range order.Items {
		items[i] = daos.OrderItemDAO{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID.Value(),
			Quantity:  item.Quantity.Value(),
			UnitPrice: item.UnitPrice.Value(),
		}
	}

	return g.datasource.Create(daos.OrderDAO{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		Amount:     order.Amount.Value(),
		Status: daos.OrderStatusDAO{
			ID:   order.Status.ID,
			Name: order.Status.Name.Value(),
		},
		Items:     items,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	})
}

func (g *OrderGateway) FindByID(id string) (*entities.Order, error) {
	orderDAO, err := g.datasource.FindByID(id)
	if err != nil {
		return nil, err
	}

	status, err := entities.NewOrderStatus(orderDAO.Status.ID, orderDAO.Status.Name)
	if err != nil {
		return nil, err
	}

	items := make([]entities.OrderItem, len(orderDAO.Items))
	for i, itemDAO := range orderDAO.Items {
		item, err := entities.NewOrderItem(
			itemDAO.ID,
			itemDAO.ProductID,
			itemDAO.OrderID,
			itemDAO.Quantity,
			itemDAO.UnitPrice,
		)
		if err != nil {
			return nil, err
		}
		items[i] = *item
	}

	order, err := entities.NewOrderWithItems(
		orderDAO.ID,
		orderDAO.CustomerID,
		orderDAO.Amount,
		*status,
		items,
		orderDAO.CreatedAt,
		orderDAO.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (g *OrderGateway) FindAll(filter dtos.OrderFilterDTO) ([]entities.Order, error) {
	orderDAOs, err := g.datasource.FindAll(filter)
	if err != nil {
		return nil, err
	}

	orders := make([]entities.Order, 0, len(orderDAOs))
	for _, orderDAO := range orderDAOs {
		status, err := entities.NewOrderStatus(orderDAO.Status.ID, orderDAO.Status.Name)
		if err != nil {
			return nil, err
		}

		items := make([]entities.OrderItem, len(orderDAO.Items))
		for i, itemDAO := range orderDAO.Items {
			item, err := entities.NewOrderItem(
				itemDAO.ID,
				itemDAO.ProductID,
				itemDAO.OrderID,
				itemDAO.Quantity,
				itemDAO.UnitPrice,
			)
			if err != nil {
				return nil, err
			}
			items[i] = *item
		}

		order, err := entities.NewOrderWithItems(
			orderDAO.ID,
			orderDAO.CustomerID,
			orderDAO.Amount,
			*status,
			items,
			orderDAO.CreatedAt,
			orderDAO.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, *order)
	}

	return orders, nil
}

func (g *OrderGateway) Update(order entities.Order) error {
	items := make([]daos.OrderItemDAO, len(order.Items))
	for i, item := range order.Items {
		items[i] = daos.OrderItemDAO{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID.Value(),
			Quantity:  item.Quantity.Value(),
			UnitPrice: item.UnitPrice.Value(),
		}
	}

	return g.datasource.Update(daos.OrderDAO{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		Amount:     order.Amount.Value(),
		Status: daos.OrderStatusDAO{
			ID:   order.Status.ID,
			Name: order.Status.Name.Value(),
		},
		Items:     items,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	})
}

func (g *OrderGateway) Delete(id string) error {
	return g.datasource.Delete(id)
}
