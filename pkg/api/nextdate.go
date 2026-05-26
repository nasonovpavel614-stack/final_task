package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

var (
	errInvalidRepeat = errors.New("некорректный формат правила повторения")
	errInvalidDate   = errors.New("некорректная дата")
)

// NextDate возвращает следующую дату выполнения задачи.
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errInvalidRepeat
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", errInvalidDate
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errInvalidRepeat
	}

	switch parts[0] {
	case "d":
		return nextByDays(date, now, parts)
	case "y":
		if len(parts) != 1 {
			return "", errInvalidRepeat
		}
		return nextByYear(date, now)
	case "w":
		return nextByWeekdays(date, now, parts)
	case "m":
		return nextByMonthDays(date, now, parts)
	default:
		return "", errInvalidRepeat
	}
}

func afterNow(date, now time.Time) bool {
	d := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	n := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	return d.After(n)
}

func nextByDays(date, now time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", errInvalidRepeat
	}
	interval, err := strconv.Atoi(parts[1])
	if err != nil || interval < 1 || interval > 400 {
		return "", errInvalidRepeat
	}

	for i := 0; i < 4000; i++ {
		date = date.AddDate(0, 0, interval)
		if afterNow(date, now) {
			return date.Format(DateFormat), nil
		}
	}
	return "", errInvalidRepeat
}

func nextByYear(date, now time.Time) (string, error) {
	for i := 0; i < 1000; i++ {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			return date.Format(DateFormat), nil
		}
	}
	return "", errInvalidRepeat
}

func nextByWeekdays(date, now time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", errInvalidRepeat
	}

	weekdays, err := parseWeekdays(parts[1])
	if err != nil {
		return "", err
	}

	return nextByMatch(date, now, func(t time.Time) bool {
		return weekdays[t.Weekday()]
	})
}

func parseWeekdays(s string) (map[time.Weekday]bool, error) {
	days := strings.Split(s, ",")
	if len(days) == 0 {
		return nil, errInvalidRepeat
	}

	weekdays := make(map[time.Weekday]bool)
	for _, d := range days {
		n, err := strconv.Atoi(strings.TrimSpace(d))
		if err != nil || n < 1 || n > 7 {
			return nil, errInvalidRepeat
		}
		weekdays[weekdayFromRule(n)] = true
	}
	return weekdays, nil
}

func weekdayFromRule(day int) time.Weekday {
	if day == 7 {
		return time.Sunday
	}
	return time.Weekday(day)
}

func nextByMonthDays(date, now time.Time, parts []string) (string, error) {
	if len(parts) < 2 || len(parts) > 3 {
		return "", errInvalidRepeat
	}

	days, err := parseMonthDayList(parts[1])
	if err != nil {
		return "", err
	}

	var months map[time.Month]bool
	if len(parts) == 3 {
		months, err = parseMonths(parts[2])
		if err != nil {
			return "", err
		}
	}

	return nextByMatch(date, now, func(t time.Time) bool {
		if months != nil && !months[t.Month()] {
			return false
		}
		return matchesMonthDay(t, days)
	})
}

func parseMonthDayList(s string) ([]int, error) {
	items := strings.Split(s, ",")
	if len(items) == 0 {
		return nil, errInvalidRepeat
	}

	days := make([]int, 0, len(items))

	for _, item := range items {
		n, err := strconv.Atoi(strings.TrimSpace(item))
		if err != nil {
			return nil, errInvalidRepeat
		}
		switch {
		case n >= 1 && n <= 31:
			days = append(days, n)
		case n == -1, n == -2:
			days = append(days, n)
		default:
			return nil, errInvalidRepeat
		}
	}

	return days, nil
}

func parseMonths(s string) (map[time.Month]bool, error) {
	items := strings.Split(s, ",")
	if len(items) == 0 {
		return nil, errInvalidRepeat
	}

	months := make(map[time.Month]bool)
	for _, item := range items {
		n, err := strconv.Atoi(strings.TrimSpace(item))
		if err != nil || n < 1 || n > 12 {
			return nil, errInvalidRepeat
		}
		months[time.Month(n)] = true
	}
	return months, nil
}

func matchesMonthDay(t time.Time, days []int) bool {
	lastDay := lastDayOfMonth(t)
	for _, d := range days {
		switch {
		case d > 0 && t.Day() == d:
			return true
		case d == -1 && t.Day() == lastDay:
			return true
		case d == -2 && t.Day() == lastDay-1:
			return true
		}
	}
	return false
}

func lastDayOfMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.Local).Day()
}

func nextByMatch(date, now time.Time, match func(time.Time) bool) (string, error) {
	for i := 0; i < 4000; i++ {
		date = date.AddDate(0, 0, 1)
		if afterNow(date, now) && match(date) {
			return date.Format(DateFormat), nil
		}
	}
	return "", errInvalidRepeat
}
