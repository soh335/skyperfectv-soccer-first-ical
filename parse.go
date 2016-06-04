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

func _parseProgramsForCategory(doc *goquery.Document, category string, now time.Time, channelSelect string) []Program {
	programs := []Program{}

	doc.Find(fmt.Sprintf(`table%s%s tbody tr:not([class="foot"])`, channelSelect, category)).Each(func(i int, s *goquery.Selection) {
		match := strings.TrimSpace(s.Find(`td.match ul li`).Text())
		matchDate, err := _newMatchDate(strings.TrimSpace(s.Find(`td.match span.date`).Text()), &now)
		if err != nil {
			logrus.Debug("_newMatchDate err:", err, s.Text())
			return
		}
		s.Find("td.channel").Each(func(i int, s *goquery.Selection) {
			live := s.Find(`span.date img[alt="LIVE"]`).Size() == 1
			date, err := _newProgramStartDate(strings.TrimSpace(s.Find(`span.date`).Text()), live, matchDate)
			if err != nil {
				logrus.Debug("_newDateWithBaseDate err:", err, s.Text())
				return
			}
			channel := s.Find(`div.cs_nambar`)
			channel.Children().RemoveFiltered("span.channelName")
			p := Program{
				StartAt: date,
				Match:   match,
				Channel: strings.TrimSpace(channel.Text()),
				Live:    live,
			}
			programs = append(programs, p)
		})
	})

	return programs
}
