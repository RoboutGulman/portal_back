package sql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"portal_back/authentication/impl/app/auth"
)

func NewAuthRepository(conn *pgx.Conn) auth.Repository {
	return &repository{conn: conn}
}

type repository struct {
	conn *pgx.Conn
}

func (r *repository) GetCompanyByUserID(ctx context.Context, id int) (int, error) {
	query := `
		SELECT MIN(company.id)
	FROM public.company
	JOIN employeeaccount ON employeeaccount.companyid=company.id
	JOIN auth_user ON employeeaccount.userid= auth_user.id
	WHERE auth_user.id = $1
	GROUP BY auth_user.id
	`

	var companyID int
	err := r.conn.QueryRow(ctx, query, id).Scan(&companyID)

	//TODO: remove
	if errors.Is(err, pgx.ErrNoRows) {
		return 1, nil
	}

	if err != nil {
		return 0, err
	}

	return companyID, nil
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
