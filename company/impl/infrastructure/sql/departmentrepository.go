package sql

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5"
	"portal_back/company/impl/app/department"
	"portal_back/company/impl/domain"
)

func NewDepartmentRepository(conn *pgx.Conn) department.Repository {
	return &repository{conn: conn}
}

type repository struct {
	conn *pgx.Conn
}

func (r repository) EditDepartment(ctx context.Context, id int, dto domain.DepartmentRequest) error {
	query := `
		UPDATE department
		SET (name, parentdepartmentid, supervisorid) = ($1, $2, $3)
		WHERE id=$4
	`
	_, err := r.conn.Exec(ctx, query, dto.Name, dto.ParentDepartmentID, dto.SupervisorID, id)
	return err
}

func (r repository) MoveDepartment(ctx context.Context, departmentID int, newParentID int) error {
	query := `
		UPDATE department
		SET parentdepartmentid=$2
		WHERE id=$1
	`
	_, err := r.conn.Exec(ctx, query, departmentID, newParentID)
	return err
}

func (r repository) MoveDepartmentToRoot(ctx context.Context, id int) error {
	query := `
		UPDATE department
		SET parentdepartmentid=NULL
		WHERE id=$1
	`
	_, err := r.conn.Exec(ctx, query, id)
	return err
}

func (r repository) DeleteDepartment(ctx context.Context, id int) error {
	query := `
		DELETE FROM department
		WHERE id=$1
	`
	_, err := r.conn.Exec(ctx, query, id)
	return err
}

func (r repository) CreateDepartment(ctx context.Context, request domain.DepartmentRequest, companyId int) (int, error) {
	query := `
		INSERT INTO department
		(name, parentdepartmentid, companyid, supervisorid)
		VALUES ($1, $2, $3,	$4)
		RETURNING id
	`
	lastInsertId := 0
	err := r.conn.QueryRow(ctx, query, request.Name, request.ParentDepartmentID, companyId, request.SupervisorID).Scan(&lastInsertId)
	return lastInsertId, err
}

func (r repository) GetDepartment(ctx context.Context, id int) (domain.Department, error) {
	query := `
		SELECT department.id, department.name, department.parentdepartmentid, parentDepartment.name, 
		       department.supervisorid, employeeaccount.firstname
		FROM department
		LEFT JOIN employeeaccount ON department.supervisorid = employeeaccount.id
		LEFT JOIN department AS parentDepartment ON department.parentdepartmentid = department.id
		WHERE department.id = $1
	`

	var departmentInfo domain.Department

	var supervisorID, parentDepartmentID sql.NullInt32
	var supervisorName, parentDepartmentName sql.NullString

	err := r.conn.QueryRow(ctx, query, id).Scan(&departmentInfo.Id, &departmentInfo.Name, &parentDepartmentID,
		&parentDepartmentName, &supervisorID, &supervisorName)
	if err == pgx.ErrNoRows {
		return departmentInfo, department.NotFound
	}
	if err != nil {
		return departmentInfo, err
	}
	if parentDepartmentID.Valid {
		departmentInfo.ParentDepartment = &domain.ParentDepartment{
			Id:   int(parentDepartmentID.Int32),
			Name: parentDepartmentName.String,
		}
	}
	if supervisorID.Valid {
		departmentInfo.Supervisor = &domain.Supervisor{
			Id:   int(supervisorID.Int32),
			Name: supervisorName.String,
		}
	}

	return departmentInfo, nil
}

func (r repository) GetChildDepartments(ctx context.Context, id int) ([]domain.Department, error) {
	query := `
		SELECT department.id, department.name, department.parentdepartmentid, parentDepartment.name, 
		       department.supervisorid, employeeaccount.firstname
		FROM department
		LEFT JOIN employeeaccount ON department.supervisorid = employeeaccount.id
		LEFT JOIN department AS parentDepartment ON department.parentdepartmentid = parentDepartment.id
		WHERE department.parentdepartmentid = $1
	`

	var childDepartments []domain.Department
	rows, err := r.conn.Query(ctx, query, id)
	if err == pgx.ErrNoRows {
		return childDepartments, department.EmployeesNotFound
	} else if err != nil {
		return childDepartments, err
	}
	defer rows.Close()

	for rows.Next() {
		var childDepartment domain.Department
		childDepartment.Supervisor = &domain.Supervisor{}
		childDepartment.ParentDepartment = &domain.ParentDepartment{}
		err := rows.Scan(&childDepartment.Id, &childDepartment.Name, &childDepartment.ParentDepartment.Id,
			&childDepartment.ParentDepartment.Name, &childDepartment.Supervisor.Id, &childDepartment.Supervisor.Name)
		if err != nil {
			return childDepartments, err
		}
		childDepartments = append(childDepartments, childDepartment)
	}

	return childDepartments, nil
}

func (r repository) GetCompanyDepartments(ctx context.Context, companyId int) ([]domain.Department, error) {
	query := `
		SELECT department.id, department.name, department.supervisorid, employeeaccount.firstname
		FROM department
		LEFT JOIN employeeaccount ON department.supervisorid = employeeaccount.id
		WHERE department.companyid = $1 AND department.parentdepartmentid IS NULL
	`

	var rootCompanyDepartments []domain.Department
	rows, err := r.conn.Query(ctx, query, companyId)
	if err == pgx.ErrNoRows {
		return rootCompanyDepartments, department.EmployeesNotFound
	} else if err != nil {
		return rootCompanyDepartments, err
	}
	defer rows.Close()

	for rows.Next() {
		var rootCompanyDepartment domain.Department
		rootCompanyDepartment.ParentDepartment = nil
		var supervisor domain.Supervisor
		rootCompanyDepartment.Supervisor = &supervisor
		err := rows.Scan(&rootCompanyDepartment.Id, &rootCompanyDepartment.Name, &rootCompanyDepartment.Supervisor.Id,
			&rootCompanyDepartment.Supervisor.Name)
		if err != nil {
			return rootCompanyDepartments, err
		}
		rootCompanyDepartments = append(rootCompanyDepartments, rootCompanyDepartment)
	}

	return rootCompanyDepartments, nil
}
