package domain

type MoveEmployeesRequest struct {
	DepartmentToID int
	Employees      []struct {
		DepartmentFromID *int
		EmployeeID       int
	}
}
