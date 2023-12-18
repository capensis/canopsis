package datastorage_test

import (
	"testing"
	"time"

	libconfig "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

func TestCanRun_GivenScheduledTimeOnToday_ShouldReturnTrue(t *testing.T) {
	now := time.Now()
	location := now.Location()
	weekday := now.Weekday()
	hour := now.Hour()
	scheduledTime := &libconfig.ScheduledTime{
		Weekday: weekday,
		Hour:    hour,
	}

	if !datastorage.CanRun(datetime.CpsTime{}, scheduledTime, location) {
		t.Errorf("exepcted true but got false")
	}
}

func TestCanRun_GivenNoScheduledTime_ShouldReturnFalse(t *testing.T) {
	now := time.Now()
	location := now.Location()

	if datastorage.CanRun(datetime.CpsTime{}, nil, location) {
		t.Errorf("exepcted false but got true")
	}
}

func TestCanRun_GivenScheduledTimeOnAnotherDay_ShouldReturnFalse(t *testing.T) {
	now := time.Now()
	location := now.Location()
	weekday := now.Weekday()
	if weekday == time.Monday {
		weekday = time.Wednesday
	} else {
		weekday = time.Monday
	}
	scheduledTime := &libconfig.ScheduledTime{
		Weekday: weekday,
		Hour:    10,
	}

	if datastorage.CanRun(datetime.CpsTime{}, scheduledTime, location) {
		t.Errorf("exepcted false but got true")
	}
}

func TestCanRun_GivenScheduledTimeOnTodayAndLastExecutedToday_ShouldReturnFalse(t *testing.T) {
	now := time.Now()
	location := now.Location()
	weekday := now.Weekday()
	hour := now.Hour()
	scheduledTime := &libconfig.ScheduledTime{
		Weekday: weekday,
		Hour:    hour,
	}

	if datastorage.CanRun(datetime.CpsTime{Time: now}, scheduledTime, location) {
		t.Errorf("exepcted false but got true")
	}
}
