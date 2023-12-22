package department

import (
	"context"
	"portal_back/company/impl/domain"
)

type Repository interface {
	GetDepartment(ctx context.Context, id int) (domain.Department, error)
	GetChildDepartments(ctx context.Context, id int) ([]domain.Department, error)
	GetDepartmentEmployees(ctx context.Context, departmentId int) ([]domain.Employee, error)
	GetCountOfEmployees(ctx context.Context, departmentId int) (int, error)
	GetCompanyDepartments(ctx context.Context, companyId int) ([]domain.Department, error)
	CreateDepartment(ctx context.Context, request domain.DepartmentRequest, companyId int) error
	DeleteDepartment(ctx context.Context, id int) error
}
