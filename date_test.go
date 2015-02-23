package main

import (
	"testing"
	"time"
)

func Test_newMatchDate(t *testing.T) {
	// same year
	type Spec struct {
		baseDate   time.Time
		matchDate  string
		expectDate time.Time
	}
	specs := []Spec{
		{
			time.Date(2014, time.Month(2), 21, 3, 10, 0, 0, asiaLoc),
			"2/21(土)前4:00 KO",
			time.Date(2014, time.Month(2), 21, 4, 0, 0, 0, asiaLoc),
		},
		{
			time.Date(2014, time.Month(2), 21, 10, 0, 0, 0, asiaLoc),
			"2/21(土)後11:00 KO",
			time.Date(2014, time.Month(2), 21, 23, 0, 0, 0, asiaLoc),
		},
		// next year
		{
			time.Date(2014, time.Month(12), 31, 3, 10, 0, 0, asiaLoc),
			"1/21(土)前4:00 KO",
			time.Date(2015, time.Month(1), 21, 4, 0, 0, 0, asiaLoc),
		},
		// 深 is next day
		{
			time.Date(2014, time.Month(2), 21, 3, 10, 0, 0, asiaLoc),
			"2/21(土)深1:00 KO",
			time.Date(2014, time.Month(2), 22, 1, 0, 0, 0, asiaLoc),
		},
	}

	for _, spec := range specs {
		d, err := _newMatchDate(spec.matchDate, &spec.baseDate)
		if err != nil {
			t.Errorf("%v:got err:%v", spec.matchDate, err)
		}
		if d.Unix() != spec.expectDate.Unix() {
			t.Errorf("should %v, but got %v", spec.expectDate, d)
		}
	}
}

func Test_newDateWithBaseDate(t *testing.T) {
	type Spec struct {
		baseDate   time.Time
		date       string
		live       bool
		expectDate time.Time
	}

	specs := []Spec{
		// live, same day
		{
			time.Date(2014, time.Month(2), 21, 4, 0, 0, 0, asiaLoc),
			"前3:10",
			true,
			time.Date(2014, time.Month(2), 21, 3, 10, 0, 0, asiaLoc),
		},
		// live, next day
		{
			time.Date(2014, time.Month(2), 21, 0, 0, 0, 0, asiaLoc),
			"後11:50",
			true,
			time.Date(2014, time.Month(2), 20, 23, 50, 0, 0, asiaLoc),
		},
		// live, next year
		{
			time.Date(2015, time.Month(1), 1, 0, 0, 0, 0, asiaLoc),
			"後11:50",
			true,
			time.Date(2014, time.Month(12), 31, 23, 50, 0, 0, asiaLoc),
		},
		// not live, same day
		{
			time.Date(2014, time.Month(2), 21, 1, 0, 0, 0, asiaLoc),
			"深1:40",
			false,
			time.Date(2014, time.Month(2), 21, 1, 40, 0, 0, asiaLoc),
		},
		// not live, next day
		{
			time.Date(2014, time.Month(2), 21, 23, 30, 0, 0, asiaLoc),
			"深1:40",
			false,
			time.Date(2014, time.Month(2), 22, 1, 40, 0, 0, asiaLoc),
		},
		// not live, next day
		{
			time.Date(2014, time.Month(2), 21, 23, 30, 0, 0, asiaLoc),
			"前4:40",
			false,
			time.Date(2014, time.Month(2), 22, 4, 40, 0, 0, asiaLoc),
		},
		// not live, next year
		{
			time.Date(2014, time.Month(12), 31, 23, 50, 0, 0, asiaLoc),
			"深1:40",
			false,
			time.Date(2015, time.Month(1), 1, 1, 40, 0, 0, asiaLoc),
		},
	}

	for _, spec := range specs {
		d, err := _newDateWithBaseDate(spec.date, spec.live, &spec.baseDate)
		if err != nil {
			t.Errorf("got err: %v spec.date:%v spec.live:%v", err, spec.date, spec.live)
		}
		if d.Unix() != spec.expectDate.Unix() {
			t.Errorf("should %v, but got %v", spec.expectDate, d)
		}
	}

}
