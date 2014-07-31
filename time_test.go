package timeutils

import (
	"testing"
	"time"
)

func TestWeekday(t *testing.T) {
	cases := map[time.Weekday]int{
		time.Monday:    0,
		time.Tuesday:   1,
		time.Wednesday: 2,
		time.Thursday:  3,
		time.Friday:    4,
		time.Saturday:  5,
		time.Sunday:    6,
	}
	for in, exp := range cases {
		out := Weekday(in)
		if exp != out {
			t.Fatalf("Weekday(%s) was %d; want %d\n", in, out, exp)
		}
	}
}

type ParseInput struct {
	year, month, day, hourMinute string
	loc                          *time.Location
}

func TestParse(t *testing.T) {
	// Note that Parse takes 0-indexed month.
	// TODO: test error handling.
	cases := map[ParseInput]time.Time{
		ParseInput{"2013", "06", "31", "23:59", time.UTC}: Must(ParseStd("2013-07-31 23:59")),
		ParseInput{"2013", "05", "30", "23:59", time.UTC}: Must(ParseStd("2013-06-30 23:59")),
	}
	for in, exp := range cases {
		out, err := Parse(in.year, in.month, in.day, in.hourMinute, in.loc)
		if err != nil {
			t.Fatalf("Parse(%v) got err %v\n", in, err)
		}
		if exp != out {
			t.Fatalf("Parse(%v) was %v; want %v\n", in, out, exp)
		}
	}
}

func TestStartOfWeek(t *testing.T) {
	cases := map[time.Time]time.Time{
		Must(time.Parse("2006-01-02 15:04:05.000", "2014-07-07 00:00:00.000")): Must(ParseStd("2014-07-07 00:00")),
		Must(ParseStd("2014-07-07 00:00")):                                     Must(ParseStd("2014-07-07 00:00")),
		Must(ParseStd("2014-07-07 00:01")):                                     Must(ParseStd("2014-07-07 00:00")),
		Must(ParseStd("2014-07-07 23:59")):                                     Must(ParseStd("2014-07-07 00:00")),
		Must(ParseStd("2014-07-13 23:59")):                                     Must(ParseStd("2014-07-07 00:00")),
		Must(time.Parse("2006-01-02 15:04:05.000", "2014-07-13 23:59:59.999")): Must(ParseStd("2014-07-07 00:00")),
		Must(time.Parse("2006-01-02 15:04:05.000", "2014-07-14 00:00:00.000")): Must(ParseStd("2014-07-14 00:00")),
	}
	for in, exp := range cases {
		out := StartOfWeek(in)
		if exp != out {
			t.Fatalf("StartOfWeek(%v) was %v; want %v\n", in, out, exp)
		}
	}
}
