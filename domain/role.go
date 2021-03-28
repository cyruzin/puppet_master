package domain

import (
	"context"
	"time"
)

// Role represent the role's model.
type Role struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Permissions []string  `json:"permissions"`
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

	AssignRole(ctx context.Context, roles []string, userID int64) error
	RemoveRole(ctx context.Context, roles []string, userID int64) error
	SyncRole(ctx context.Context, roles []string, userID int64) error
}
