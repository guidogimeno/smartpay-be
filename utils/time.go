package utils

import "time"

const (
	layout = "2006-01-02"
)

type Clock struct {
	time time.Time
}

func NewClock() *Clock {
	return &Clock{
		time: time.Now(),
	}
}

func (t *Clock) AddDays(days int) *Clock {
	t.time = t.time.AddDate(0, 0, days)
	return t
}

func (t *Clock) AddMonths(months int) *Clock {
	t.time = t.time.AddDate(0, months, 0)
	return t
}

func (t *Clock) Format() string {
	return t.time.Format(layout)
}
