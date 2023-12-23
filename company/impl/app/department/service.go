package department

import (
	"context"
	"errors"
	"portal_back/company/impl/app/employeeaccount"
	"portal_back/company/impl/domain"
	"portal_back/core/network"
)

type Service interface {
	GetDepartments(ctx context.Context, companyId int) ([]domain.DepartmentPreview, error)
	CreateDepartment(ctx context.Context, dto domain.DepartmentRequest, companyId int) error
	GetDepartment(ctx context.Context, id int) (domain.DepartmentWithEmployees, error)
	DeleteDepartment(ctx context.Context, id int) error
	EditDepartment(ctx context.Context, id int, dto domain.DepartmentRequest, requestInfo network.RequestInfo) error
	GetEmployees(ctx context.Context, companyId int) error
}

var EmployeesNotFound = errors.New("employees in this department not found")
var NotFound = errors.New("department not found")

func NewService(repository Repository, employeeService employeeaccount.Service) Service {
	return &service{repository: repository, employeeService: employeeService}
}

type service struct {
	repository      Repository
	employeeService employeeaccount.Service
}

func (s *service) GetDepartments(ctx context.Context, companyId int) ([]domain.DepartmentPreview, error) {
	rootDepartments, err := s.repository.GetCompanyDepartments(ctx, companyId)

	if err != nil {
		return nil, err
	}
	var resultDeps []domain.DepartmentPreview

	for _, dep := range rootDepartments {
		count, _ := s.employeeService.GetCountOfEmployees(ctx, dep.Id)
		var arr []domain.DepartmentPreview
		depPreview := domain.DepartmentPreview{
			CountOfEmployees: count,
			Departments:      &arr, Id: dep.Id, Name: dep.Name,
		}
		resultDeps = append(resultDeps, depPreview)
		err := s.findChildren(ctx, depPreview)
		if err != nil {
			return nil, err
		}
	}
	return resultDeps, nil
}

func (s *service) findChildren(ctx context.Context, department domain.DepartmentPreview) error {
	childDepartments, err := s.repository.GetChildDepartments(ctx, department.Id)
	if err != nil {
		return err
	}
	for _, dep := range childDepartments {
		count, _ := s.employeeService.GetCountOfEmployees(ctx, dep.Id)
		var arr []domain.DepartmentPreview
		normDep := domain.DepartmentPreview{
			CountOfEmployees: count,
			Departments:      &arr, Id: dep.Id, Name: dep.Name,
		}
		*department.Departments = append(*department.Departments, normDep)
		err := s.findChildren(ctx, normDep)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) CreateDepartment(ctx context.Context, dto domain.DepartmentRequest, companyId int) error {
	id, err := s.repository.CreateDepartment(ctx, dto, companyId)
	if err != nil {
		return err
	}
	var moveArray []domain.MoveEmployeeInfo
	request := domain.MoveEmployeesRequest{
		DepartmentToID: id,
		Employees:      &moveArray,
	}

	for _, id := range dto.EmployeeIDs {
		moveArray = append(moveArray, domain.MoveEmployeeInfo{
			DepartmentFromID: nil,
			EmployeeID:       id,
		})
	}

	err = s.employeeService.MoveEmployeesToDepartment(ctx, request)
	return err
}

func (s *service) GetDepartment(ctx context.Context, id int) (domain.DepartmentWithEmployees, error) {
	department, err := s.repository.GetDepartment(ctx, id)
	if err != nil {
		return domain.DepartmentWithEmployees{}, err
	}
	var arr []domain.DepartmentWithEmployees
	employees, _ := s.employeeService.GetDepartmentEmployees(ctx, id)
	result := domain.DepartmentWithEmployees{
		Departments:      &arr,
		Employees:        employees,
		Id:               department.Id,
		Name:             department.Name,
		ParentDepartment: department.ParentDepartment,
		Supervisor:       department.Supervisor,
	}

	err = s.findChildrenWithEmployees(ctx, result)
	if err != nil {
		return domain.DepartmentWithEmployees{}, err
	}

	return result, nil
}

func (s *service) findChildrenWithEmployees(ctx context.Context, department domain.DepartmentWithEmployees) error {
	childDepartments, err := s.repository.GetChildDepartments(ctx, department.Id)
	if err != nil {
		return err
	}
	for _, dep := range childDepartments {
		employees, _ := s.employeeService.GetDepartmentEmployees(ctx, dep.Id)
		var arr []domain.DepartmentWithEmployees
		normDep := domain.DepartmentWithEmployees{
			Departments: &arr,
			Employees:   employees,
			Id:          dep.Id,
			Name:        dep.Name,
			ParentDepartment: &domain.ParentDepartment{
				Id:   department.Id,
				Name: department.Name,
			},
			Supervisor: dep.Supervisor,
		}
		*department.Departments = append(*department.Departments, normDep)
		err := s.findChildrenWithEmployees(ctx, normDep)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) DeleteDepartment(ctx context.Context, id int) error {
	department, err := s.repository.GetDepartment(ctx, id)
	if err != nil {
		return err
	}

	childDepartments, err := s.repository.GetChildDepartments(ctx, id)
	employees, err := s.employeeService.GetDepartmentEmployees(ctx, id)
	if department.ParentDepartment != nil {
		for _, child := range childDepartments {
			err = s.repository.MoveDepartment(ctx, child.Id, department.ParentDepartment.Id)
			if err != nil {
				return err
			}
		}

		var moveArray []domain.MoveEmployeeInfo
		request := domain.MoveEmployeesRequest{
			DepartmentToID: department.ParentDepartment.Id,
			Employees:      &moveArray,
		}

		for _, employee := range employees {
			moveArray = append(moveArray, domain.MoveEmployeeInfo{
				DepartmentFromID: &id,
				EmployeeID:       employee.Id,
			})
		}

		err = s.employeeService.MoveEmployeesToDepartment(ctx, request)
	} else {
		for _, child := range childDepartments {
			err = s.repository.MoveDepartmentToRoot(ctx, child.Id)
			if err != nil {
				return err
			}
		}

		for _, employee := range employees {
			err := s.employeeService.DeleteEmployeeFromDepartment(ctx, employee.Id, id)
			if err != nil {
				return err
			}
		}
	}

	return s.repository.DeleteDepartment(ctx, id)
}

func (s *service) EditDepartment(ctx context.Context, id int, dto domain.DepartmentRequest, requestInfo network.RequestInfo) error {
	//TODO implement me
	panic("implement me")
}

func (s *service) GetEmployees(ctx context.Context, companyId int) error {
	//TODO implement me
	panic("implement me")
}
