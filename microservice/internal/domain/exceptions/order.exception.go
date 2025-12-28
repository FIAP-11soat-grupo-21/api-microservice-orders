package exceptions

type OrderNotFoundException struct {
	Message string
}

type InvalidOrderDataException struct {
	Message string
}

type AmountNotValidException struct {
	Message string
}

type OrderStatusNotFoundException struct {
	Message string
}

func (e *OrderNotFoundException) Error() string {
	if e.Message == "" {
		return "Order not found"
	}
	return e.Message
}

func (e *InvalidOrderDataException) Error() string {
	if e.Message == "" {
		return "Invalid order data"
	}
	return e.Message
}

func (e *AmountNotValidException) Error() string {
	if e.Message == "" {
		return "Amount is not valid"
	}
	return e.Message
}

func (e *OrderStatusNotFoundException) Error() string {
	if e.Message == "" {
		return "Order Status not found"
	}
	return e.Message
}
