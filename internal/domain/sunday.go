package domain

import (
	jsonparse "devopegin/pkg/customs/json_parse"
	"time"
)

type Employee struct {
	Name       string
	Document   string
	Group      int
	Position   *Position
	ExtraHours []*ExtraHour
}

type Position struct {
	Name string `json:"name"`
}

type ExtraHour struct {
	ID            int
	Date          time.Time
	NumberOfHours int
}

type Sunday struct {
	Month               string
	Year                string
	Responsible         Responsible
	ImmediateBoss       ImmediateBoss
	EntryTime           jsonparse.CustomTime
	SundayEntryTime     jsonparse.CustomTime
	SundayDepartureTime jsonparse.CustomTime
	Justification       string
	CompanyImage        []byte
}

type SundayForm struct {
	Month               string               `json:"month"`
	Year                string               `json:"year"`
	Responsible         Responsible          `json:"responsible"`
	ImmediateBoss       ImmediateBoss        `json:"immediate_boss"`
	EntryTime           jsonparse.CustomTime `json:"entry_time"`
	SundayEntryTime     jsonparse.CustomTime `json:"entry_time_sunday"`
	SundayDepartureTime jsonparse.CustomTime `json:"departure_time_sunday"`
	Justification       string               `json:"justification"`
}

type Responsible struct {
	Name     string   `json:"name"`
	Position Position `json:"position"`
}
type ImmediateBoss struct {
	Name       string `json:"name"`
	Location   string `json:"location"`
	Department string `json:"department"`
}
