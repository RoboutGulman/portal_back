package internalapi

import (
	"context"
	"errors"
	"portal_back/authentication/api/internalapi/model"
)

var UserAlreadyExists = errors.New("user with this email already exists")

type AuthRequestService interface {
	IsAuthenticated(token string) model.AuthValidationResult
	CreateNewUser(ctx context.Context, email string) error
	GetUserIdByEmail(ctx context.Context, email string) (*int, error)
}
