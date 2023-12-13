package auth

import "context"

type Repository interface {
	GetSaltByEmail(ctx context.Context, email string) (string, error)
	GetPasswordByEmail(ctx context.Context, email string) (string, error)
	GetUserIDByEmail(ctx context.Context, email string) (int, error)
}
