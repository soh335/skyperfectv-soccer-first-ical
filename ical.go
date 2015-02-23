package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/soh335/ical"
)

func ICal(doc *goquery.Document, categories []string, channels []string, now time.Time, calname string, liveonly bool) (*ical.VCalendar, error) {
	programs_map := map[string][]Program{}
	for _, category := range categories {
		programs := []Program{}
		_programs := _parseProgramsForCategory(doc, category, now)
		for _, program := range _programs {
			if !_inslice(program.Channel, channels) {
				continue
			}
			programs = append(programs, program)
		}
		programs_map[category] = programs
	}

	components := []ical.VComponent{}
	for category, programs := range programs_map {
		for _, program := range programs {
			if liveonly && !program.Live {
				continue
			}

			live := ""
			if program.Live {
				live = "[live]"
			}
			h := sha1.New()
			io.WriteString(h, fmt.Sprintf("%v%v%v", program.StartAt.String(), program.Match, program.Channel))
			uid := hex.EncodeToString(h.Sum(nil))
			component := &ical.VEvent{
				UID:     uid,
				DTSTAMP: *program.StartAt,
				DTSTART: *program.StartAt,
				DTEND:   program.StartAt.Add(time.Hour * 2),
				SUMMARY: fmt.Sprintf("%v[%v][%v]%v", live, category, program.Channel, program.Match),
				TZID:    now.Location().String(),
			}
			components = append(components, component)
		}
	}

	cal := ical.NewBasicVCalendar()
	cal.PRODID = calname
	cal.X_WR_CALNAME = calname
	cal.X_WR_CALDESC = calname
	cal.X_WR_TIMEZONE = now.Location().String()
	cal.VComponent = components

	return cal, nil
}
