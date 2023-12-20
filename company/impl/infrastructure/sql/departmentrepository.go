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

type repository struct {
	conn *pgx.Conn
}

func (r repository) CreateDepartment(ctx context.Context, request domain.DepartmentRequest, companyId int) error {
	query := `
		INSERT INTO department
		(name, parentdepartmentid, companyid, supervisorid)
		VALUES ($1, $2, $3,	$4)
	`

	_, err := r.conn.Exec(ctx, query, request.Name, request.ParentDepartmentID, companyId, request.SupervisorID)
	return err
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
	departmentInfo.Supervisor = &domain.Supervisor{}
	departmentInfo.ParentDepartment = &domain.ParentDepartment{}
	err := r.conn.QueryRow(ctx, query, id).Scan(&departmentInfo.Id, &departmentInfo.Name, &departmentInfo.ParentDepartment.Id,
		&departmentInfo.ParentDepartment.Name, &departmentInfo.Supervisor.Id, &departmentInfo.Supervisor.Name)
	if err == pgx.ErrNoRows {
		return departmentInfo, department.EmployeesNotFound
	}
	if departmentInfo.Supervisor.Id == 0 {
		departmentInfo.Supervisor = nil
	}
	if departmentInfo.ParentDepartment.Id == 0 {
		departmentInfo.ParentDepartment = nil
	}
	//if err != nil {
	//	return departmentInfo, err
	//}

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
	rows, err := r.conn.Query(ctx, query, departmentId)
	if err == pgx.ErrNoRows {
		return departmentEmployees, department.EmployeesNotFound
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

func (r repository) GetCountOfEmployees(ctx context.Context, departmentId int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM employee_department
		WHERE departmentid = $1
	`

	var countOfDepartmentEmployees int
	err := r.conn.QueryRow(ctx, query, departmentId).Scan(&countOfDepartmentEmployees)
	if err == pgx.ErrNoRows {
		return 0, department.EmployeesNotFound
	} else if err != nil {
		return 0, err
	}

	return countOfDepartmentEmployees, nil
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
