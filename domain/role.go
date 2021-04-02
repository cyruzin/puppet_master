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
	Permissions []int     `json:"permissions,omitempty"`
	UpdatedAt   time.Time `json:"updated_at" db:"created_at"`
	CreatedAt   time.Time `json:"created_at" db:"updated_at"`
}

// RoleUsecase represent the role's usecases.
type RoleUsecase interface {
	Fetch(ctx context.Context) ([]*Role, error)
	GetByID(ctx context.Context, id int64) (*Role, error)
	Store(ctx context.Context, role *Role) (*Role, error)
	Update(ctx context.Context, role *Role) (*Role, error)
	Delete(ctx context.Context, id int64) error

	GetRoleByUserID(ctx context.Context, userID int64) (*Role, error)
}

// RoleRepository represent the role's repository contract.
type RoleRepository interface {
	Fetch(ctx context.Context) ([]*Role, error)
	GetByID(ctx context.Context, id int64) (*Role, error)
	Store(ctx context.Context, role *Role) (*Role, error)
	Update(ctx context.Context, role *Role) (*Role, error)
	Delete(ctx context.Context, id int64) error

	GetRoleByUserID(ctx context.Context, userID int64) (*Role, error)
	AssignRole(ctx context.Context, role int, userID int64) error
	RemoveRole(ctx context.Context, role int, userID int64) error
	SyncRole(ctx context.Context, role int, userID int64) error
}
