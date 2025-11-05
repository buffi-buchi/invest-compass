package date

import (
	"errors"
	"time"
)

type Date struct {
	year  int
	month time.Month
	day   int
}

func (d *Date) MarshalJSON() ([]byte, error) {
	t := time.Date(d.year, d.month, d.day, 0, 0, 0, 0, time.UTC)
	data := make([]byte, 0, len(time.DateOnly)+len(`""`))
	data = append(data, '"')
	data = t.AppendFormat(data, time.DateOnly)
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

	d.year, d.month, d.day = t.Date()

	return nil
}

func (d *Date) Time() time.Time {
	return time.Date(d.year, d.month, d.day, 0, 0, 0, 0, time.UTC)
}

func (d *Date) String() string {
	return d.Time().Format(time.DateOnly)
}
