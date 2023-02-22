package domain

import "time"

type Employee struct {
	Name       string
	Document   string
	Group      int
	Location   *Location
	ExtraHours []*ExtraHour
}

type Location struct {
	Name string
}

type ExtraHour struct {
	ID            int
	Date          time.Time
	NumberOfHours int
}
