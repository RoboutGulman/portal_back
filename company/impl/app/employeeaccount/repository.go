package employeeaccount

import (
	"context"
	"portal_back/company/impl/domain"
)

type Repository interface {
	CreateEmployee(ctx context.Context, dto domain.EmployeeRequest, userID int, companyID int) error
	GetEmployee(ctx context.Context, id int) (domain.EmployeeWithConnections, error)
	DeleteEmployee(ctx context.Context, id int) error
	DeleteEmployeeFromDepartment(ctx context.Context, id int, departmentID int) error
	EditEmployee(ctx context.Context, id int, dto domain.EmployeeRequest) error

	GetRootEmployees(ctx context.Context, companyID int) ([]domain.Employee, error)
	GetDepartmentEmployees(ctx context.Context, departmentID int) ([]domain.Employee, error)
	GetCountOfEmployees(ctx context.Context, departmentID int) (int, error)
	GetCompanyEmployee(ctx context.Context, userID int, companyID int) (domain.EmployeeWithConnections, error)

	MoveEmployeeToDepartment(ctx context.Context, employeeID int, departmentFromID int, departmentToID int) error
	AddEmployeeToDepartment(ctx context.Context, employeeID int, departmentID int) error
}
