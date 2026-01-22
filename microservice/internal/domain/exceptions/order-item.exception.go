package exceptions

type InvalidOrderItemData struct {
	Message string
}

func (e *InvalidOrderItemData) Error() string {
	if e.Message == "" {
		return "Invalid order item data"
	}
	return e.Message
}

type OrderItemProductInactiveException struct {
	ProductID string
}

func (e *OrderItemProductInactiveException) Error() string {
	return "The product is inactive: " + e.ProductID
}
