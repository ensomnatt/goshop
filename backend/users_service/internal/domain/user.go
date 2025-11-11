package domain

import (
	"context"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRepository interface {
	Create(user *User, ctx context.Context) error
	GetByEmail(email string, ctx context.Context) (*User, error)
}
