package domain

import (
	"context"
)

type User struct {
	ID       string `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByCustomerEmail(ctx context.Context, email string) (Customers, error)
}
