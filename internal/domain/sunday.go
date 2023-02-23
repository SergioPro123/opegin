package domain

import "time"

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

type SundayForm struct {
	Month         string        `json:"month"`
	Year          string        `json:"year"`
	Responsible   Responsible   `json:"responsible"`
	ImmediateBoss ImmediateBoss `json:"immediate_boss"`
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
