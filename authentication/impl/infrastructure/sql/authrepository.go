package sql

import (
	"context"
	"github.com/jackc/pgx/v5"
	"portal_back/authentication/impl/app/auth"
)

func NewAuthRepository(conn *pgx.Conn) auth.Repository {
	return &repository{conn: conn}
}

type repository struct {
	conn *pgx.Conn
}

func (r *repository) CreateUser(ctx context.Context, email string) error {
	query := `
		INSERT INTO auth_user
		(password, salt, email)
		VALUES ($1, $2, $3)
	`
	_, err := r.conn.Exec(ctx, query, "qwerty123", "test_salt", email)
	return err
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (int, error) {
	query := `
		SELECT id FROM auth_user 
        WHERE email=$1
	`
	var id int
	err := r.conn.QueryRow(ctx, query, email).Scan(&id)
	if err == pgx.ErrNoRows {
		return 0, auth.ErrUserNotFound
	} else if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *repository) GetSaltByEmail(ctx context.Context, email string) (string, error) {
	query := `
		SELECT salt FROM auth_user 
        WHERE email=$1
	`
	var salt string
	err := r.conn.QueryRow(ctx, query, email).Scan(&salt)
	if err == pgx.ErrNoRows {
		return "", auth.ErrUserNotFound
	} else if err != nil {
		return "", err
	}
	return salt, nil
}

func (r *repository) GetPasswordByEmail(ctx context.Context, email string) (string, error) {
	query := `
		SELECT password FROM auth_user 
        WHERE email=$1
	`
	var password string
	err := r.conn.QueryRow(ctx, query, email).Scan(&password)
	if err == pgx.ErrNoRows {
		return "", auth.ErrUserNotFound
	} else if err != nil {
		return "", err
	}
	return password, nil
}

func (r *repository) GetUserIDByEmail(ctx context.Context, email string) (int, error) {
	query := `
		SELECT id FROM auth_user 
        WHERE email=$1
	`
	var userID int
	err := r.conn.QueryRow(ctx, query, email).Scan(&userID)
	if err == pgx.ErrNoRows {
		return 0, auth.ErrUserNotFound
	} else if err != nil {
		return 0, err
	}
	return userID, nil
}
