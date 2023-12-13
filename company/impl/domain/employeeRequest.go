package domain

import "time"

type EmployeeRequest struct {
	FirstName       string
	SecondName      string
	Surname         string
	DateOfBirth     time.Time
	Email           string
	Icon            string
	TelephoneNumber string
	DepartmentID    *int
	RoleIDs         []int
}
