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
