package product

// Service implementa la lógica de negocio para productos.
type Service struct {
	repo Repository
}

// NewService crea una nueva instancia de Service.
func NewService(r Repository) *Service {
	return &Service{repo: r}
}
