package domain

type MoveEmployeesRequest struct {
	DepartmentToID int
	Employees      *[]MoveEmployeeInfo
}

type MoveEmployeeInfo struct {
	DepartmentFromID *int
	EmployeeID       int
}
