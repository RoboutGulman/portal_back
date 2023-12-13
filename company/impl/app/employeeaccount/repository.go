package employeeaccount

import (
	"context"
	"portal_back/company/impl/domain"
)

type Repository interface {
	CreateEmployee(ctx context.Context, dto domain.EmployeeRequest, userId int, companyId int) error
	GetEmployee(ctx context.Context, id int) (domain.EmployeeWithConnections, error)
	GetEmployeeByUserAndCompanyIds(ctx context.Context, userId int, companyId int) (domain.EmployeeWithConnections, error)
	DeleteEmployee(ctx context.Context, id int) error
	EditEmployee(ctx context.Context, id int, dto domain.EmployeeRequest) error
	MoveEmployeeToDepartment(ctx context.Context, employeeId int, departmentFromId int, departmentToId int) error
	AddEmployeeToDepartment(ctx context.Context, employeeId int, departmentId int) error
}
