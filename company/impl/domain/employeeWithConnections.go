package domain

import "time"

type DepartmentInfo struct {
	Id   int
	Name string
}

type RoleInfo struct {
	Id   int
	Name string
}

type EmployeeWithConnections struct {
	Id              int
	FirstName       string
	SecondName      string
	Surname         string
	DateOfBirth     time.Time
	Email           string
	Icon            string
	TelephoneNumber string
	Company         struct {
		Id   int
		Name string
	}
	Departments []DepartmentInfo
	Roles       []RoleInfo
}
