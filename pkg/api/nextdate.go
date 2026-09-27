package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	return date.Format(DateFormat) > now.Format(DateFormat)
}

func parseCSVInts(s string) ([]int, error) {
	if s == "" {
		return nil, fmt.Errorf("empty list")
	}
	parts := strings.Split(s, ",")
	values := make([]int, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			return nil, fmt.Errorf("invalid empty value")
		}
		n, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid number %q", part)
		}
		values = append(values, n)
	}
	return values, nil
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func matchesMonthRule(date time.Time, days []int, months map[int]bool) bool {
	if len(months) > 0 && !months[int(date.Month())] {
		return false
	}
	last := daysInMonth(date.Year(), date.Month())
	for _, day := range days {
		switch day {
		case -1:
			if date.Day() == last {
				return true
			}
		case -2:
			if date.Day() == last-1 {
				return true
			}
		default:
			if day <= last && date.Day() == day {
				return true
			}
		}
	}
	return false
}

// NextDate returns the next task date strictly later than now.
func NextDate(now time.Time, dstart, repeat string) (string, error) {
	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid start date: %w", err)
	}
	if repeat == "" {
		return "", fmt.Errorf("repeat rule is empty")
	}

	fields := strings.Fields(repeat)
	if len(fields) == 0 {
		return "", fmt.Errorf("repeat rule is empty")
	}

	switch fields[0] {
	case "y":
		if len(fields) != 1 {
			return "", fmt.Errorf("invalid yearly rule")
		}
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				return date.Format(DateFormat), nil
			}
		}

	case "d":
		if len(fields) != 2 {
			return "", fmt.Errorf("invalid daily rule")
		}
		interval, err := strconv.Atoi(fields[1])
		if err != nil || interval < 1 || interval > 400 {
			return "", fmt.Errorf("invalid day interval")
		}
		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				return date.Format(DateFormat), nil
			}
		}

	case "w":
		if len(fields) != 2 {
			return "", fmt.Errorf("invalid weekly rule")
		}
		values, err := parseCSVInts(fields[1])
		if err != nil {
			return "", err
		}
		weekdays := make(map[int]bool, len(values))
		for _, v := range values {
			if v < 1 || v > 7 {
				return "", fmt.Errorf("invalid weekday")
			}
			weekdays[v] = true
		}
		candidate := date.AddDate(0, 0, 1)
		for {
			wd := int(candidate.Weekday())
			if wd == 0 {
				wd = 7
			}
			if weekdays[wd] && afterNow(candidate, now) {
				return candidate.Format(DateFormat), nil
			}
			candidate = candidate.AddDate(0, 0, 1)
		}

	case "m":
		if len(fields) < 2 || len(fields) > 3 {
			return "", fmt.Errorf("invalid monthly rule")
		}
		days, err := parseCSVInts(fields[1])
		if err != nil {
			return "", err
		}
		for _, d := range days {
			if d == 0 || d < -2 || d > 31 {
				return "", fmt.Errorf("invalid month day")
			}
		}

		months := map[int]bool{}
		if len(fields) == 3 {
			monthValues, err := parseCSVInts(fields[2])
			if err != nil {
				return "", err
			}
			for _, m := range monthValues {
				if m < 1 || m > 12 {
					return "", fmt.Errorf("invalid month")
				}
				months[m] = true
			}
		}

		candidate := date.AddDate(0, 0, 1)
		for {
			if matchesMonthRule(candidate, days, months) && afterNow(candidate, now) {
				return candidate.Format(DateFormat), nil
			}
			candidate = candidate.AddDate(0, 0, 1)
		}

	default:
		return "", fmt.Errorf("unsupported repeat rule")
	}
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
		return
	}

	now := time.Now()
	if value := r.FormValue("now"); value != "" {
		parsed, err := time.Parse(DateFormat, value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		now = parsed
	}

	next, err := NextDate(now, r.FormValue("date"), r.FormValue("repeat"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, err := fmt.Fprint(w, next); err != nil {
		log.Printf("write next date: %v", err)
	}
}
