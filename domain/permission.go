package domain

import (
	"context"
	"time"
)

// Permission represent the permission's model.
type Permission struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
	UpdatedAt   time.Time `json:"updated_at" db:"created_at"`
	CreatedAt   time.Time `json:"created_at" db:"updated_at"`
}

// PermissionUsecase represent the permission's usecases.
type PermissionUsecase interface {
	Fetch(ctx context.Context) ([]*Permission, error)
	GetByID(ctx context.Context, id int64) (*Permission, error)
	Store(ctx context.Context, permission *Permission) error
	Update(ctx context.Context, permission *Permission) error
	Delete(ctx context.Context, id int64) error
}

// PermissionRepository represent the permission's repository contract.
type PermissionRepository interface {
	Fetch(ctx context.Context) ([]*Permission, error)
	GetByID(ctx context.Context, id int64) (*Permission, error)
	Store(ctx context.Context, permission *Permission) error
	Update(ctx context.Context, permission *Permission) error
	Delete(ctx context.Context, id int64) error
}
