package mapper

import (
	frontendapi "portal_back/company/api/frontend"
	"portal_back/company/impl/domain"
)

func MapEmployeeRequest(department frontendapi.EmployeeRequest) domain.EmployeeRequest {
	return domain.EmployeeRequest{
		FirstName:       department.FirstName,
		SecondName:      department.SecondName,
		Surname:         department.Surname,
		DateOfBirth:     department.DateOfBirth.Time,
		Email:           department.Email,
		Icon:            department.Icon,
		TelephoneNumber: department.TelephoneNumber,
		DepartmentID:    department.DepartmentID,
		RoleIDs:         department.RoleIDs,
	}
}
