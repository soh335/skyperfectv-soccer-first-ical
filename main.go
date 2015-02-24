package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/codegangsta/cli"
	"github.com/sirupsen/logrus"
)

const url = "http://soccer.skyperfectv.co.jp/static/first/"

func main() {
	app := cli.NewApp()
	app.Usage = "generate ical of soccer.skyperfectv.co.jp/static/first"
	app.Commands = []cli.Command{
		{
			Name:   "category",
			Usage:  "print categories",
			Action: CategoriesAction,
		},
		{
			Name:   "ical",
			Usage:  "print ical",
			Action: ICalAction,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "category",
				},
				cli.StringFlag{
					Name: "channel",
				},
				cli.StringFlag{
					Name:  "tzid",
					Value: "Asia/Tokyo",
				},
				cli.StringFlag{
					Name: "calname",
				},
				cli.BoolFlag{
					Name: "liveonly",
				},
			},
		},
	}
	app.Run(os.Args)
}

func CategoriesAction(c *cli.Context) {
	f := func(c *cli.Context) error {
		r, err := _fetch()
		if err != nil {
			return err
		}
		defer r.Close()

		doc, err := goquery.NewDocumentFromReader(r)
		if err != nil {
			return err
		}

		for _, category := range _parseCategories(doc) {
			fmt.Println(category)
		}
		return nil
	}
	if err := f(c); err != nil {
		logrus.Fatal(err)
	}
}

func ICalAction(c *cli.Context) {
	f := func(c *cli.Context) error {
		categories := strings.Split(c.String("category"), ",")
		if len(categories) < 1 {
			return fmt.Errorf("required category")
		}

		whiteSpace := regexp.MustCompile(`\s`)
		for _, category := range categories {
			if whiteSpace.MatchString(category) {
				return fmt.Errorf("should not contain white space in category")
			}
		}

		loc, err := time.LoadLocation(c.String("tzid"))
		if err != nil {
			return err
		}
		now := time.Now().In(loc)

		channels := strings.Split(c.String("channel"), ",")
		if len(channels) < 1 {
			return fmt.Errorf("required channel")
		}

		for i, channel := range channels {
			channels[i] = strings.TrimSpace(channel)
		}

		r, err := _fetch()
		if err != nil {
			return err
		}
		defer r.Close()

		doc, err := goquery.NewDocumentFromReader(r)
		if err != nil {
			return err
		}

		cal, err := ICal(doc, categories, channels, now, c.String("calname"), c.Bool("liveonly"))
		if err != nil {
			return err
		}
		return cal.Encode(os.Stdout)
	}
	if err := f(c); err != nil {
		logrus.Fatal(err)
	}
}

func _inslice(s string, strs []string) bool {
	for _, str := range strs {
		if s == str {
			return true
		}
	}
	return false
}

func _fetch() (io.ReadCloser, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code is not 200. got %d", res.StatusCode)
	}
	return res.Body, nil
}

type Program struct {
	StartAt *time.Time
	Match   string
	Channel string
	Live    bool
}
