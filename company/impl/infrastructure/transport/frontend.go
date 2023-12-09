package transport

import (
	"encoding/json"
	"errors"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"net/http"
	"portal_back/company/api/frontend"
	"portal_back/company/impl/app/employeeaccount"
	"portal_back/company/impl/domain"
	"portal_back/roles/api/internalapi"
)

func NewServer(accountService employeeaccount.Service, rolesService internalapi.RolesRequestService) frontendapi.ServerInterface {
	return &frontendServer{accountService, rolesService}
}

type frontendServer struct {
	accountService employeeaccount.Service
	rolesService   internalapi.RolesRequestService
}

func (f frontendServer) GetCompanyDepartments(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (f frontendServer) CreateNewDepartment(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
	//проверка на права (обращение в модуль ролей)
}

func (f frontendServer) GetDepartment(w http.ResponseWriter, r *http.Request, departmentId int) {
	//TODO implement me
	panic("implement me")
}

func (f frontendServer) DeleteDepartment(w http.ResponseWriter, r *http.Request, departmentId int) {
	//TODO implement me
	panic("implement me")
	//проверка на права (обращение в модуль ролей)
}

func (f frontendServer) EditDepartment(w http.ResponseWriter, r *http.Request, departmentId int) {
	//TODO implement me
	panic("implement me")
	//проверка на права (обращение в модуль ролей)
}

func (f frontendServer) GetCompanyDepartmentsWithEmployees(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (f frontendServer) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
	//проверка на права (обращение в модуль ролей)
}

func (f frontendServer) MoveEmployeesToDepartment(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
	//проверка на права (обращение в модуль ролей)
}

func (f frontendServer) GetEmployee(w http.ResponseWriter, r *http.Request, employeeId int) {
	employee, err := f.accountService.GetEmployee(
		r.Context(),
		employeeId)
	if errors.Is(err, employeeaccount.EmployeeNotFound) {
		w.WriteHeader(http.StatusNotFound)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else if err == nil {
		resp, err := json.Marshal(frontendapi.EmployeeWithConnections{
			Id:              employee.Id,
			Company:         frontendapi.Company{Id: employee.Company.Id, Name: employee.Company.Name},
			DateOfBirth:     openapi_types.Date{Time: employee.DateOfBirth},
			Departments:     mapDepartments(employee.Departments),
			Email:           employee.Email,
			FirstName:       employee.FirstName,
			SecondName:      employee.SecondName,
			Surname:         employee.Surname,
			Icon:            employee.Icon,
			TelephoneNumber: employee.TelephoneNumber,
			Roles:           mapRoles(employee.Roles),
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func mapDepartments(departments []domain.DepartmentInfo) []frontendapi.DepartmentInfo {
	var result []frontendapi.DepartmentInfo
	for _, d := range departments {
		result = append(result, frontendapi.DepartmentInfo{Id: d.Id, Name: d.Name})
	}
	return result
}

func mapRoles(departments []domain.RoleInfo) []frontendapi.RoleInfo {
	var result []frontendapi.RoleInfo
	for _, d := range departments {
		result = append(result, frontendapi.RoleInfo{Id: d.Id, Name: d.Name})
	}
	return result
}

func (f frontendServer) DeleteEmployee(w http.ResponseWriter, r *http.Request, employeeId int) {
	//TODO implement me
	panic("implement me")
	//проверка на права (обращение в модуль ролей)
}

func (f frontendServer) EditEmployee(w http.ResponseWriter, r *http.Request, employeeId int) {
	//TODO implement me
	panic("implement me")
	//проверка на права (обращение в модуль ролей)
}
