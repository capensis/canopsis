package datetime

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

const (
	DurationUnitSecond = "s"
	DurationUnitMinute = "m"
	DurationUnitHour   = "h"
	DurationUnitDay    = "d"
	DurationUnitWeek   = "w"
	DurationUnitMonth  = "M"
	DurationUnitYear   = "y"

	daysInWeek      = 7
	hoursInDay      = 24
	minutesInHour   = 60
	secondsInMinute = minutesInHour
	secondsInHour   = secondsInMinute * minutesInHour
	secondsInDay    = hoursInDay * secondsInHour
	minutesInDay    = hoursInDay * minutesInHour
	hoursInWeek     = daysInWeek * hoursInDay
	minutesInWeek   = hoursInWeek * minutesInHour
	secondsInWeek   = minutesInWeek * secondsInMinute
)

// DurationWithUnit represent duration with user-preferred units
type DurationWithUnit struct {
	Value int64  `bson:"value" json:"value" binding:"required,min=1"`
	Unit  string `bson:"unit" json:"unit" binding:"required,oneof=s m h d w M y"`
}

func NewDurationWithUnit(value int64, unit string) DurationWithUnit {
	return DurationWithUnit{
		Value: value,
		Unit:  unit,
	}
}

func (d DurationWithUnit) AddTo(t CpsTime) CpsTime {
	var r time.Time

	switch d.Unit {
	case DurationUnitSecond:
		r = t.Add(time.Duration(d.Value) * time.Second)
	case DurationUnitMinute:
		r = t.Add(time.Duration(d.Value) * time.Minute)
	case DurationUnitHour:
		r = t.Add(time.Duration(d.Value) * time.Hour)
	case DurationUnitDay:
		r = t.AddDate(0, 0, int(d.Value))
	case DurationUnitWeek:
		r = t.AddDate(0, 0, 7*int(d.Value))
	case DurationUnitMonth:
		r = t.AddDate(0, int(d.Value), 0)
	case DurationUnitYear:
		r = t.AddDate(int(d.Value), 0, 0)
	default:
		r = t.Add(time.Duration(d.Value) * time.Second)
	}

	return CpsTime{Time: r}
}

func (d DurationWithUnit) SubFrom(t CpsTime) CpsTime {
	var r time.Time

	switch d.Unit {
	case DurationUnitSecond:
		r = t.Add(-time.Duration(d.Value) * time.Second)
	case DurationUnitMinute:
		r = t.Add(-time.Duration(d.Value) * time.Minute)
	case DurationUnitHour:
		r = t.Add(-time.Duration(d.Value) * time.Hour)
	case DurationUnitDay:
		r = t.AddDate(0, 0, -int(d.Value))
	case DurationUnitWeek:
		r = t.AddDate(0, 0, -7*int(d.Value))
	case DurationUnitMonth:
		r = t.AddDate(0, -int(d.Value), 0)
	case DurationUnitYear:
		r = t.AddDate(-int(d.Value), 0, 0)
	default:
		r = t.Add(-time.Duration(d.Value) * time.Second)
	}

	return CpsTime{Time: r}
}

func (d DurationWithUnit) To(unit string) (DurationWithUnit, error) {
	if d.Unit == unit {
		return d, nil
	}

	newDuration := DurationWithUnit{
		Value: d.Value,
		Unit:  unit,
	}

	if d.Value == 0 {
		return newDuration, nil
	}

	var scale int64

	switch d.Unit {
	case DurationUnitMinute:
		if unit == DurationUnitSecond {
			scale = secondsInMinute
		}
	case DurationUnitHour:
		switch unit {
		case DurationUnitMinute:
			scale = minutesInHour
		case DurationUnitSecond:
			scale = secondsInHour
		}
	case DurationUnitDay:
		switch unit {
		case DurationUnitHour:
			scale = hoursInDay
		case DurationUnitMinute:
			scale = minutesInDay
		case DurationUnitSecond:
			scale = secondsInDay
		}
	case DurationUnitWeek:
		switch unit {
		case DurationUnitDay:
			scale = daysInWeek
		case DurationUnitHour:
			scale = hoursInWeek
		case DurationUnitMinute:
			scale = minutesInWeek
		case DurationUnitSecond:
			scale = secondsInWeek
		}
	}

	if scale == 0 {
		return DurationWithUnit{}, fmt.Errorf("cannot transform unit %q to %q", d.Unit, unit)
	}

	newDuration.Value *= scale
	return newDuration, nil
}

func (d DurationWithUnit) String() string {
	unit := d.Unit
	if unit == DurationUnitMonth {
		unit = "mons"
	}
	return fmt.Sprintf("%d%s", d.Value, unit)
}

func (d DurationWithUnit) IsZero() bool {
	return d == DurationWithUnit{}
}

func ParseDurationWithUnit(str string) (DurationWithUnit, error) {
	d := DurationWithUnit{}
	if str == "" {
		return d, fmt.Errorf("invalid duration %q", str)
	}

	r := regexp.MustCompile(`^(-?)(?P<val>\d+)(?P<t>[smhdwMy])$`)
	res := r.FindStringSubmatch(str)
	if len(res) == 0 {
		return d, fmt.Errorf("invalid duration %q", str)
	}

	val, err := strconv.Atoi(res[2])
	if err != nil {
		return d, fmt.Errorf("invalid duration %q: %w", str, err)
	}

	d.Value = int64(val)
	d.Unit = res[3]

	if res[1] == "-" {
		d.Value = -d.Value
	}

	return d, nil
}

type DurationWithEnabled struct {
	DurationWithUnit `bson:",inline"`
	Enabled          *bool `bson:"enabled" json:"enabled" binding:"required"`
}

func IsDurationEnabledAndValid(durationWithEnabled *DurationWithEnabled) bool {
	return durationWithEnabled != nil &&
		durationWithEnabled.Enabled != nil &&
		*durationWithEnabled.Enabled &&
		durationWithEnabled.Value > 0
}
