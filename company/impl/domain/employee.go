package domain

import "time"

type Employee struct {
	Id              int
	FirstName       string
	SecondName      string
	Surname         string
	DateOfBirth     time.Time
	Email           string
	Icon            string
	TelephoneNumber string
}
