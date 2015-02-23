package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func _newMatchDate(matchDate string, baseDate *time.Time) (*time.Time, error) {
	matchDateRegexp := regexp.MustCompile(`(\d{1,2})/(\d{1,2})\s?\(.\)(.)(\d{1,2}):(\d{1,2}) KO`)

	m := matchDateRegexp.FindAllStringSubmatch(matchDate, -1)

	if len(m) == 0 {
		return nil, fmt.Errorf("match failed: %v", matchDate)
	}
	month, err := strconv.Atoi(m[0][1])
	if err != nil {
		return nil, err
	}
	day, err := strconv.Atoi(m[0][2])
	if err != nil {
		return nil, err
	}

	hour, err := strconv.Atoi(m[0][4])
	if err != nil {
		return nil, err
	}
	s := m[0][3]
	switch s {
	case "前":
	case "後":
		hour += 12
	case "深":
		hour += 24
	default:
		return nil, fmt.Errorf("not support:", s)
	}

	min, err := strconv.Atoi(m[0][5])
	if err != nil {
		return nil, err
	}

	year := baseDate.Year()
	// when current month is december, but match date seems to be january. it is next year case.
	if int(baseDate.Month()) > month {
		year++
	}
	d := time.Date(year, time.Month(month), day, hour, min, 0, 0, baseDate.Location())
	return &d, nil
}

func _newDateWithBaseDate(date string, live bool, baseDate *time.Time) (*time.Time, error) {
	dateRegexp := regexp.MustCompile(`(.)(\d{1,2}):(\d{1,2})`)

	m := dateRegexp.FindAllStringSubmatch(date, -1)

	if len(m) == 0 {
		return nil, fmt.Errorf("match failed: %v", date)
	}

	hour, err := strconv.Atoi(m[0][2])
	if err != nil {
		return nil, err
	}

	s := m[0][1]
	switch s {
	case "前", "深":
	case "後":
		hour += 12
	default:
		return nil, fmt.Errorf("not support:", s)
	}
	min, err := strconv.Atoi(m[0][3])
	if err != nil {
		return nil, err
	}

	diffDay := 0

	if live {
		// live program may start before kick off.

		// seems to be before day. ( for example 0 hour and 23 hour )
		if baseDate.Hour() < hour {
			diffDay--
		}

	} else {
		// not live program bay start after kick off.

		// seems to be before day. ( for example 23 hour and 0 hour )
		if baseDate.Hour() > hour {
			diffDay++
		}
	}

	d := time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day()+diffDay, hour, min, 0, 0, baseDate.Location())
	return &d, nil
}
