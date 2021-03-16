package domain

import (
	"context"
)

// AuthUsecase represent the auth's usecases.
type AuthUsecase interface {
	Login(ctx context.Context, username, password string) (bool, error)
	HashPassword(ctx context.Context, password string, cost int) (string, error)
	CheckPasswordHash(ctx context.Context, password, hash string) bool
	ChangePassword(ctx context.Context) (bool, error)
	ResetPassword(ctx context.Context) (bool, error)
	GenerateToken(ctx context.Context) (string, error)
}
