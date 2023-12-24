package sql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"portal_back/company/impl/app/department"
	"portal_back/company/impl/app/employeeaccount"
	"portal_back/company/impl/domain"
)

func NewEmployeeAccountRepository(conn *pgx.Conn) employeeaccount.Repository {
	return &employeeAccountRepository{conn: conn}
}

type employeeAccountRepository struct {
	conn *pgx.Conn
}

func (r employeeAccountRepository) DeleteEmployeeFromDepartment(ctx context.Context, id int, departmentID int) error {
	query := `
		DELETE FROM employee_department
		WHERE id=$1
		WHERE accountid = $1 ANDdepartmentid = $2
	`
	_, err := r.conn.Exec(ctx, query, id, departmentID)
	return err
}

func (r employeeAccountRepository) GetRootEmployees(ctx context.Context, companyID int) ([]domain.Employee, error) {
	query := `
		SELECT employeeaccount.id, employeeaccount.firstname, employeeaccount.secondname, employeeaccount.surname,
			employeeaccount.dateofbirth, auth_user.email, employeeaccount.telephonenumber
		FROM employeeaccount
		JOIN auth_user ON employeeaccount.userid = auth_user.id
		WHERE employeeaccount.id NOT IN (SELECT accountid FROM employee_department) AND companyID = $1
	`

	var departmentEmployees []domain.Employee
	rows, err := r.conn.Query(ctx, query, companyID)
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
			&departmentEmployee.TelephoneNumber)
		if err != nil {
			return departmentEmployees, err
		}
		departmentEmployees = append(departmentEmployees, departmentEmployee)
	}

	return departmentEmployees, nil
}

func (r employeeAccountRepository) GetCountOfEmployees(ctx context.Context, departmentID int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM employee_department
		WHERE departmentid = $1
	`

	var countOfDepartmentEmployees int
	err := r.conn.QueryRow(ctx, query, departmentID).Scan(&countOfDepartmentEmployees)
	if err == pgx.ErrNoRows {
		return 0, department.EmployeesNotFound
	} else if err != nil {
		return 0, err
	}

	return countOfDepartmentEmployees, nil
}

func (r employeeAccountRepository) GetDepartmentEmployees(ctx context.Context, departmentId int) ([]domain.Employee, error) {
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

func (r employeeAccountRepository) CreateEmployee(ctx context.Context, dto domain.EmployeeRequest, userId int, companyId int) error {
	query := `
		INSERT INTO employeeaccount
		(firstname, secondname, surname,
		telephonenumber, avatarurl, dateofbirth, userid, companyid, job)
		VALUES ($1, $2, $3,	$4, $5, $6, $7, $8, $9)
	`
	_, err := r.conn.Exec(ctx, query, dto.FirstName, dto.SecondName, dto.Surname, dto.TelephoneNumber, dto.Icon, dto.DateOfBirth, userId, companyId, "a")
	return err
}

func (r employeeAccountRepository) GetEmployee(ctx context.Context, id int) (domain.EmployeeWithConnections, error) {
	var employee domain.EmployeeWithConnections
	var rows pgx.Rows
	query := `
		SELECT employeeaccount.id, employeeaccount.firstname,
		       employeeaccount.secondname, employeeaccount.surname,
		       employeeaccount.dateofbirth, employeeaccount.telephonenumber,
		       employeeaccount.avatarurl, company.id, company.name, auth_user.email
		FROM employeeaccount
		LEFT JOIN company ON employeeaccount.companyid=company.id
		LEFT JOIN auth_user ON employeeaccount.userid=auth_user.id
        WHERE employeeaccount.id=$1
	`
	err := r.conn.QueryRow(ctx, query, id).Scan(&employee.Id, &employee.FirstName,
		&employee.SecondName, &employee.Surname,
		&employee.DateOfBirth, &employee.TelephoneNumber,
		&employee.Icon, &employee.Company.Id, &employee.Company.Name, &employee.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		return employee, employeeaccount.EmployeeNotFound
	}
	if err != nil {
		return employee, err
	}
	query = `
		SELECT department.id, department.name
		FROM department
		RIGHT JOIN employee_department ON employee_department.departmentid=department.id
		WHERE employee_department.accountid=$1
	`
	rows, err = r.conn.Query(ctx, query, id)
	for rows.Next() {
		var departmentInfo domain.DepartmentInfo
		err = rows.Scan(&departmentInfo.Id, &departmentInfo.Name)
		if err != nil {
			return employee, err
		}
		employee.Departments = append(employee.Departments, departmentInfo)
	}
	query = `
		SELECT role.id, role.title
		FROM role
		RIGHT JOIN employee_roles ON employee_roles.roleid=role.id
		WHERE employee_roles.accountid=$1
	`
	rows, err = r.conn.Query(ctx, query, id)
	for rows.Next() {
		var roleInfo domain.RoleInfo
		err = rows.Scan(&roleInfo.Id, &roleInfo.Name)
		if err != nil {
			return employee, err
		}
		employee.Roles = append(employee.Roles, roleInfo)
	}
	return employee, nil
}

func (r employeeAccountRepository) GetCompanyEmployee(ctx context.Context, userId int, companyId int) (domain.EmployeeWithConnections, error) {
	var employee domain.EmployeeWithConnections
	var err error
	var rows pgx.Rows
	query := `
		SELECT employeeaccount.id, employeeaccount.firstname,
		       employeeaccount.secondname, employeeaccount.surname,
		       employeeaccount.dateofbirth, employeeaccount.telephonenumber,
		       employeeaccount.avatarurl, company.id, company.name, auth_user.email
		FROM employeeaccount
		LEFT JOIN company ON employeeaccount.companyid=company.id
		LEFT JOIN auth_user ON employeeaccount.userid=auth_user.id
        WHERE employeeaccount.id=$1 AND employeeaccount.companyid=$2
	`
	err = r.conn.QueryRow(ctx, query, userId, companyId).Scan(&employee.Id, &employee.FirstName,
		&employee.SecondName, &employee.Surname,
		&employee.DateOfBirth, &employee.TelephoneNumber,
		&employee.Icon, &employee.Company.Id,
		&employee.Company.Name, &employee.Email)
	if err == pgx.ErrNoRows {
		return employee, employeeaccount.EmployeeNotFound
	} else if err != nil {
		return employee, err
	}
	query = `
		SELECT department.id, department.name
		FROM department
		RIGHT JOIN employee_department ON employee_department.departmentid=department.id
		RIGHT JOIN employeeaccount ON employeeaccount.id=employee_department.accountid
		WHERE employeeaccount.id=&1 AND department.companyid=$2
	`
	rows, err = r.conn.Query(ctx, query, userId, companyId)
	for rows.Next() {
		var departmentInfo domain.DepartmentInfo
		err = rows.Scan(departmentInfo.Id, departmentInfo.Name)
		if err != nil {
			return employee, err
		}
		employee.Departments = append(employee.Departments, departmentInfo)
	}
	query = `
		SELECT role.id, role.title
		FROM role
		RIGHT JOIN employee_roles ON employee_roles.roleid=role.id
		RIGHT JOIN employeeaccount ON employeeaccount.id=employee_roles.accountid
		WHERE employeeaccount.id=$1 AND employeeaccount.companyid=$2
	`
	rows, err = r.conn.Query(ctx, query, userId, companyId)
	for rows.Next() {
		var roleInfo domain.RoleInfo
		err = rows.Scan(roleInfo.Id, roleInfo.Name)
		if err != nil {
			return employee, err
		}
		employee.Roles = append(employee.Roles, roleInfo)
	}
	return employee, nil
}

func (r employeeAccountRepository) DeleteEmployee(ctx context.Context, id int) error {
	query := `
		DELETE FROM employeeaccount
		WHERE id=$1
	`
	_, err := r.conn.Exec(ctx, query, id)
	return err
}

func (r employeeAccountRepository) EditEmployee(ctx context.Context, id int, dto domain.EmployeeRequest) error {
	query := `
		UPDATE employeeaccount
		SET firstname=$2, secondname=$3, surname=$4,
		    telephonenumber=$5, avatarurl=$6, dateofbirth=$7,
		WHERE id=$1
	`
	_, err := r.conn.Exec(ctx, query, id, dto.FirstName, dto.SecondName, dto.Surname, dto.TelephoneNumber, dto.Icon, dto.DateOfBirth)
	return err
}

func (r employeeAccountRepository) MoveEmployeeToDepartment(ctx context.Context, employeeId int, departmentFromId int, departmentToId int) error {
	query := `
		UPDATE employee_department
		SET departmentid=$3
		WHERE accountid=$1 AND departmentid=$2
	`
	_, err := r.conn.Exec(ctx, query, employeeId, departmentFromId, departmentToId)
	return err
}

func (r employeeAccountRepository) AddEmployeeToDepartment(ctx context.Context, employeeId int, departmentId int) error {
	query := `
		INSERT INTO employee_department
		(accountid, departmentid)
		VALUES ($1, $2)
	`
	_, err := r.conn.Exec(ctx, query, employeeId, departmentId)
	return err
}
