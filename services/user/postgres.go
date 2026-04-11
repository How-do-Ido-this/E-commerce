package user

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

func (r *PostgresRepository) Save(ctx context.Context, u *User) error {
	query := `INSERT INTO users (id, name, email) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, u.ID, u.Name, u.Email)
	if err != nil {
		return fmt.Errorf("error saving user: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*User, error) {
	query := `SELECT id, name, email FROM users WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var u User
	err := row.Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		return nil, fmt.Errorf("error querying user: %w", err)
	}
	return &u, nil
}
