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
	newDuration := DurationWithUnit{
		Value: d.Value,
		Unit:  unit,
	}

	if d.Unit == unit || d.Value == 0 {
		return newDuration, nil
	}

	in := int64(0)

	switch d.Unit {
	case DurationUnitMinute:
		if unit == DurationUnitSecond {
			in = 60
		}
	case DurationUnitHour:
		switch unit {
		case DurationUnitMinute:
			in = 60
		case DurationUnitSecond:
			in = 60 * 60
		}
	case DurationUnitDay:
		switch unit {
		case DurationUnitHour:
			in = 24
		case DurationUnitMinute:
			in = 24 * 60
		case DurationUnitSecond:
			in = 24 * 60 * 60
		}
	case DurationUnitWeek:
		switch unit {
		case DurationUnitDay:
			in = 7
		case DurationUnitHour:
			in = 7 * 24
		case DurationUnitMinute:
			in = 7 * 24 * 60
		case DurationUnitSecond:
			in = 7 * 24 * 60 * 60
		}
	}

	if in > 0 {
		newDuration.Value *= in
		return newDuration, nil
	}

	return DurationWithUnit{}, fmt.Errorf("cannot transform unit %q to %q", d.Unit, unit)
}

func (d DurationWithUnit) String() string {
	return fmt.Sprintf("%d%s", d.Value, d.Unit)
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
