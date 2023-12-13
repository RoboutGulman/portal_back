package domain

type DepartmentPreview struct {
	CountOfEmployees int
	Departments      *[]DepartmentPreview
	Id               int
	Name             string
}
