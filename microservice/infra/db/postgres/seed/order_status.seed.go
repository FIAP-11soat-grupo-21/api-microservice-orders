package seed

import (
	"log"

	"gorm.io/gorm"

	"microservice/infra/db/postgres/models"
)

const (
	ORDER_STATUS_RECEIVED_ID  = "56d3b3c3-1801-49cd-bae7-972c78082012"
	ORDER_STATUS_CONFIRMED_ID = "3f9a1c98-7b2f-4f3b-8a96-c0b7c761a123"
	ORDER_STATUS_PREPARING_ID = "5a8b2b16-9b47-4e35-ae27-28f7994ef456"
	ORDER_STATUS_READY_ID     = "bd91a1ee-1234-4cde-9c2a-efb1d2a3a789"
	ORDER_STATUS_DELIVERED_ID = "f1e2d3c4-5b6a-7c8d-9e0f-1a2b3c4d5e6f"
)

func SeedOrderStatus(db *gorm.DB) {
	defaults := []models.OrderStatusModel{
		{ID: ORDER_STATUS_RECEIVED_ID, Name: "Recebido"},
		{ID: ORDER_STATUS_CONFIRMED_ID, Name: "Confirmado"},
		{ID: ORDER_STATUS_PREPARING_ID, Name: "Em preparação"},
		{ID: ORDER_STATUS_READY_ID, Name: "Pronto"},
		{ID: ORDER_STATUS_DELIVERED_ID, Name: "Entregue"},
	}

	for _, status := range defaults {
		var existing models.OrderStatusModel
		if err := db.Where("id = ?", status.ID).First(&existing).Error; err == gorm.ErrRecordNotFound {
			statusCopy := status
			if err := db.Create(&statusCopy).Error; err != nil {
				log.Printf("Erro ao criar status %s: %v", status.ID, err)
			} else {
				log.Printf("Status %s criado com sucesso", status.Name)
			}
		} else {
			log.Printf("Status %s já existe", status.Name)
		}
	}
}
