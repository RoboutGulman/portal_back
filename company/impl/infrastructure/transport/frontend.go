package transport

import (
	"encoding/json"
	"errors"
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

func (f frontendServer) GetEmployeeList(w http.ResponseWriter, r *http.Request) {
	network.WrapWithBody(f.authRequestService, w, r, func(info network.RequestInfo, request frontendapi.GetEmployeeListRequest) {
		ctx := r.Context()
		var employees []frontendapi.EmployeeWithConnections
		for _, id := range request {
			employee, err := f.accountService.GetEmployee(ctx, id)
			if errors.Is(err, employeeaccount.EmployeeNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			employees = append(employees, mapper.MapEmployeeWithConnections(employee))
		}
		resp, err := json.Marshal(employees)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = w.Write(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (f frontendServer) GetDepartments(w http.ResponseWriter, r *http.Request) {
	network.Wrap(f.authRequestService, w, r, func(info network.RequestInfo) {
		departments, err := f.departmentService.GetDepartments(r.Context(), info.CompanyId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := frontendapi.GetAllDepartmentsResponse{
			Departments: mapper.MapDepartmentsPreview(departments),
			IsEditable:  true,
		}

		resp, err := json.Marshal(response)

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

	response := frontendapi.GetDepartmentResponse{
		Department: mapper.MapDepartment(departmentWithEmployees),
		IsEditable: true,
	}

	resp, err := json.Marshal(response)

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
	network.WrapWithBody(f.authRequestService, w, r, func(info network.RequestInfo, request frontendapi.DepartmentRequest) {
		err := f.departmentService.EditDepartment(r.Context(), departmentId, mapper.MapDepartmentRequest(request))

		if errors.Is(err, department.NotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
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
		resp, err := json.Marshal(frontendapi.GetEmployeeResponse{
			Employee:   mapper.MapEmployeeWithConnections(employee),
			IsEditable: true,
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
	network.WrapWithBody(f.authRequestService, w, r, func(info network.RequestInfo, request frontendapi.EmployeeRequest) {
		err := f.accountService.EditEmployee(r.Context(), employeeId, mapper.MapEmployeeRequest(request))

		if errors.Is(err, employeeaccount.EmployeeNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

}
