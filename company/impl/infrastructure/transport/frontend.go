package transport

import (
	"encoding/json"
	"errors"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"net/http"
	authInteralapi "portal_back/authentication/api/internalapi"
	"portal_back/company/api/frontend"
	"portal_back/company/impl/app/department"
	"portal_back/company/impl/app/employeeaccount"
	"portal_back/company/impl/infrastructure/mapper"
	"portal_back/core/network"
	"portal_back/roles/api/internalapi"
)

func NewServer(accountService employeeaccount.Service, departmentService department.Service, rolesService internalapi.RolesRequestService, authRequestService authInteralapi.AuthRequestService) frontendapi.ServerInterface {
	return &frontendServer{accountService, departmentService, rolesService, authRequestService}
}

type frontendServer struct {
	accountService     employeeaccount.Service
	departmentService  department.Service
	rolesService       internalapi.RolesRequestService
	authRequestService authInteralapi.AuthRequestService
}

func (f frontendServer) GetDepartments(w http.ResponseWriter, r *http.Request) {
	network.Wrap(f.authRequestService, w, r, func(info network.RequestInfo) {
		departments, err := f.departmentService.GetDepartments(r.Context(), info.CompanyId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(mapper.MapDepartmentsPreview(departments))

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

	})
}

func (f frontendServer) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	network.WrapWithBody(f.authRequestService, w, r, func(info network.RequestInfo, request frontendapi.DepartmentRequest) {
		err := f.departmentService.CreateDepartment(r.Context(), mapper.MapDepartmentRequest(request), info.CompanyId)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	})
}

func (f frontendServer) GetDepartment(w http.ResponseWriter, r *http.Request, departmentId int) {
	departmentWithEmployees, err := f.departmentService.GetDepartment(r.Context(), departmentId)
	if errors.Is(err, department.NotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(mapper.MapDepartment(departmentWithEmployees))

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

func (f frontendServer) DeleteDepartment(w http.ResponseWriter, r *http.Request, departmentId int) {
	err := f.departmentService.DeleteDepartment(r.Context(), departmentId)
	if errors.Is(err, department.NotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (f frontendServer) EditDepartment(w http.ResponseWriter, r *http.Request, departmentId int) {
	//TODO implement me
	panic("implement me")
}

func (f frontendServer) GetEmployees(w http.ResponseWriter, r *http.Request) {
	network.Wrap(f.authRequestService, w, r, func(info network.RequestInfo) {
		departments, err := f.departmentService.GetDepartments(r.Context(), info.CompanyId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var result frontendapi.GetAllEmployeesResponse

		result.IsEditable = true

		rootEmployees, err := f.accountService.GetRootEmployees(r.Context(), info.CompanyId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var employees []frontendapi.Employee
		for _, emp := range rootEmployees {
			employees = append(employees, mapper.MapEmployee(emp))
		}
		result.Employees = employees

		for _, dep := range departments {
			departmentWithEmployees, err := f.departmentService.GetDepartment(r.Context(), dep.Id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			result.Departments = append(result.Departments, mapper.MapDepartment(departmentWithEmployees))
		}

		resp, err := json.Marshal(result)

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

	})
}

func (f frontendServer) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	network.WrapWithBody(f.authRequestService, w, r, func(info network.RequestInfo, request frontendapi.EmployeeRequest) {
		err := f.accountService.CreateEmployee(r.Context(), mapper.MapEmployeeRequest(request), info.CompanyId)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	})
}

func (f frontendServer) MoveEmployeesToDepartment(w http.ResponseWriter, r *http.Request) {
	network.WrapWithBody(f.authRequestService, w, r, func(info network.RequestInfo, request frontendapi.MoveEmployeesRequest) {
		err := f.accountService.MoveEmployeesToDepartment(r.Context(), mapper.MapMoveEmployeeRequest(request))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	})
}

func (f frontendServer) GetEmployee(w http.ResponseWriter, r *http.Request, employeeId int) {
	employee, err := f.accountService.GetEmployee(
		r.Context(),
		employeeId)
	if errors.Is(err, employeeaccount.EmployeeNotFound) {
		w.WriteHeader(http.StatusNotFound)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		resp, err := json.Marshal(frontendapi.EmployeeWithConnections{
			Id:              employee.Id,
			Company:         frontendapi.Company{Id: employee.Company.Id, Name: employee.Company.Name},
			DateOfBirth:     openapi_types.Date{Time: employee.DateOfBirth},
			Departments:     mapper.MapDepartments(employee.Departments),
			Email:           employee.Email,
			FirstName:       employee.FirstName,
			SecondName:      employee.SecondName,
			Surname:         employee.Surname,
			Icon:            employee.Icon,
			TelephoneNumber: employee.TelephoneNumber,
			Roles:           mapper.MapRoles(employee.Roles),
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

func (f frontendServer) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	network.WrapWithBody(f.authRequestService, w, r, func(info network.RequestInfo, request frontendapi.DeleteEmployeeRequest) {
		deleteRequestArray := mapper.MapDeleteEmployeeRequest(request)
		for _, req := range deleteRequestArray {
			err := f.accountService.DeleteEmployee(r.Context(), req.EmployeeID, req.DepartmentID)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

	})
}

func (f frontendServer) EditEmployee(w http.ResponseWriter, r *http.Request, employeeId int) {
	//TODO implement me
	panic("implement me")
}
