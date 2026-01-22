package factories

import (
	"microservice/infra/api/client"
	"microservice/infra/db/postgres/data_source"
	"microservice/infra/messaging"
	"microservice/internal/adapters/brokers"
	"microservice/internal/adapters/gateways"
	"microservice/internal/interfaces"
	"microservice/internal/use_cases"
	"microservice/utils/config"
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

var newApiClient func() client.IApiClient = func() client.IApiClient {
	cfg := config.LoadConfig()

	return client.NewApiClient(cfg.APIGatewayURL)
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

func NewApiClient() client.IApiClient {
	return newApiClient()
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

func SetNewApiClient(fn func() client.IApiClient) {
	if fn == nil {
		newApiClient = func() client.IApiClient {
			cfg := config.LoadConfig()
			return client.NewApiClient(cfg.APIGatewayURL)
		}
		return
	}
	newApiClient = fn
}
