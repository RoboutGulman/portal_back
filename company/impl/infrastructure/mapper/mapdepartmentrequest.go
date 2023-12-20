package mapper

import (
	frontendapi "portal_back/company/api/frontend"
	"portal_back/company/impl/domain"
)

func MapDepartmentRequest(department frontendapi.DepartmentRequest) domain.DepartmentRequest {
	return domain.DepartmentRequest{EmployeeIDs: department.EmployeeIDs, Name: department.Name, ParentDepartmentID: department.ParentDepartmentID, SupervisorID: department.SupervisorID}
}
