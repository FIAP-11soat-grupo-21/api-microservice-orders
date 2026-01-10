package handlers

import (
	"microservice/internal/interfaces"
	"microservice/internal/test_helpers"
	"microservice/utils/factories"
)

func SetupTestMocks(orderDS *test_helpers.MockOrderDataSource, statusDS *test_helpers.MockOrderStatusDataSource) func() {
	factories.SetNewOrderDataSource(func() interfaces.IOrderDataSource {
		return orderDS
	})
	factories.SetNewOrderStatusDataSource(func() interfaces.IOrderStatusDataSource {
		return statusDS
	})

	return func() {
		factories.SetNewOrderDataSource(nil)
		factories.SetNewOrderStatusDataSource(nil)
	}
}
