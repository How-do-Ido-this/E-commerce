package inventory

import (
	"context"
	"e-commerce/pkg/events"
	"e-commerce/pkg/rabbitmq"
	"encoding/json"
	"github.com/google/uuid"
)

// Service implementa la lógica de negocio para el inventario.
type Service struct {
	repo         Repository
	rabbitClient *rabbitmq.Client
}

// NewService crea una nueva instancia de Service.
func NewService(r Repository, rc *rabbitmq.Client) *Service {
	return &Service{repo: r, rabbitClient: rc}
}

// ReserveStock valida y reserva stock.
func (s *Service) ReserveStock(ctx context.Context, event events.OrderCreated) error {
	// 1. Obtener inventario
	productID := event.Items[0].ProductID.String()
	inv, err := s.repo.GetByProductID(ctx, productID)
	if err != nil {
		return err
	}

	// 2. Validar
	if inv.QuantityAvailable < event.Items[0].Quantity {
		// Publicar evento de error
		errorEvent := events.StockInsufficient{
			OrderID:   event.OrderID,
			ProductID: event.Items[0].ProductID,
			Reason:    "insufficient stock",
		}
		body, _ := json.Marshal(errorEvent)
		return s.rabbitClient.Publish("orders", "stock.insufficient", body)
	}

	// 3. Reservar
	inv.QuantityAvailable -= event.Items[0].Quantity
	inv.ReservedQuantity += event.Items[0].Quantity
	if err := s.repo.Update(ctx, inv); err != nil {
		return err
	}

	// 4. Publicar evento de éxito
	successEvent := events.InventoryReserved{
		OrderID:                event.OrderID,
		InventoryReservationID: uuid.New(),
	}
	body, _ := json.Marshal(successEvent)
	return s.rabbitClient.Publish("orders", "inventory.reserved", body)
}
