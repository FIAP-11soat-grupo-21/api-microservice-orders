package factories

import (
	"microservice/infra/db/postgres/data_source"
	"microservice/internal/interfaces"
)

var newOrderDataSource func() interfaces.IOrderDataSource = func() interfaces.IOrderDataSource {
	return data_source.NewGormOrderDataSource()
}

var newOrderStatusDataSource func() interfaces.IOrderStatusDataSource = func() interfaces.IOrderStatusDataSource {
	return data_source.NewGormOrderStatusDataSource()
}

func NewOrderDataSource() interfaces.IOrderDataSource {
	return newOrderDataSource()
}

func NewOrderStatusDataSource() interfaces.IOrderStatusDataSource {
	return newOrderStatusDataSource()
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
