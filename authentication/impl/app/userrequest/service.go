package userrequest

import (
	"context"
	"errors"
	"portal_back/authentication/api/internalapi"
	"portal_back/authentication/impl/app/auth"
)

func NewService(authService auth.Service) internalapi.UserRequestService {
	return &service{authService: authService}
}

type service struct {
	authService auth.Service
}

func (s service) CreateNewUser(ctx context.Context, email string) error {
	err := s.authService.CreateUser(ctx, email)
	if errors.Is(err, auth.ErrUserAlreadyExists) {
		return internalapi.UserAlreadyExists
	}
	return err
}

func (s service) GetUserId(ctx context.Context, email string) (int, error) {
	return s.authService.GetUserByEmail(ctx, email)
}
