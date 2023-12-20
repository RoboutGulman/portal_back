package mapper

import (
	openapi_types "github.com/oapi-codegen/runtime/types"
	frontendapi "portal_back/company/api/frontend"
	"portal_back/company/impl/domain"
)

func MapEmployee(employee domain.Employee) frontendapi.Employee {
	return frontendapi.Employee{
		DateOfBirth:     openapi_types.Date{Time: employee.DateOfBirth},
		Email:           employee.Email,
		FirstName:       employee.FirstName,
		Icon:            employee.Icon,
		Id:              employee.Id,
		SecondName:      employee.SecondName,
		Surname:         employee.Surname,
		TelephoneNumber: employee.TelephoneNumber,
	}
}
