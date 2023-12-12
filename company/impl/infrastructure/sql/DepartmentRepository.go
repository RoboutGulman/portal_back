package sql

import (
	"context"
	"github.com/jackc/pgx/v5"
	"portal_back/company/impl/app/department"
	"portal_back/company/impl/domain"
)

func NewDepartmentRepository(conn *pgx.Conn) department.Repository {
	return &repository{conn: conn}
}

type departmentRepository struct {
	conn *pgx.Conn
}

func (r repository) GetDepartment(ctx context.Context, id int) (*domain.Department, error) {
	query := `
		SELECT department.id, department.name, department.parentdepartmentid, parentDepartment.name, 
		       department.supervisorid, employeeaccount.firstname
		FROM department
		LEFT JOIN employeeaccount ON department.supervisorid = employeeaccount.id
		LEFT JOIN department AS parentDepartment ON department.parentdepartmentid = department.id
		WHERE department.id = $1
	`

	var departmentInfo *domain.Department
	err := r.conn.QueryRow(ctx, query, departmentInfo).Scan(&departmentInfo.Id, &departmentInfo.Name, &departmentInfo.ParentDepartment.Id,
		&departmentInfo.ParentDepartment.Name, &departmentInfo.Supervisor.Id, &departmentInfo.Supervisor.Name)
	if err == pgx.ErrNoRows {
		return departmentInfo, department.DepartmentEmployeesNotFound
	} else if err != nil {
		return departmentInfo, err
	}

	return departmentInfo, nil
}

func (r repository) GetDepartmentEmployees(ctx context.Context, departmentId int) ([]domain.Employee, error) {
	query := `
		SELECT employeeaccount.id, employeeaccount.firstname, employeeaccount.secondname, employeeaccount.surname,
			employeeaccount.dateofbirth, auth_user.email
		FROM employeeaccount
		JOIN employee_department ON employeeaccount.id = employee_department.accountid
		JOIN auth_user ON employeeaccount.userid = auth_user.id
		WHERE employee_department.departmentid = $1
	`

	var departmentEmployees []domain.Employee
	rows, err := r.conn.Query(ctx, query)
	if err == pgx.ErrNoRows {
		return departmentEmployees, department.DepartmentEmployeesNotFound
	} else if err != nil {
		return departmentEmployees, err
	}
	defer rows.Close()

	for rows.Next() {
		var departmentEmployee domain.Employee
		err := rows.Scan(&departmentEmployee.Id, &departmentEmployee.FirstName, &departmentEmployee.SecondName,
			&departmentEmployee.Surname, &departmentEmployee.DateOfBirth, &departmentEmployee.Email,
			&departmentEmployee.Icon, &departmentEmployee.TelephoneNumber)
		if err != nil {
			return departmentEmployees, err
		}
		departmentEmployees = append(departmentEmployees, departmentEmployee)
	}

	return departmentEmployees, nil
}

func (r repository) GetChildDepartments(ctx context.Context, id int) ([]domain.Department, error) {
	query := `
		SELECT department.id, department.name, department.parentdepartmentid, parentDepartment.name, 
		       department.supervisorid, employeeaccount.firstname
		FROM department
		LEFT JOIN employeeaccount ON department.supervisorid = employeeaccount.id
		LEFT JOIN department AS parentDepartment ON department.parentdepartmentid = department.id
		WHERE department.parentdepartmentid = $1
	`

	var childDepartments []domain.Department
	rows, err := r.conn.Query(ctx, query)
	if err == pgx.ErrNoRows {
		return childDepartments, department.DepartmentEmployeesNotFound
	} else if err != nil {
		return childDepartments, err
	}
	defer rows.Close()

	for rows.Next() {
		var childDepartment domain.Department
		err := rows.Scan(&childDepartment.Id, &childDepartment.Name, &childDepartment.ParentDepartment.Id,
			&childDepartment.ParentDepartment.Name, &childDepartment.Supervisor.Id, &childDepartment.Supervisor.Name)
		if err != nil {
			return childDepartments, err
		}
		childDepartments = append(childDepartments, childDepartment)
	}

	return childDepartments, nil
}

func (r repository) GetCountOfDepartmentEmployees(ctx context.Context, departmentId int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM employee_department
		WHERE departmentid = $1
	`

	var countOfDepartmentEmployees int
	err := r.conn.QueryRow(ctx, query, departmentId).Scan(&countOfDepartmentEmployees)
	if err == pgx.ErrNoRows {
		return 0, department.DepartmentEmployeesNotFound
	} else if err != nil {
		return 0, err
	}

	return countOfDepartmentEmployees, nil
}

func (r repository) GetRootCompanyDepartments(ctx context.Context, companyId int) ([]domain.Department, error) {
	query := `
		SELECT department.id, department.name, department.supervisorid, employeeaccount.firstname
		FROM department
		LEFT JOIN employeeaccount ON department.supervisorid = employeeaccount.id
		WHERE department.companyid = $1 AND department.parentdepartmentid IS NULL
	`

	var rootCompanyDepartments []domain.Department
	rows, err := r.conn.Query(ctx, query, companyId)
	if err == pgx.ErrNoRows {
		return rootCompanyDepartments, department.DepartmentEmployeesNotFound
	} else if err != nil {
		return rootCompanyDepartments, err
	}
	defer rows.Close()

	for rows.Next() {
		var rootCompanyDepartment domain.Department
		rootCompanyDepartment.ParentDepartment.Id = 0
		rootCompanyDepartment.ParentDepartment.Name = ``
		err := rows.Scan(&rootCompanyDepartment.Id, &rootCompanyDepartment.Name, &rootCompanyDepartment.Supervisor.Id,
			&rootCompanyDepartment.Supervisor.Name)
		if err != nil {
			return rootCompanyDepartments, err
		}
		rootCompanyDepartments = append(rootCompanyDepartments, rootCompanyDepartment)
	}

	return rootCompanyDepartments, nil
}
