package factories

import (
	"microservice/internal/adapters/gateways"
	"microservice/internal/use_cases"
)

func NewProcessPaymentConfirmationUseCase() *use_cases.ProcessPaymentConfirmationUseCase {
	orderDataSource := NewOrderDataSource()
	orderStatusDataSource := NewOrderStatusDataSource()
	
	orderGateway := gateways.NewOrderGateway(orderDataSource)
	orderStatusGateway := gateways.NewOrderStatusGateway(orderStatusDataSource)
	
	return use_cases.NewProcessPaymentConfirmationUseCase(orderGateway, orderStatusGateway)
}