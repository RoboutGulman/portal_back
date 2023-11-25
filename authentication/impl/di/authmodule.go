package di

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"net/http"
	"os"
	"portal_back/authentication/api/frontend"
	"portal_back/authentication/api/internalapi"
	"portal_back/authentication/impl/app/auth"
	"portal_back/authentication/impl/app/authrequest"
	"portal_back/authentication/impl/app/token"
	"portal_back/authentication/impl/infrastructure/sql"
	"portal_back/authentication/impl/infrastructure/transport"
)

func InitAuthModule() (internalapi.AuthRequestService, *pgx.Conn) {
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "password"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "app"
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", dbUser, dbPassword, dbHost, dbName)

	conn, _ := pgx.Connect(context.Background(), connStr)

	repo := sql.NewTokenStorage(conn)
	tokenService := token.NewService(repo)

	authRepo := sql.NewAuthRepository(conn)
	authService := auth.NewService(authRepo, tokenService)
	server := transport.NewServer(authService, tokenService)
	authRequestService := authrequest.NewService()

	http.Handle("/authorization/", frontendapi.Handler(server))

	return authRequestService, conn
}
