package sql

import (
	"context"
	"portal_back/role/impl/app/role"
	"portal_back/role/impl/domain"

	"github.com/jackc/pgx/v5"
)

type repository struct {
	conn *pgx.Conn
}

func rowsToArray(rows pgx.Rows) []domain.Role {
	var roles []domain.Role
	for rows.Next() {
		var role domain.Role
		rows.Scan(&role.Id, &role.Title, &role.Description, &role.RoleType)
		roles = append(roles, role)
	}
	return roles
}

// AssignRoleToUser implements role.RoleRepository.
func (repo *repository) AssignRoleToUser(context context.Context, roleId, userId int) error {
	query := `
		INSERT INTO employee_roles(account_id, role_id)
		SELECT $1, $2
		WHERE 
			NOT EXISTS(
				SELECT  account_id, role_id FROM employee_roles
				WHERE account_id=$1 AND role_id=$2
			)
	`
	_, err := repo.conn.Query(context, query, userId, roleId)
	return err
}

// GetAllRoles implements role.RoleRepository.
func (repo *repository) GetAllRoles(context context.Context) ([]domain.Role, error) {
	query := `
		SELECT role.id, role.title, role.description, role.role_type 
		FROM role
	`
	rows, err := repo.conn.Query(context, query)
	defer rows.Close()

	roles := rowsToArray(rows)
	if err == pgx.ErrNoRows {
		return []domain.Role{}, nil
	} else if err != nil {
		return nil, err
	}
	return roles, nil
}

// GetUserRoles implements role.RoleRepository.
func (repo *repository) GetUserRoles(context context.Context, userId int) ([]domain.Role, error) {
	query := `
		SELECT role.id, role.title, role.description, role.role_type  FROM role
		RIGHT JOIN employee_roles ON role.id=employee_roles.role_id
		AND employee_roles.account_id=$1
	`
	rows, err := repo.conn.Query(context, query, userId)
	defer rows.Close()

	roles := rowsToArray(rows)
	if err == pgx.ErrNoRows {
		return []domain.Role{}, nil
	} else if err != nil {
		return nil, err
	}
	return roles, nil
}

// RemoveRoleFromUser implements role.RoleRepository.
func (repo *repository) RemoveRoleFromUser(context context.Context, roleId, userId int) error {
	query := `
		DELETE FROM  employee_roles
		WHERE employee_roles.account_id=$1 
		AND employee_roles.role_id=$2
		RETURNING account_id
	`
	var id int
	return repo.conn.QueryRow(context, query, userId, roleId).Scan(&id)
}

func NewRepository(conn *pgx.Conn) role.RoleRepository {
	return &repository{conn: conn}
}
