package department

import (
	"context"
	"portal_back/company/impl/domain"
)

type Repository interface {
	GetDepartment(ctx context.Context, id int) (domain.Department, error)
	CreateDepartment(ctx context.Context, request domain.DepartmentRequest, companyId int) (int, error)
	DeleteDepartment(ctx context.Context, id int) error
	EditDepartment(ctx context.Context, id int, dto domain.DepartmentRequest) error

	GetChildDepartments(ctx context.Context, id int) ([]domain.Department, error)
	GetCompanyDepartments(ctx context.Context, companyID int) ([]domain.Department, error)

	MoveDepartment(ctx context.Context, departmentID int, newParentID int) error
	MoveDepartmentToRoot(ctx context.Context, id int) error
}
