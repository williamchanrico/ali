package util

import (
	"time"
)

// ParseTimeStr will parse date string into time
func ParseTimeStr(str string) (time.Time, error) {
	if str == "" {
		return time.Time{}, nil
	}

	timeLayout := []string{
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04Z",
	}

	var err error
	var t time.Time
	for i := range timeLayout {
		t, err = time.Parse(timeLayout[i], str)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, err
}
