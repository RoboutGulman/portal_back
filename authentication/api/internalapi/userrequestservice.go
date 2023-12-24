package internalapi

import (
	"context"
	"errors"
)

var UserAlreadyExists = errors.New("user with this email already exists")

type UserRequestService interface {
	CreateNewUser(ctx context.Context, email string) error
	GetUserId(ctx context.Context, email string) (int, error)
}
