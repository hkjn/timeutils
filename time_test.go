package timeutils

import (
	"testing"
	"time"
)

func TestWeekday(t *testing.T) {
	cases := []struct {
		in   time.Weekday
		want int
	}{
		{
			in:   time.Monday,
			want: 0,
		},
		{
			in:   time.Tuesday,
			want: 1,
		},
		{
			in:   time.Wednesday,
			want: 2,
		},
		{
			in:   time.Thursday,
			want: 3,
		},
		{
			in:   time.Friday,
			want: 4,
		},
		{
			in:   time.Saturday,
			want: 5,
		},
		{
			in:   time.Sunday,
			want: 6,
		},
	}
	for i, tt := range cases {
		out := Weekday(tt.in)
		if tt.want != out {
			t.Errorf("[%d] Weekday(%s) => %d; want %d\n", i, tt.in, out, tt.want)
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
	for i, tt := range cases {
		got := AsMillis(tt.t, tt.offset)
		if got != tt.want {
			t.Errorf("[%d] AsMillis(%v, %d) => %d; want %d\n", i, tt.t, tt.offset, got, tt.want)
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
	for i, tt := range cases {
		got, err := Parse(tt.year, tt.month, tt.day, tt.hourMinute, tt.loc)
		if err != nil {
			t.Errorf("[%d] Parse(%q, %q, %q, %q, %v) got err %v\n", i, tt.year, tt.month, tt.day, tt.hourMinute, tt.loc, err)
		}
		if got != tt.want {
			t.Errorf("[%d] Parse(%q, %q, %q, %q, %v) => %v; want %v\n", i, tt.year, tt.month, tt.day, tt.hourMinute, tt.loc, got, tt.want)
		}
	}
}

func TestStartOfWeek(t *testing.T) {
	cases := []struct {
		in   time.Time
		want time.Time
	}{
		{
			in:   Must(time.Parse("2006-01-02 15:04:05.000", "2014-07-07 00:00:00.000")),
			want: Must(ParseStd("2014-07-07 00:00")),
		},
		{
			in:   Must(ParseStd("2014-07-07 00:00")),
			want: Must(ParseStd("2014-07-07 00:00")),
		},
		{
			in:   Must(ParseStd("2014-07-07 00:01")),
			want: Must(ParseStd("2014-07-07 00:00")),
		},
		{
			in:   Must(ParseStd("2014-07-07 23:59")),
			want: Must(ParseStd("2014-07-07 00:00")),
		},
		{
			in:   Must(ParseStd("2014-07-13 23:59")),
			want: Must(ParseStd("2014-07-07 00:00")),
		},
		{
			in:   Must(time.Parse("2006-01-02 15:04:05.000", "2014-07-13 23:59:59.999")),
			want: Must(ParseStd("2014-07-07 00:00")),
		},
		{
			in:   Must(time.Parse("2006-01-02 15:04:05.000", "2014-07-14 00:00:00.000")),
			want: Must(ParseStd("2014-07-14 00:00")),
		},
	}
	for i, tt := range cases {
		out := StartOfWeek(tt.in)
		if tt.want != out {
			t.Errorf("[%d] StartOfWeek(%v) => %v; want %v\n", i, tt.in, out, tt.want)
		}
	}
}

func TestDescDuration(t *testing.T) {
	cases := map[time.Duration]string{
		time.Millisecond:                 "0.0 sec ago",
		time.Millisecond * 49:            "0.0 sec ago",
		time.Millisecond * 50:            "0.1 sec ago",
		time.Second:                      "1.0 sec ago",
		time.Millisecond * 500:           "0.5 sec ago",
		time.Second * 59:                 "59.0 sec ago",
		time.Second*60 - 1:               "60.0 sec ago",
		time.Second*60 + 1:               "1.0 min ago",
		time.Minute:                      "1.0 min ago",
		time.Minute * 60:                 "1.0 hrs ago",
		time.Hour*24 - 1:                 "24.0 hrs ago",
		time.Hour * 24:                   "1.0 days ago",
		time.Hour*24*10e4 + time.Hour*12: "100000.5 days ago",
	}
	for in, exp := range cases {
		out := DescDuration(in)
		if exp != out {
			t.Errorf("DescDuration(%v) => %q; want %q\n", in, out, exp)
		}
	}
}
