package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

func _parseCategories(doc *goquery.Document) []string {
	categories := []string{}

	doc.Find("div.subject-category ul li a").Each(func(i int, s *goquery.Selection) {
		h, exist := s.Attr("href")
		if !exist {
			return
		}
		h = strings.TrimPrefix(h, "#")
		categories = append(categories, h)
	})

	return categories
}

func _parseProgramsForCategory(doc *goquery.Document, category string, now time.Time) []Program {
	programs := []Program{}

	doc.Find(fmt.Sprintf(`div#%s tbody tr:not([class="foot"])`, category)).Each(func(i int, s *goquery.Selection) {
		match := strings.TrimSpace(s.Find(`td.match ul li`).Text())
		matchDate, err := _newMatchDate(strings.TrimSpace(s.Find(`td.match span.date`).Text()), &now)
		if err != nil {
			logrus.Debug("_newMatchDate err:", err)
			return
		}
		s.Find("td.channel").Each(func(i int, s *goquery.Selection) {
			live := s.Find(`span.date img[alt="LIVE"]`).Size() == 1
			date, err := _newProgramStartDate(strings.TrimSpace(s.Find(`span.date`).Text()), live, matchDate)
			if err != nil {
				logrus.Debug("_newDateWithBaseDate err:", err)
				return
			}
			channel := strings.TrimSpace(s.Find(`div.cs_nambar`).Text())
			p := Program{
				StartAt: date,
				Match:   match,
				Channel: channel,
				Live:    live,
			}
			programs = append(programs, p)
		})
	})

	return programs
}
