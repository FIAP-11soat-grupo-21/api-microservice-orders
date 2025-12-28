package models

import (
	"testing"
)

func TestOrderModel_TableName(t *testing.T) {
	model := OrderModel{}
	tableName := model.TableName()

	if tableName != "orders" {
		t.Errorf("OrderModel.TableName() = %v, want orders", tableName)
	}
}

func TestOrderItemModel_TableName(t *testing.T) {
	model := OrderItemModel{}
	tableName := model.TableName()

	if tableName != "order_items" {
		t.Errorf("OrderItemModel.TableName() = %v, want order_items", tableName)
	}
}

func TestOrderStatusModel_TableName(t *testing.T) {
	model := OrderStatusModel{}
	tableName := model.TableName()

	if tableName != "order_status" {
		t.Errorf("OrderStatusModel.TableName() = %v, want order_status", tableName)
	}
}
