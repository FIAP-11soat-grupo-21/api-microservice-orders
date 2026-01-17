package factories

import (
	"microservice/infra/db/postgres/data_source"
	"microservice/infra/messaging"
	"microservice/internal/adapters/brokers"
	"microservice/internal/adapters/gateways"
	"microservice/internal/interfaces"
	"microservice/internal/use_cases"
)

var newOrderDataSource func() interfaces.IOrderDataSource = func() interfaces.IOrderDataSource {
	return data_source.NewGormOrderDataSource()
}

var newOrderStatusDataSource func() interfaces.IOrderStatusDataSource = func() interfaces.IOrderStatusDataSource {
	return data_source.NewGormOrderStatusDataSource()
}

var newMessageBroker func() brokers.MessageBroker = func() brokers.MessageBroker {
	return messaging.GetBroker()
}

func NewOrderDataSource() interfaces.IOrderDataSource {
	return newOrderDataSource()
}

func NewOrderStatusDataSource() interfaces.IOrderStatusDataSource {
	return newOrderStatusDataSource()
}

func NewMessageBroker() brokers.MessageBroker {
	return newMessageBroker()
}

func NewUpdateOrderStatusUseCase() *use_cases.UpdateOrderStatusUseCase {
	orderDataSource := NewOrderDataSource()
	orderStatusDataSource := NewOrderStatusDataSource()
	orderGateway := gateways.NewOrderGateway(orderDataSource)
	orderStatusGateway := gateways.NewOrderStatusGateway(orderStatusDataSource)
	return use_cases.NewUpdateOrderStatusUseCase(orderGateway, orderStatusGateway)
}

func SetNewOrderDataSource(fn func() interfaces.IOrderDataSource) {
	if fn == nil {
		newOrderDataSource = func() interfaces.IOrderDataSource {
			return data_source.NewGormOrderDataSource()
		}
		return
	}
	newOrderDataSource = fn
}

func SetNewOrderStatusDataSource(fn func() interfaces.IOrderStatusDataSource) {
	if fn == nil {
		newOrderStatusDataSource = func() interfaces.IOrderStatusDataSource {
			return data_source.NewGormOrderStatusDataSource()
		}
		return
	}
	newOrderStatusDataSource = fn
}

func SetNewMessageBroker(fn func() brokers.MessageBroker) {
	if fn == nil {
		newMessageBroker = func() brokers.MessageBroker {
			return messaging.GetBroker()
		}
		return
	}
	newMessageBroker = fn
}
