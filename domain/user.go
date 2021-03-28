package domain

import (
	"context"
	"time"
)

// User represent the user's model.
type User struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Email       string    `json:"email" validate:"required,email"`
	Password    string    `json:"-"`
	Roles       []string  `json:"roles"`
	Permissions []string  `json:"permissions"`
	UpdatedAt   time.Time `json:"updated_at" db:"created_at"`
	CreatedAt   time.Time `json:"created_at" db:"updated_at"`
}

// UserUsecase represent the user's usecases.
type UserUsecase interface {
	Fetch(ctx context.Context) ([]*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	Store(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
}

// UserRepository represent the user's repository contract.
type UserRepository interface {
	Fetch(ctx context.Context) ([]*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	Store(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
}
