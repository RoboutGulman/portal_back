package employeeaccount

import (
	"context"
	"errors"
	"portal_back/authentication/api/internalapi"
	"portal_back/company/impl/domain"
)

var EmployeeAlreadyExists = errors.New("employee with this email already exists in your company")
var EmployeeNotFound = errors.New("employee not found")

type Service interface {
	GetEmployee(ctx context.Context, id int) (domain.EmployeeWithConnections, error)
	CreateEmployee(ctx context.Context, dto domain.EmployeeRequest, companyID int) error
	DeleteEmployee(ctx context.Context, id int, departmentID int) error
	EditEmployee(ctx context.Context, id int, dto domain.EmployeeRequest) error

	GetCountOfEmployees(ctx context.Context, departmentID int) (int, error)
	GetDepartmentEmployees(ctx context.Context, departmentID int) ([]domain.Employee, error)
	GetRootEmployees(ctx context.Context, companyID int) ([]domain.Employee, error)

	MoveEmployeesToDepartment(ctx context.Context, dto domain.MoveEmployeesRequest) error
	DeleteEmployeeFromDepartment(ctx context.Context, id int, departmentID int) error
}

func NewService(repository Repository, userService internalapi.UserRequestService) Service {
	return &service{repository: repository, userService: userService}
}

type service struct {
	repository  Repository
	userService internalapi.UserRequestService
}

func (s *service) GetRootEmployees(ctx context.Context, companyID int) ([]domain.Employee, error) {
	return s.repository.GetRootEmployees(ctx, companyID)
}

func (s *service) CreateEmployee(ctx context.Context, dto domain.EmployeeRequest, companyID int) error {
	userId, err := s.createUserOrGetExisting(ctx, dto.Email)
	if err != nil {
		return nil
	}

	_, err = s.repository.GetCompanyEmployee(ctx, userId, companyID)
	if !errors.Is(err, EmployeeNotFound) {
		return EmployeeAlreadyExists
	}

	err = s.repository.CreateEmployee(ctx, dto, userId, companyID)
	if err != nil {
		return err
	}

	if dto.DepartmentID != nil {
		createdEmployee, err := s.repository.GetCompanyEmployee(ctx, userId, companyID)
		if err != nil {
			return err
		}
		err = s.repository.AddEmployeeToDepartment(ctx, createdEmployee.Id, *dto.DepartmentID)
	}

	return nil
}

func (s *service) createUserOrGetExisting(ctx context.Context, email string) (int, error) {
	err := s.userService.CreateNewUser(ctx, email)

	if err != nil && !errors.Is(err, internalapi.UserAlreadyExists) {
		return -1, err
	}

	userId, err := s.userService.GetUserId(ctx, email)

	if err != nil {
		return -1, err
	}

	return userId, nil
}

func (s *service) GetEmployee(ctx context.Context, id int) (domain.EmployeeWithConnections, error) {
	return s.repository.GetEmployee(ctx, id)
}

func (s *service) DeleteEmployee(ctx context.Context, id int, departmentID int) error {
	err := s.repository.DeleteEmployeeFromDepartment(ctx, id, departmentID)
	if err != nil {
		return err
	}

	employee, err := s.repository.GetEmployee(ctx, id)
	if len(employee.Departments) < 1 {
		err = s.repository.DeleteEmployee(ctx, id)
		if err != nil {
			return err
		}
	}

	return nil

}

func (s *service) EditEmployee(ctx context.Context, id int, dto domain.EmployeeRequest) error {
	_, err := s.repository.GetEmployee(ctx, id)
	if err != nil {
		return err
	}
	return s.repository.EditEmployee(ctx, id, dto)
}

func (s *service) MoveEmployeesToDepartment(ctx context.Context, dto domain.MoveEmployeesRequest) error {
	for _, l := range *dto.Employees {
		if l.DepartmentFromID != nil {
			err := s.repository.MoveEmployeeToDepartment(ctx, l.EmployeeID, *l.DepartmentFromID, dto.DepartmentToID)
			if err != nil {
				return err
			}
		} else {
			err := s.repository.AddEmployeeToDepartment(ctx, l.EmployeeID, dto.DepartmentToID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *service) GetDepartmentEmployees(ctx context.Context, departmentID int) ([]domain.Employee, error) {
	return s.repository.GetDepartmentEmployees(ctx, departmentID)
}

func (s *service) GetCountOfEmployees(ctx context.Context, departmentID int) (int, error) {
	return s.repository.GetCountOfEmployees(ctx, departmentID)
}

func (s *service) DeleteEmployeeFromDepartment(ctx context.Context, id int, departmentID int) error {
	return s.repository.DeleteEmployeeFromDepartment(ctx, id, departmentID)
}
