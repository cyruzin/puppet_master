package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type contextKey int

const ContextKeyID contextKey = iota

// Auth represent the auth's model.
type Auth struct {
	UserID    int64     `json:"user_id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password,omitempty" validate:"required,gte=8"`
	Roles     []string  `json:"roles"`
	TokenUUID uuid.UUID `json:"token_uuid"`
	Token     string    `json:"token,omitempty"`
	Revoked   bool      `json:"revoked"`
}

// AuthUsecase represent the auth's usecases.
type AuthUsecase interface {
	Authenticate(ctx context.Context, email, password string) (string, error)
	GenerateToken(auth *Auth, expiration time.Time) (string, error)
	// ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) (bool, error)
	// ResetPassword(ctx context.Context, email string) (bool, error)
}

// AuthRepository represent the auth's repository contract.
type AuthRepository interface {
	Authenticate(ctx context.Context, email string) (*User, error)
	// ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) (bool, error)
	// ResetPassword(ctx context.Context, email string) (bool, error)
}
