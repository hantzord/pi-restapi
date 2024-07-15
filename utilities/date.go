package utilities

import "time"

func StringToDate(date string) (time.Time, error) {
	newDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Time{}, err
	}
	return newDate, nil
}

func StringToTime(timeString string) (time.Time, error) {
	newTime, err := time.Parse("15:04", timeString)
	if err != nil {
		return time.Time{}, err
	}
	return newTime, nil
}
