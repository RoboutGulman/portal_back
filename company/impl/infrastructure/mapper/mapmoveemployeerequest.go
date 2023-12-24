package mapper

import (
	frontendapi "portal_back/company/api/frontend"
	"portal_back/company/impl/domain"
)

func MapMoveEmployeeRequest(request frontendapi.MoveEmployeesRequest) domain.MoveEmployeesRequest {
	var infoArray []domain.MoveEmployeeInfo
	for _, info := range request.Employees {
		infoArray = append(infoArray, domain.MoveEmployeeInfo{
			DepartmentFromID: info.DepartmentFromID,
			EmployeeID:       *info.EmployeeID,
		})
	}
	return domain.MoveEmployeesRequest{
		DepartmentToID: *request.DepartmentToID,
		Employees:      &infoArray,
	}
}
