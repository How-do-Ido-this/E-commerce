package product

import "context"

type Repository interface {
	GetByID(ctx context.Context, id string) (*Product, error)
	List(ctx context.Context) ([]*Product, error)
}
