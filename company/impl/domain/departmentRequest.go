package domain

type DepartmentRequest struct {
	EmployeeIDs        []int
	Name               string
	ParentDepartmentID *int
	SupervisorID       int
}
