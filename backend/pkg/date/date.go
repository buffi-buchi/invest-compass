package date

import (
	"errors"
	"time"
)

type Date struct {
	time time.Time
}

func NewDate(year int, month time.Month, day int) Date {
	return Date{
		time: time.Date(year, month, day, 0, 0, 0, 0, time.UTC),
	}
}

func (d Date) Time() time.Time {
	return d.time
}

func (d Date) String() string {
	return d.time.Format(time.DateOnly)
}

func (d Date) MarshalJSON() ([]byte, error) {
	data := make([]byte, 0, len(time.DateOnly)+len(`""`))
	data = append(data, '"')
	data = d.time.AppendFormat(data, time.DateOnly)
	data = append(data, '"')
	return data, nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("Date.UnmarshalJSON: input is not a JSON string")
	}

	data = data[len(`"`) : len(data)-len(`"`)]

	if len(data) != len(time.DateOnly) {
		return errors.New("Date.UnmarshalJSON: input is not a JSON string")
	}

	t, err := time.Parse(time.DateOnly, string(data))
	if err != nil {
		return err
	}

	d.time = t

	return nil
}
