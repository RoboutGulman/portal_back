package domain

type Department struct {
	Id               int
	Name             string
	ParentDepartment struct {
		Id   int
		Name string
	}
	Supervisor struct {
		Id   int
		Name string
	}
}
