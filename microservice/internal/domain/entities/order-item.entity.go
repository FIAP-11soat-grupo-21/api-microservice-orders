package entities

import "microservice/internal/domain/value_objects"

type OrderItem struct {
	ID        string
	OrderID   string
	ProductID value_objects.ProductID
	Quantity  value_objects.Quantity
	UnitPrice value_objects.UnitPrice
}

func NewOrderItem(id string, productID string, orderID string, quantity int, unitPrice float64) (*OrderItem, error) {
	productIDValueObject, err := value_objects.NewProductID(productID)
	if err != nil {
		return nil, err
	}

	quantityValueObject, err := value_objects.NewQuantity(quantity)
	if err != nil {
		return nil, err
	}

	unitPriceValueObject, err := value_objects.NewUnitPrice(unitPrice)
	if err != nil {
		return nil, err
	}

	return &OrderItem{
		ID:        id,
		OrderID:   orderID,
		ProductID: productIDValueObject,
		Quantity:  quantityValueObject,
		UnitPrice: unitPriceValueObject,
	}, nil
}

func (oi *OrderItem) GetTotal() float64 {
	return oi.UnitPrice.Value() * float64(oi.Quantity.Value())
}
