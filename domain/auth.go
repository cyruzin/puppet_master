package domain

import (
	"context"
)

type contextKey int

const ContextKeyID contextKey = iota

// Auth represent the auth's model.
type Auth struct {
	ID       int64  `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required,gte=8"`
	Token    string `json:"token,omitempty"`
}

// AuthUsecase represent the auth's usecases.
type AuthUsecase interface {
	Authenticate(ctx context.Context, email, password string) (string, error)

	GenerateToken() (string, error)

	// ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) (bool, error)
	// ResetPassword(ctx context.Context, email string) (bool, error)
}

// AuthRepository represent the auth's repository contract.
type AuthRepository interface {
	Authenticate(ctx context.Context, email string) (string, error)

	// ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) (bool, error)
	// ResetPassword(ctx context.Context, email string) (bool, error)

	// AssignRole(role string, id int64) error
	// RemoveRole(role string, id int64) error

	// GivePermissionTo(permission string, id int64) error
	// RemovePermissionTo(permission string, id int64) error
	// SyncPermission(permission string, id int64) error
}
