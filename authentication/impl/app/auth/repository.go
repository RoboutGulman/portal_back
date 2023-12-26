package auth

import "context"

type Repository interface {
	GetSaltByEmail(ctx context.Context, email string) (string, error)
	GetPasswordByEmail(ctx context.Context, email string) (string, error)
	GetUserIDByEmail(ctx context.Context, email string) (int, error)
	CreateUser(ctx context.Context, email string) error
	GetUserByEmail(ctx context.Context, email string) (int, error)
	GetCompanyByUserID(ctx context.Context, id int) (int, error)
}
