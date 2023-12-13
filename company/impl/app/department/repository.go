package department

import (
	"context"
	"portal_back/company/impl/domain"
)

type Repository interface {
	GetDepartment(ctx context.Context, id int) (*domain.Department, error)
	GetChildDepartments(ctx context.Context, id int) ([]domain.Department, error)
	GetDepartmentEmployees(ctx context.Context, departmentId int) ([]domain.Employee, error)
	GetCountOfDepartmentEmployees(ctx context.Context, departmentId int) (int, error)
	GetRootCompanyDepartments(ctx context.Context, companyId int) ([]domain.Department, error)
}
