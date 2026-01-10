package consumers

import "microservice/internal/use_cases"

type IProcessPaymentConfirmationUseCase interface {
	Execute(dto use_cases.PaymentConfirmationDTO) (*use_cases.PaymentConfirmationResult, error)
}

type IUpdateOrderStatusUseCase interface {
	Execute(dto use_cases.UpdateOrderStatusDTO) (*use_cases.UpdateOrderStatusResult, error)
}
