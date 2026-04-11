package order

import (
	"context"
	"e-commerce/pkg/events"
	"e-commerce/pkg/rabbitmq"
	"encoding/json"
	"log"
)

// Consumer se encarga de escuchar eventos de RabbitMQ para el servicio de órdenes.
type Consumer struct {
	client  *rabbitmq.Client
	service *Service
}

// NewConsumer crea una nueva instancia de Consumer.
func NewConsumer(rc *rabbitmq.Client, s *Service) *Consumer {
	return &Consumer{client: rc, service: s}
}

// StartListening inicia el proceso de consumo de eventos.
func (c *Consumer) StartListening() {
	deliveries, err := c.client.Consume("inventory_events_queue") // Cola para escuchar eventos de inventario
	if err != nil {
		log.Fatalf("failed to register consumer: %v", err)
	}

	for d := range deliveries {
		// Aquí escuchamos dos tipos de eventos, podríamos usar una cabecera para diferenciar
		// pero para simplicidad supongamos que procesamos ambos.
		var reserved events.InventoryReserved
		if err := json.Unmarshal(d.Body, &reserved); err == nil {
			log.Printf("Received InventoryReserved event: %s", reserved.OrderID)
			c.service.UpdateOrderStatus(context.Background(), reserved.OrderID.String(), StatusCreated)
			continue
		}

		var insufficient events.StockInsufficient
		if err := json.Unmarshal(d.Body, &insufficient); err == nil {
			log.Printf("Received StockInsufficient event: %s", insufficient.OrderID)
			c.service.UpdateOrderStatus(context.Background(), insufficient.OrderID.String(), StatusFailed)
		}
	}
}
