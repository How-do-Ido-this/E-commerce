package events

import "github.com/google/uuid"

type OrderItem struct {
	ProductID uuid.UUID `json:"productId"`
	Quantity  int       `json:"quantity"`
}

type OrderCreated struct {
	OrderID     uuid.UUID   `json:"orderId"`
	UserID      uuid.UUID   `json:"userId"`
	Items       []OrderItem `json:"items"`
	TotalAmount float64     `json:"totalAmount"`
}

type InventoryReserved struct {
	OrderID                uuid.UUID `json:"orderId"`
	InventoryReservationID uuid.UUID `json:"inventoryReservationId"`
}

type StockInsufficient struct {
	OrderID   uuid.UUID `json:"orderId"`
	ProductID uuid.UUID `json:"productId"`
	Reason    string    `json:"reason"`
}

type OrderCancelled struct {
	OrderID uuid.UUID `json:"orderId"`
	Reason  string    `json:"reason"`
}
