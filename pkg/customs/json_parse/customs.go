package jsonparse

import (
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
}

const hourLayout = "03:04"

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(hourLayout, s)
	return
}
