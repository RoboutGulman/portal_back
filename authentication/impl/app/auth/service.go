package auth

import (
	"context"
	"errors"
	"portal_back/authentication/impl/app/token"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserNotLogged = errors.New("user not logged")

type LoginData struct {
	Email    string
	Password string
}

type Service interface {
	GetSaltByEmail(ctx context.Context, email string) (string, error)
	Login(ctx context.Context, logData LoginData) (token.LoginData, error)
	Logout(ctx context.Context, token string) error
}

func NewService(repository Repository, tokenService token.Service) Service {
	return &service{repository: repository, tokenService: tokenService}
}

type service struct {
	repository   Repository
	tokenService token.Service
}

func (s *service) GetSaltByEmail(ctx context.Context, email string) (string, error) {
	return s.repository.GetSaltByEmail(ctx, email)
}

func (s *service) Login(ctx context.Context, logData LoginData) (token.LoginData, error) {
	password, err := s.repository.GetPasswordByEmail(ctx, logData.Email)
	if err != nil {
		return token.LoginData{}, err
	}
	if password != logData.Password {
		return token.LoginData{}, ErrUserNotFound
	}

	userID, err := s.repository.GetUserIDByEmail(ctx, logData.Email)
	if err != nil {
		return token.LoginData{}, err
	}

	tokens, err := s.tokenService.GenerateTokensForUser(ctx, userID)
	if err != nil {
		return token.LoginData{}, err
	}
	return token.LoginData{
		Tokens: tokens,
		User: token.User{
			Id: userID,
		},
	}, nil
}

func (s *service) Logout(ctx context.Context, userToken string) error {
	err := s.tokenService.RemoveToken(ctx, userToken)
	if err == token.ErrUserWithTokenNotFound {
		return ErrUserNotLogged
	}
	return err
}
