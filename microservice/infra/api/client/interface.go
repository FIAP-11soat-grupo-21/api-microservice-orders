package client

type ProductResponseDTO struct {
	ID     string  `json:"id"`
	Price  float64 `json:"price"`
	Active bool    `json:"active"`
}

type IApiClient interface {
	Get(path string, obj any) error
}
