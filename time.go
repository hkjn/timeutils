// Package timeutils provides some convenience functions around time.
package timeutils

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// MustLoadLoc loads the time.Location specified by the string, or panics.
func MustLoadLoc(l string) *time.Location {
	loc, err := time.LoadLocation(l)
	if err != nil {
		log.Fatal("bad location: %v\n", err)
	}
	return loc
}

// daysIn returns the number of days in a month for a given year.
func daysIn(m time.Month, year int) int {
	// This is equivalent to time.daysIn(m, year).o
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// Weekday gets the Monday-indexed number for the time.Weekday.
func Weekday(d time.Weekday) int {
	day := (d - 1) % 7
	if day < 0 {
		day += 7
	}
	return int(day)
}

// StartOfWeek returns the start of the current week for the time.
func StartOfWeek(t time.Time) time.Time {
	// Figure out number of days to back up until Mon:
	// Sun is 0 -> 6, Sat is 6 -> 5, etc.
	toMon := Weekday(t.Weekday())
	y, m, d := t.AddDate(0, 0, -int(toMon)).Date()
	// Result is 00:00:00 on that year, month, day.
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// GetWeek describes the week + year as a string, given week offset.
func GetWeek(offset int, loc *time.Location) string {
	now := time.Now().In(loc)
	fromTs := StartOfWeek(now.AddDate(0, 0, 7*offset))
	toTs := fromTs.AddDate(0, 0, 6)
	_, w := fromTs.ISOWeek()
	return fmt.Sprintf("Week %d (%s — %s)", w, fromTs.Format("Jan 2 2006"), toTs.Format("Jan 2 2006"))
}

// GetWeekOffset parses the week offset from current, which is 0.
func GetWeekOffset(week string) (int, error) {
	if week == "" {
		// Special-case for missing value, which is interpreted as "current week".
		return 0, nil
	}
	w, err := strconv.ParseInt(week, 10, 0)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("bad week: %s", err.Error()))
	}
	if w > 0 || w < -200000 {
		return 0, errors.New(fmt.Sprintf("bad week: %v", w))
	}
	return int(w), nil
}

// Parse extracts time from string-based info, with some constraints.
//
// The described time cannot be in the future, or more than 1000 years in the past
func Parse(year string, month string, day string, hourMinute string, loc *time.Location) (time.Time, error) {
	now := time.Now().In(loc)

	y64, err := strconv.ParseInt(year, 10, 0)
	y := int(y64)
	if err != nil {
		return time.Time{}, err
	}
	if y < now.Year()-1000 {
		return time.Time{}, errors.New(fmt.Sprintf("bad year; %d is too far in the past", y))
	}
	m, err := strconv.ParseInt(month, 10, 0)
	if err != nil {
		return time.Time{}, err
	}
	if m < 0 || m > 11 {
		return time.Time{}, errors.New(fmt.Sprintf("bad month: %d", m))
	}
	d64, err := strconv.ParseInt(day, 10, 0)
	d := int(d64)
	if err != nil {
		return time.Time{}, err
	}
	if d < 1 || d > daysIn(time.Month(m), y) {
		return time.Time{}, errors.New(fmt.Sprintf("bad day: %d", d))
	}
	parts := strings.Split(hourMinute, ":")
	if len(parts) != 2 {
		return time.Time{}, errors.New(fmt.Sprintf("bad hour/minute: %s", hourMinute))
	}
	h, err := strconv.ParseInt(parts[0], 10, 0)
	if err != nil {
		return time.Time{}, err
	}
	if h < 0 || h > 60 {
		return time.Time{}, errors.New(fmt.Sprintf("bad hour: %d", h))
	}
	min, err := strconv.ParseInt(parts[1], 10, 0)
	if err != nil {
		return time.Time{}, err
	}
	if min < 0 || min > 60 {
		return time.Time{}, errors.New(fmt.Sprintf("bad minute: %d", min))
	}

	// Month is +1 since time.Month is [1, 12].
	t := time.Time(time.Date(int(y), time.Month(m+1), int(d), int(h), int(min), 0, 0, loc))
	if t.After(now) {
		return time.Time{}, errors.New(fmt.Sprintf("bad time; %v is in the future", time.Time(t)))
	}
	return t, nil
}

// Selector holds info useful to make time selections.
type Selector struct {
	SelectedDay   int
	SelectedMonth time.Month
	SelectedYear  int
	SelectedTime  string
	Months        []time.Month
	Years         []int
	DaysInMonth   []int
}

// Create populates a Selector from given starting point.
func (s *Selector) Create(from time.Time) {
	days := make([]int, 31)
	for d := 0; d < 31; d++ { // TODO: Actual number of days / month (change dynamically on selection?).
		days[d] = d + 1
	}
	numYears := 5
	years := make([]int, numYears)
	for i := 0; i < numYears; i++ {
		years[i] = from.Year() - i
	}
	*s = Selector{
		SelectedYear:  from.Year(),
		SelectedMonth: from.Month() - 1, // -1 to give [0, 11]
		SelectedDay:   from.Day(),
		SelectedTime:  from.Format("15:04"),
		DaysInMonth:   days,
		Months:        make([]time.Month, 12),
		Years:         years,
	}
	for i := 1; i <= 12; i++ {
		s.Months[i-1] = time.Month(i)
	}
}
