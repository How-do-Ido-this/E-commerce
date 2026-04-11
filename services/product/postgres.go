package product

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Product, error) {
	query := `SELECT id, name, price FROM products WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var p Product
	err := row.Scan(&p.ID, &p.Name, &p.Price)
	if err != nil {
		return nil, fmt.Errorf("error querying product: %w", err)
	}
	return &p, nil
}

func (r *PostgresRepository) List(ctx context.Context) ([]*Product, error) {
	query := `SELECT id, name, price FROM products`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error listing products: %w", err)
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	return products, nil
}
