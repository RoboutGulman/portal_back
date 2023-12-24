package di

import (
	"github.com/jackc/pgx/v5"
	"net/http"
	"portal_back/authentication/api/frontend"
	"portal_back/authentication/api/internalapi"
	"portal_back/authentication/impl/app/auth"
	"portal_back/authentication/impl/app/authrequest"
	"portal_back/authentication/impl/app/token"
	"portal_back/authentication/impl/app/userrequest"
	"portal_back/authentication/impl/infrastructure/sql"
	"portal_back/authentication/impl/infrastructure/transport"
)

func InitAuthModule(conn *pgx.Conn) (internalapi.AuthRequestService, internalapi.UserRequestService) {
	repo := sql.NewTokenStorage(conn)
	tokenService := token.NewService(repo)

	authRepo := sql.NewAuthRepository(conn)
	authService := auth.NewService(authRepo, tokenService)
	server := transport.NewServer(authService, tokenService)
	authRequestService := authrequest.NewService()
	userRequestService := userrequest.NewService(authService)

	http.Handle("/authorization/", frontendapi.Handler(server))

	return authRequestService, userRequestService
}
