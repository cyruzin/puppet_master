package domain

import (
	"context"
	"time"
)

type contextKey int

const ContextKeyID contextKey = iota

// Auth represent the auth's model.
type Auth struct {
	UserID   int64  `json:"user_id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required,gte=8"`
	Role     string `json:"role"`
	Token    string `json:"token,omitempty"`
}

// AuthToken represent the token and refresh_token payload.
type AuthToken struct {
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// AuthUsecase represent the auth's usecases.
type AuthUsecase interface {
	Authenticate(ctx context.Context, email, password string) (*AuthToken, error)
	GenerateToken(claimKey string, claimValue interface{}, expiration time.Time) (string, error)
	// ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) (bool, error)
	// ResetPassword(ctx context.Context, email string) (bool, error)
	GetCache(ctx context.Context, key string, destination interface{}) error
}

// AuthRepository represent the auth's repository contract.
type AuthRepository interface {
	Authenticate(ctx context.Context, email string) (*User, error)
	// ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) (bool, error)
	// ResetPassword(ctx context.Context, email string) (bool, error)
}
