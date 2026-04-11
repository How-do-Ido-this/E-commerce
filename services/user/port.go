package user

import "context"

type Repository interface {
	Save(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
}
