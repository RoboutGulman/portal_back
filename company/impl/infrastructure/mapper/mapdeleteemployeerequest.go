package mapper

import (
	frontendapi "portal_back/company/api/frontend"
	"portal_back/company/impl/domain"
)

func MapDeleteEmployeeRequest(dto frontendapi.DeleteEmployeeRequest) []domain.DeleteEmployeeRequest {
	var result []domain.DeleteEmployeeRequest
	for _, request := range dto {
		result = append(result, domain.DeleteEmployeeRequest{
			DepartmentID: request.DepartmentID,
			EmployeeID:   request.EmployeeID,
		})
	}
	return result
}
