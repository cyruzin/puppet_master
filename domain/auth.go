package domain

import (
	"context"
)

// Auth represent the auth's model.
type Auth struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required,min=8"`
}

// AuthUsecase represent the auth's usecases.
type AuthUsecase interface {
	Authenticate(ctx context.Context, email, password string) (string, error)
	// ChangePassword(
	// 	ctx context.Context,
	// 	userID int64,
	// 	oldPassword string,
	// 	newPassword string,
	// ) (bool, error)
	// ResetPassword(ctx context.Context, email string) (bool, error)
	GenerateToken() (string, error)
	ParseToken(token string) (interface{}, error)
}

// AuthRepository represent the auth's repository contract.
type AuthRepository interface {
	Authenticate(ctx context.Context, email, password string) (string, error)
	// ChangePassword(
	// 	ctx context.Context,
	// 	userID int64,
	// 	oldPassword string,
	// 	newPassword string,
	// ) (bool, error)
	// ResetPassword(ctx context.Context, email string) (bool, error)
}
