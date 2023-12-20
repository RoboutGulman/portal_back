package domain

type DepartmentWithEmployees struct {
	Departments      *[]DepartmentWithEmployees
	Employees        []Employee
	Id               int
	Name             string
	ParentDepartment *ParentDepartment
	Supervisor       *Supervisor
}
