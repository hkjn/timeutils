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
			t.Fatalf("Weekday(%s) => %d; want %d\n", in, out, exp)
		}
	}
}

func TestAsMillis(t *testing.T) {
	cases := []struct {
		t      time.Time
		offset int
		want   int
	}{
		{
			t:      Must(ParseStd("2013-07-31 23:59")),
			offset: 0,
			want:   1375315140000,
		},
		{
			t:      Must(ParseStd("2013-07-31 23:59")),
			offset: 2 * 60 * 60,
			want:   1375322340000,
		},
		{
			t:      Must(ParseStd("2013-07-31 23:59")),
			offset: -8 * 60 * 60,
			want:   1375286340000,
		},
	}
	for _, tt := range cases {
		got := AsMillis(tt.t, tt.offset)
		if got != tt.want {
			t.Fatalf("AsMillis(%v, %d) => %d; want %d\n", tt.t, tt.offset, got, tt.want)
		}
	}
}

func TestParse(t *testing.T) {
	// Note that Parse takes 0-indexed month.
	// TODO: test error handling?
	cases := []struct {
		year, month, day, hourMinute string
		loc                          *time.Location
		want                         time.Time
	}{
		{"2013", "06", "31", "23:59", time.UTC, Must(ParseStd("2013-07-31 23:59"))},
		{"2013", "05", "30", "23:59", time.UTC, Must(ParseStd("2013-06-30 23:59"))},
	}
	for _, tt := range cases {
		got, err := Parse(tt.year, tt.month, tt.day, tt.hourMinute, tt.loc)
		if err != nil {
			t.Fatalf("Parse(%q, %q, %q, %q, %v) got err %v\n", tt.year, tt.month, tt.day, tt.hourMinute, tt.loc, err)
		}
		if got != tt.want {
			t.Fatalf("Parse(%q, %q, %q, %q, %v) => %v; want %v\n", tt.year, tt.month, tt.day, tt.hourMinute, got, tt.want)
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
			t.Fatalf("StartOfWeek(%v) => %v; want %v\n", in, out, exp)
		}
	}
}
