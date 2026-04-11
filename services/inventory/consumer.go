package inventory

import (
	"context"
	"e-commerce/pkg/events"
	"e-commerce/pkg/rabbitmq"
	"encoding/json"
	"log"
)

// Consumer se encarga de escuchar eventos de RabbitMQ.
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
	deliveries, err := c.client.Consume("order_created_queue")
	if err != nil {
		log.Fatalf("failed to register consumer: %v", err)
	}

	for d := range deliveries {
		var event events.OrderCreated
		if err := json.Unmarshal(d.Body, &event); err != nil {
			log.Printf("failed to unmarshal event: %v", err)
			continue
		}

		log.Printf("Received OrderCreated event: %s", event.OrderID)

		// Llamar a la lógica de negocio
		if err := c.service.ReserveStock(context.Background(), event); err != nil {
			log.Printf("failed to reserve stock: %v", err)
		}
	}
}
