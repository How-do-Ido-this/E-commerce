package order

type OrderStatus string

const (
	StatusPending OrderStatus = "PENDING"
	StatusCreated OrderStatus = "CREATED"
	StatusFailed  OrderStatus = "FAILED"
)

type Order struct {
	ID        string
	UserID    string
	ProductID string
	Quantity  int
	Status    OrderStatus
}
