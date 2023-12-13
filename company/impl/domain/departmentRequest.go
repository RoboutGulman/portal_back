package domain

type DepartmentRequest struct {
	EmployeeIDs        []int
	Name               string
	ParentDepartmentID *int
	Supervisor         *struct {
		Id   int
		Name string
	}
}
