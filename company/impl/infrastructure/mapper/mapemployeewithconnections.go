package mapper

import (
	openapi_types "github.com/oapi-codegen/runtime/types"
	frontendapi "portal_back/company/api/frontend"
	"portal_back/company/impl/domain"
)

func MapEmployeeWithConnections(employee domain.EmployeeWithConnections) frontendapi.EmployeeWithConnections {
	return frontendapi.EmployeeWithConnections{
		Id:              employee.Id,
		Company:         frontendapi.Company{Id: employee.Company.Id, Name: employee.Company.Name},
		DateOfBirth:     openapi_types.Date{Time: employee.DateOfBirth},
		Departments:     MapDepartments(employee.Departments),
		Email:           employee.Email,
		FirstName:       employee.FirstName,
		SecondName:      employee.SecondName,
		Surname:         employee.Surname,
		Icon:            employee.Icon,
		TelephoneNumber: employee.TelephoneNumber,
		Roles:           MapRoles(employee.Roles),
	}
}
