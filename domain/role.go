package domain

import (
	"context"
	"time"
)

// Role represent the role's model.
type Role struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at" db:"created_at"`
	CreatedAt   time.Time `json:"created_at" db:"updated_at"`
}

// RoleUsecase represent the role's usecases.
type RoleUsecase interface {
	Fetch(ctx context.Context) ([]*Role, error)
	GetByID(ctx context.Context, id int64) (*Role, error)
	Store(ctx context.Context, role *Role) error
	Update(ctx context.Context, role *Role) error
	Delete(ctx context.Context, id int64) error
}

// RoleRepository represent the role's repository contract.
type RoleRepository interface {
	Fetch(ctx context.Context) ([]*Role, error)
	GetByID(ctx context.Context, id int64) (*Role, error)
	Store(ctx context.Context, role *Role) error
	Update(ctx context.Context, role *Role) error
	Delete(ctx context.Context, id int64) error
}
