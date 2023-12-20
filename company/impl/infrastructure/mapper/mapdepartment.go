package mapper

import (
	frontendapi "portal_back/company/api/frontend"
	"portal_back/company/impl/domain"
)

func MapDepartment(department domain.DepartmentWithEmployees) frontendapi.Department {
	var childDepartments []frontendapi.Department
	if department.Departments != nil {
		for _, d := range *department.Departments {
			childDepartments = append(childDepartments, MapDepartment(d))
		}
	}
	var employees []frontendapi.Employee
	for _, e := range department.Employees {
		employees = append(employees, MapEmployee(e))
	}

	var parentDepartment *frontendapi.DepartmentInfo
	if department.ParentDepartment != nil {
		parentDepartment = &frontendapi.DepartmentInfo{
			Id:   department.ParentDepartment.Id,
			Name: department.ParentDepartment.Name,
		}
	}

	var supervisor *frontendapi.Supervisor
	if department.Supervisor != nil {
		supervisor = &frontendapi.Supervisor{
			Id:   department.Supervisor.Id,
			Name: department.Supervisor.Name,
		}
	}

	result := frontendapi.Department{
		Departments:      &childDepartments,
		Employees:        employees,
		Id:               department.Id,
		Name:             department.Name,
		ParentDepartment: parentDepartment,
		Supervisor:       supervisor,
	}

	return result

}
