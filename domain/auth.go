package domain

import (
	"context"
)

// AuthUsecase represent the auth's usecases.
type AuthUsecase interface {
	Login(ctx context.Context, username, password string) (bool, error)
	ChangePassword(ctx context.Context) (bool, error)
	ResetPassword(ctx context.Context) (bool, error)
	GenerateToken(ctx context.Context) (string, error)
	GetTokenInfo(ctx context.Context) (string, error)
}
